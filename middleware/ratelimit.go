package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"xiuyanzhe/redis"
)

type RateLimitConfig struct {
	Window  time.Duration // 时间窗口
	MaxReq  int           // 窗口内最大请求数
	Burst   int           // 突发允许量
	KeyFunc func(c *gin.Context) string
}

var (
	// GeneralLimit 通用API限制: 每分钟60次，允许突发80
	GeneralLimit = RateLimitConfig{
		Window: time.Minute,
		MaxReq: 60,
		Burst:  80,
		KeyFunc: func(c *gin.Context) string {
			return "ratelimit:general:" + c.ClientIP()
		},
	}

	// AuthLimit 登录注册限制: 每分钟10次
	AuthLimit = RateLimitConfig{
		Window: time.Minute,
		MaxReq: 10,
		Burst:  15,
		KeyFunc: func(c *gin.Context) string {
			return "ratelimit:auth:" + c.ClientIP()
		},
	}

	// SearchLimit 搜索限制: 每分钟30次
	SearchLimit = RateLimitConfig{
		Window: time.Minute,
		MaxReq: 30,
		Burst:  40,
		KeyFunc: func(c *gin.Context) string {
			uid, _ := c.Get("user_id")
			return "ratelimit:search:user:" + strconv.Itoa(int(uid.(uint)))
		},
	}

	// LearnLimit 学习API限制: 每分钟120次
	LearnLimit = RateLimitConfig{
		Window: time.Minute,
		MaxReq: 120,
		Burst:  150,
		KeyFunc: func(c *gin.Context) string {
			uid, _ := c.Get("user_id")
			return "ratelimit:learn:user:" + strconv.Itoa(int(uid.(uint)))
		},
	}
)

// RateLimitMiddleware 基于Redis滑动窗口的限流中间件
// 使用Redis有序集合(Sorted Set)，以毫秒时间戳为score
// 每次请求清理窗口外的旧记录，再判断当前窗口内请求数
func RateLimitMiddleware(config RateLimitConfig) gin.HandlerFunc {
	now := time.Now()
	_ = now
	return func(c *gin.Context) {
		key := config.KeyFunc(c)
		windowMs := config.Window.Milliseconds()
		nowMs := time.Now().UnixMilli()
		threshold := nowMs - windowMs

		ctx := c.Request.Context()

		// 使用Redis EVAL执行Lua脚本保证原子性
		script := `
			redis.call('ZREMRANGEBYSCORE', KEYS[1], 0, ARGV[1])
			local count = redis.call('ZCARD', KEYS[1])
			if tonumber(count) < tonumber(ARGV[2]) then
				redis.call('ZADD', KEYS[1], ARGV[3], ARGV[4])
				redis.call('EXPIRE', KEYS[1], ARGV[5])
				return {1, count + 1}
			end
			return {0, count}
		`

		result, err := redis.Client.Eval(ctx, script, []string{key},
			threshold,               // ARGV[1]: 窗口起始时间戳
			config.MaxReq+config.Burst, // ARGV[2]: 最大容量
			nowMs,                   // ARGV[3]: 当前时间戳(score)
			key+":"+strconv.FormatInt(nowMs, 10), // ARGV[4]: 唯一成员
			int(config.Window.Seconds())+10, // ARGV[5]: TTL = 窗口+10秒
		).Result()

		if err != nil {
			// Redis不可用时，放行（降级策略）
			c.Next()
			return
		}

		vals, ok := result.([]interface{})
		if !ok || len(vals) < 2 {
			c.Next()
			return
		}

		allowed, _ := vals[0].(int64)
		currentCount, _ := vals[1].(int64)

		// 设置响应头，告知客户端限流状态
		c.Header("X-RateLimit-Limit", strconv.Itoa(config.MaxReq))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(max(0, config.MaxReq-int(currentCount))))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(nowMs+windowMs, 10))

		if allowed == 0 {
			c.Header("Retry-After", strconv.Itoa(int(config.Window.Seconds())))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			return
		}

		c.Next()
	}
}
