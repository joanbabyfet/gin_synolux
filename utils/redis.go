package utils

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var (
	Redis        *redis.Client
	RedisTimeout time.Duration //默认过期时间1天
)

func init() {
	//redis初始化
	Redis = InitRedis()
	RedisTimeout = 86400 * time.Second
}

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.host"),
		Password: viper.GetString("redis.password"), // no password set
		DB:       viper.GetInt("redis.database"),    // use default DB
	})
	result := rdb.Ping()
	if result.Val() != "PONG" {
		// 连接有问题
		return nil
	}
	return rdb
}
