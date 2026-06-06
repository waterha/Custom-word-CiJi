package redis

import (
	"encoding/json"
	"time"
)

// CacheConfig 缓存配置
type CacheConfig struct {
	TTL      time.Duration // 过期时间
	Prefix   string        // Key前缀
	KeyFunc  func(args ...string) string
}

// 缓存键前缀常量
const (
	CachePrefixSearch     = "cache:search:"
	CachePrefixStats      = "cache:stats:"
	CachePrefixWord       = "cache:word:"
	CachePrefixCustomWord = "cache:custom_word:"
)

// SetCache 将数据存入Redis缓存（JSON序列化）
func SetCache(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return Client.Set(Ctx, key, string(data), ttl).Err()
}

// GetCache 从Redis缓存获取数据（JSON反序列化）
func GetCache(key string, dest interface{}) error {
	data, err := Client.Get(Ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

// DelCache 删除缓存
func DelCache(key string) error {
	return Client.Del(Ctx, key).Err()
}

// DelCacheByPrefix 按前缀批量删除缓存（用于数据变更后的缓存清理）
// 使用Redis SCAN命令逐条删除，避免KEYS阻塞
func DelCacheByPrefix(prefix string) error {
	iter := Client.Scan(Ctx, 0, prefix+"*", 100).Iterator()
	keys := make([]string, 0, 100)
	for iter.Next(Ctx) {
		keys = append(keys, iter.Val())
		if len(keys) >= 100 {
			if err := Client.Del(Ctx, keys...).Err(); err != nil {
				return err
			}
			keys = keys[:0]
		}
	}
	if len(keys) > 0 {
		return Client.Del(Ctx, keys...).Err()
	}
	return iter.Err()
}
