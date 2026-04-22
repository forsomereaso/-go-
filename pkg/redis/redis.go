package redis

import (
	"context"
	"game-server/pkg/config"
	"game-server/pkg/logger"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var RDB *redis.Client
var ctx = context.Background()

// 初始化Redis
func Init() {
	cfg := config.Conf.Redis

	RDB = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		logger.Error("redis连接失败", zap.Error(err))
		panic(err)
	}

	logger.Info("Redis连接成功")
}

// 简单封装Set方法，方便后续用
func Set(key, value string, expire time.Duration) error {
	return RDB.Set(ctx, key, value, expire).Err()
}

// Get方法
func Get(key string) (string, error) {
	return RDB.Get(ctx, key).Result()
}
