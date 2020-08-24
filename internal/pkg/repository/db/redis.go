package db

import (
	"github.com/go-redis/redis/v8"
	"github.com/whale-team/whaleEcho/internal/pkg/app"
	"go.uber.org/fx"
)

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Params struct {
	fx.In
	Redis *redis.Client
}

func New(params Params) app.Repositorier {
	return &Repo{
		redisDB: params.Redis,
	}
}

func NewRedis(config RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password, // no password set
		DB:       config.DB,
	})
	return rdb
}
