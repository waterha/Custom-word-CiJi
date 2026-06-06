package main

import (
	"xiuyanzhe/routes"
	"xiuyanzhe/database"
	"xiuyanzhe/middleware"
	"xiuyanzhe/redis"
)

func main() {
	database.InitDB()
	redis.InitRedis()

	// 初始化数据库索引
	middleware.EnsureIndexes()

	// 启动异步队列处理器（批量写入非关键数据，降低数据库压力）
	redis.StartAsyncQueue()
	defer redis.StopAsyncQueue()

	r := routes.SetupRouter()
	r.Run(":8080")
}
