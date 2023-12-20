package db

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/wonderivan/logger"
	"kubeops/config"
)

var Rdb *redis.Client

func CloseRedis() error {
	return Rdb.Close()
}

func InitRedis() {
	ctx := context.Background()
	Rdb = redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: config.RedisPassword, // no password set
		DB:       0,                    // use default DB
	})
	pong, err := Rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	logger.Info("redis 初始化成功", pong)
}
