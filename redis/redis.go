package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"os"
)

var Client *redis.Client
var Ctx = context.Background()

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func InitRedis() {
	host := getEnv("REDIS_HOST", "redis") // 本地测试改成localhost
	port := getEnv("REDIS_PORT", "6379")
	password := getEnv("REDIS_PASSWORD", "")

	Client = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password, // 无密码
		DB:       0,        // 默认DB
	})

	// 测试连接
	_, err := Client.Ping(Ctx).Result()
	if err != nil {
		panic("连接Redis失败")
	}
}

func Set(key string, value interface{}, expiration int) error {
	return Client.Set(Ctx, key, value, 0).Err()
}

func Get(key string) (string, error) {
	return Client.Get(Ctx, key).Result()
}

func Del(key string) error {
	return Client.Del(Ctx, key).Err()
}

func LPush(key string, values ...interface{}) error {
	return Client.LPush(Ctx, key, values...).Err()
}

func RPop(key string) (string, error) {
	return Client.RPop(Ctx, key).Result()
}

func LLen(key string) (int64, error) {
	return Client.LLen(Ctx, key).Result()
}
