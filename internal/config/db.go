package config

import (
	"os"
	"strconv"

	"github.com/pkg/errors"
)

// DB holds the DB configuration
type DB struct {
	MongoUri string
	Debug    bool
	RedisUri string
}

var db = &DB{}

// DBCfg returns the default DB configuration
func DBCfg() *DB {
	return db
}

// LoadDBCfg loads DB configuration
func LoadDBCfg() error {
	db.RedisUri = os.Getenv("REDIS_URL")
	db.MongoUri = os.Getenv("DB_MONGO")
	db.Debug, _ = strconv.ParseBool(os.Getenv("DB_DEBUG"))

	if db.MongoUri == "" {
		return errors.New("DB_MONGO is required")
	}

	return nil
}
