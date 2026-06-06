package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"xiuyanzhe/controllers"
	"xiuyanzhe/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORSMiddleware())

	// 提供前端静态文件
	r.Static("/assets", "./frontend/dist/assets")

	// 提供根目录下的静态文件
	r.StaticFile("/favicon.svg", "./frontend/dist/favicon.svg")
	r.StaticFile("/icons.svg", "./frontend/dist/icons.svg")

	// 公共路由（带限流）
	r.POST("/api/register", middleware.RateLimitMiddleware(middleware.AuthLimit), controllers.Register)
	r.POST("/api/login", middleware.RateLimitMiddleware(middleware.AuthLimit), controllers.Login)
	// 健康检查（不限流）
	r.GET("/api/health", controllers.HealthCheck)

	// 需要认证的路由
	auth := r.Group("/api")
	auth.Use(middleware.JWTMiddleware())
	auth.Use(middleware.RateLimitMiddleware(middleware.GeneralLimit))
	{
		// 学习路由
		learn := auth.Group("/learn")
		learn.Use(middleware.RateLimitMiddleware(middleware.LearnLimit))
		{
			// 搜索：缓存30秒
			learn.GET("/search", middleware.CacheMiddleware(30*time.Second), controllers.SearchWords)
			// 统计：缓存60秒
			learn.GET("/stats", middleware.CacheMiddleware(60*time.Second), controllers.GetLearningStats)
			learn.GET("/next", controllers.GetNextWord)
			learn.GET("/words/:id", controllers.GetWord)
			learn.POST("/answer", controllers.SubmitAnswer)
			learn.GET("/progress", controllers.GetProgress)

			// 自定义单词学习路由
			learn.GET("/custom/next", controllers.GetNextCustomWord)
			learn.POST("/custom/answer", controllers.SubmitCustomAnswer)
			learn.GET("/custom/progress", controllers.GetCustomProgress)
		}

		// 管理员路由
		admin := auth.Group("/admin")
		admin.Use(middleware.AdminMiddleware())
		{
			// 分页查询（分页后不缓存，每次请求量小无需缓存）
			admin.GET("/words", controllers.GetWords)
			admin.GET("/words/:id", middleware.CacheMiddleware(30*time.Second), controllers.GetWord)
			// 写操作同时失效搜索缓存和管理员单词列表缓存
			admin.POST("/words", middleware.InvalidateCacheHandler("cache:search:"), middleware.InvalidateCacheHandler("cache:api:/api/admin/words"), controllers.AddWord)
			admin.PUT("/words/:id", middleware.InvalidateCacheHandler("cache:search:"), middleware.InvalidateCacheHandler("cache:api:/api/admin/words"), controllers.UpdateWord)
			admin.DELETE("/words/:id", middleware.InvalidateCacheHandler("cache:search:"), middleware.InvalidateCacheHandler("cache:api:/api/admin/words"), controllers.DeleteWord)
		}

		// 监控路由
		monitor := auth.Group("/monitor")
		monitor.Use(middleware.AdminMiddleware())
		{
			monitor.GET("/overview", middleware.CacheMiddleware(30*time.Second), controllers.GetMonitorOverview)
			monitor.GET("/hourly-visits", middleware.CacheMiddleware(30*time.Second), controllers.GetHourlyVisits)
			monitor.GET("/daily-registrations", middleware.CacheMiddleware(30*time.Second), controllers.GetDailyRegistrations)
		}

		// 用户自定义单词路由
		custom := auth.Group("/custom")
		{
			custom.GET("/words", middleware.CacheMiddleware(30*time.Second), controllers.GetCustomWords)
			custom.GET("/words/:id", middleware.CacheMiddleware(30*time.Second), controllers.GetCustomWord)
			custom.POST("/words", middleware.InvalidateCacheHandler("cache:search:"), controllers.AddCustomWord)
			custom.PUT("/words/:id", middleware.InvalidateCacheHandler("cache:search:"), controllers.UpdateCustomWord)
			custom.DELETE("/words/:id", middleware.InvalidateCacheHandler("cache:search:"), controllers.DeleteCustomWord)
		}
	}

	// 处理前端路由，确保单页应用的路由正常工作
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	return r
}
