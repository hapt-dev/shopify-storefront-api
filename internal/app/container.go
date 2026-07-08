package app

import (
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ALLINBYTE-COMPANY/shopify-storefront-api/internal/repository"
	"github.com/ALLINBYTE-COMPANY/shopify-storefront-api/internal/service"
)

// Repositories groups the repositories of every module.
type Repositories struct {
	Store repository.IStoreRepository
}

// Services groups the services of every module.
type Services struct {
	Store service.IStoreService
}

// initModules builds the repositories and then the services from the initialized infrastructure (Mongo/Redis...).
func (a *App) initModules() {
	a.Repositories = newRepositories(a.MongoDB, a.Redis)
	a.Services = newServices(a.Repositories)
}

func newRepositories(db *mongo.Database, redisClient *redis.Client) *Repositories {
	_ = db
	return &Repositories{
		Store: repository.NewStoreRepository(redisClient),
	}
}

func newServices(repos *Repositories) *Services {
	return &Services{
		Store: service.NewStoreService(repos.Store),
	}
}
