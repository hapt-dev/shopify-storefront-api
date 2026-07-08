package base

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IBaseRepository[T any, ID any] interface {
	BatchCreate(ctx context.Context, entities []*T) error
	Create(ctx context.Context, entity *T) error
	FindByID(ctx context.Context, id ID) (*T, error)
	First(ctx context.Context, conditions []Condition) (*T, error)
	FindAll(ctx context.Context, conditions []Condition) ([]T, error)
	FindPaged(ctx context.Context, conditions []Condition, paging PagingInput) (*PagingResponseDto, error)
	Update(ctx context.Context, id ID, entity *T) error
	Delete(ctx context.Context, id ID) error
}

type BaseRepository[T any, ID any] struct {
	Collection string
	DB         *mongo.Database
	Redis      *redis.Client
}

func NewBaseRepository[T any, ID any](db *mongo.Database, collection string, redisClient *redis.Client) *BaseRepository[T, ID] {
	repo := &BaseRepository[T, ID]{DB: db, Collection: collection, Redis: redisClient}
	// check compile-time with type assertion
	var _ IBaseRepository[T, ID] = repo
	return repo
}

func (r *BaseRepository[T, ID]) collection() *mongo.Collection {
	return r.DB.Collection(r.Collection)
}

// BatchCreate inserts multiple entities into the collection.
func (r *BaseRepository[T, ID]) BatchCreate(ctx context.Context, entities []*T) error {
	if len(entities) == 0 {
		return errors.Wrap(ErrInvalidInput, "entities list cannot be empty")
	}

	documents := make([]interface{}, 0, len(entities))
	for _, entity := range entities {
		r.applyTimestamps(entity)
		r.applyGeneratedID(entity)
		documents = append(documents, entity)
	}

	_, err := r.collection().InsertMany(ctx, documents)
	return err
}

// Create inserts a single new entity.
func (r *BaseRepository[T, ID]) Create(ctx context.Context, entity *T) error {
	r.applyTimestamps(entity)
	r.applyGeneratedID(entity)

	_, err := r.collection().InsertOne(ctx, entity)
	return err
}

// FindByID returns the entity matching the given _id.
func (r *BaseRepository[T, ID]) FindByID(ctx context.Context, id ID) (*T, error) {
	var result T
	err := r.collection().FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// First returns the first entity matching the given conditions.
func (r *BaseRepository[T, ID]) First(ctx context.Context, conditions []Condition) (*T, error) {
	if len(conditions) == 0 {
		return nil, errors.Wrap(ErrInvalidInput, "conditions cannot be empty")
	}

	var result T
	filter := buildFilter(conditions)
	err := r.collection().FindOne(ctx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// FindAll returns every entity matching the given conditions.
func (r *BaseRepository[T, ID]) FindAll(ctx context.Context, conditions []Condition) ([]T, error) {
	filter := buildFilter(conditions)

	cursor, err := r.collection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []T
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

// FindPaged returns entities for a page, with sorting and the total record count.
func (r *BaseRepository[T, ID]) FindPaged(ctx context.Context, conditions []Condition, paging PagingInput) (*PagingResponseDto, error) {
	paging.Normalize()

	filter := buildFilter(conditions)

	totalCount, err := r.collection().CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	findOptions := options.Find().
		SetLimit(paging.PageSize).
		SetSkip((paging.PageNum - 1) * paging.PageSize)

	if paging.Sort != "" {
		sortOrder := 1
		if paging.SortDesc {
			sortOrder = -1
		}
		findOptions.SetSort(bson.M{paging.Sort: sortOrder})
	}

	cursor, err := r.collection().Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []T
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return &PagingResponseDto{
		Data:        results,
		Count:       totalCount,
		PageNum:     paging.PageNum,
		PageSize:    paging.PageSize,
		HasNextPage: paging.PageNum*paging.PageSize < totalCount,
		HasPrevPage: paging.PageNum > 1,
	}, nil
}

// Update replaces the whole document matching the given _id.
func (r *BaseRepository[T, ID]) Update(ctx context.Context, id ID, entity *T) error {
	if updater, ok := any(entity).(interface{ SetUpdatedAt() }); ok {
		updater.SetUpdatedAt()
	}

	res, err := r.collection().ReplaceOne(ctx, bson.M{"_id": id}, entity)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete removes the document matching the given _id.
func (r *BaseRepository[T, ID]) Delete(ctx context.Context, id ID) error {
	res, err := r.collection().DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}

// applyTimestamps sets timestamps when the entity implements Timestamper.
func (r *BaseRepository[T, ID]) applyTimestamps(entity *T) {
	if ts, ok := any(entity).(Timestamper); ok {
		ts.SetTimestamps()
	}
}

// applyGeneratedID generates a new ObjectID when the ID is still the zero value and its type is primitive.ObjectID.
func (r *BaseRepository[T, ID]) applyGeneratedID(entity *T) {
	setter, ok := any(entity).(IDSetter[ID])
	if !ok {
		return
	}

	var zeroID ID
	if _, isObjectID := any(zeroID).(primitive.ObjectID); !isObjectID {
		return
	}

	idField := reflect.ValueOf(entity).Elem().FieldByName("ID")
	if idField.IsValid() && reflect.DeepEqual(idField.Interface(), zeroID) {
		newID := primitive.NewObjectID()
		setter.SetID(any(newID).(ID))
	}
}

// buildFilter builds a bson filter from the condition list, defaulting the operator to $eq.
func buildFilter(conditions []Condition) bson.M {
	filter := bson.M{}
	for _, cond := range conditions {
		if cond.Field == "" {
			continue
		}
		operator := cond.Operator
		if operator == "" {
			operator = "$eq"
		}
		filter[cond.Field] = bson.M{operator: cond.Value}
	}
	return filter
}
