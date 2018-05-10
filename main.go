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

	// 获取或初始化数据库内的配置信息
	model.InitSysConfig()

	// 初始化路由
	router.InitRouter()

	// 开启HTTP服务
	StartServer()
}
