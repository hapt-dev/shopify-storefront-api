// Package app initializes the shared infrastructure components (config, logger,
// Mongo, Redis) so they can be reused by both the HTTP server and service jobs.
package app

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/ALLINBYTE-COMPANY/shopify-storefront-api/internal/config"
	"github.com/ALLINBYTE-COMPANY/shopify-storefront-api/platform/database"
	"github.com/ALLINBYTE-COMPANY/shopify-storefront-api/platform/logger"
)

// App holds the initialized infrastructure dependencies together with the modules (repos/services).
type App struct {
	Config  *config.App
	DB      *config.DB
	Log     *zap.Logger
	Mongo   *mongo.Client
	MongoDB *mongo.Database
	Redis   *redis.Client

	Repositories *Repositories
	Services     *Services
}

// New loads env, config, logger and opens the Mongo/Redis connections.
// It does not touch HTTP/router so it can be reused by service jobs.
func New() (*App, error) {
	_ = godotenv.Load()

	if err := config.LoadApiCfg(); err != nil {
		return nil, err
	}

	if err := logger.SetUpLogger(); err != nil {
		return nil, err
	}

	dbCfg := config.DBCfg()

	mongoInst, err := database.InitMongoConnection(dbCfg.MongoUri)
	if err != nil {
		return nil, err
	}

	var redisClient *redis.Client
	redisClient, err = database.InitRedisConnection(dbCfg.RedisUri)
	if err != nil {
		return nil, err
	}

	a := &App{
		Config:  config.GetAppCfg(),
		DB:      dbCfg,
		Log:     logger.GetLogger(),
		Mongo:   mongoInst.Client,
		MongoDB: mongoInst.DB,
		Redis:   redisClient,
	}

	a.initModules()

	return a, nil
}

// Close releases infrastructure resources: closes Mongo/Redis and flushes the logger.
func (a *App) Close(ctx context.Context) {
	if a.Mongo != nil {
		if err := a.Mongo.Disconnect(ctx); err != nil {
			a.Log.Error("mongo disconnect error", zap.Error(err))
		}
	}
	if a.Redis != nil {
		if err := a.Redis.Close(); err != nil {
			a.Log.Error("redis close error", zap.Error(err))
		}
	}
	_ = logger.Sync()
}
