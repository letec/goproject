package main

import (
	"model"
)

func main() {

	// 连接MYSQL
	model.MysqlConnect()

	// 连接REDIS
	model.RedisConnect()

	// 初始化路由
	initRouter()

	// 开启HTTP服务
	startServer()

}
