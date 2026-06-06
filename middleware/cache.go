package middleware

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"xiuyanzhe/redis"
)

type cacheResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *cacheResponseWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

// CacheMiddleware API响应缓存中间件
// 仅缓存GET请求的JSON响应，适用于：
// - 搜索接口缓存30秒
// - 统计接口缓存60秒
// - 单词详情缓存300秒
// 使用Redis存储，结合缓存预热和主动失效
func CacheMiddleware(ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "GET" {
			c.Next()
			return
		}

		// 构建缓存key：包含路径和查询参数
		cacheKey := "cache:api:" + c.Request.URL.Path + "?" + c.Request.URL.RawQuery

		// 尝试从缓存读取
		var cachedData interface{}
		err := redis.GetCache(cacheKey, &cachedData)
		if err == nil && cachedData != nil {
			c.JSON(http.StatusOK, cachedData)
			c.Abort()
			return
		}

		// 包装ResponseWriter以捕获响应
		w := &cacheResponseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = w

		c.Next()

		// 只缓存成功的JSON响应
		if c.Writer.Status() == http.StatusOK {
			contentType := c.Writer.Header().Get("Content-Type")
			if len(contentType) >= 16 && contentType[:16] == "application/json" {
				var data interface{}
				if err := json.Unmarshal(w.body.Bytes(), &data); err == nil {
					if err := redis.SetCache(cacheKey, data, ttl); err != nil {
						log.Printf("[Cache] 写入缓存失败: %v", err)
					}
				}
			}
		}
	}
}

// InvalidateCacheHandler 手动失效缓存的辅助函数（在写操作后调用）
// 调用方式: middleware.InvalidateCacheHandler(prefix)
// 示例: 添加单词后失效搜索缓存
func InvalidateCacheHandler(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 写操作成功后（2xx, 3xx），失效相关缓存
		if c.Writer.Status() < 400 {
			go func() {
				if err := redis.DelCacheByPrefix(prefix); err != nil {
					log.Printf("[Cache] 失效缓存失败 [%s]: %v", prefix, err)
				}
			}()
		}
	}
}
