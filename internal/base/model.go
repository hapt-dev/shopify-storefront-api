package base

import (
	"time"
)

type BaseModel[T any] struct {
	ID        T         `bson:"_id" json:"id"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

// Timestamper lets the base repository set created/updated timestamps.
type Timestamper interface {
	SetTimestamps()
	SetUpdatedAt()
}

func (m *BaseModel[T]) SetTimestamps() {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
}

func (m *BaseModel[T]) SetUpdatedAt() {
	m.UpdatedAt = time.Now()
}

// IDSetter lets the base repository assign an auto-generated ID when the entity has none.
type IDSetter[T any] interface {
	SetID(id T)
}
