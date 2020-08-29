package db

import (
	"github.com/go-redis/redis/v8"
	"github.com/whale-team/whaleEcho/internal/pkg/app"
	"go.uber.org/fx"
)

// RedisConfig config for redis db
type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// Params redis db dependencies
type Params struct {
	fx.In
	Redis *redis.Client
}

// New create a repo instance
func New(params Params) app.Repositorier {
	return &Repo{
		redisDB: params.Redis,
	}
}

// NewRedis create a redis client instance
func NewRedis(config RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password, // no password set
		DB:       config.DB,
	})
	return rdb
}
