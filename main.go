package main

import (
	"model"
	"router"
)

func main() {
	// 连接MYSQL
	model.MysqlConnect()

	// 连接REDIS
	model.RedisConnect()

	// 初始化路由
	router.InitRouter()

	// 开启HTTP服务
	router.StartServer()
}
