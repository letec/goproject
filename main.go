package main

import (
	"goproject/src/common"
	"goproject/src/model"
	"goproject/src/router"
	"os"
)

const sysLog = "log/sysLog.log"

// startServer 运行服务器
func startServer() {
	config, err := common.LoadServerConfig()
	if err != nil {
		common.WriteLog(sysLog, "加载server配置文件失败,退出程序!")
		os.Exit(0)
	}
	if !common.CheckParamsExist([]string{"port"}, config) {
		common.WriteLog(sysLog, "server配置文件没有填写正确,退出程序!")
		os.Exit(0)
	}
	port := config["port"]

	router.InitRouter(port)
}

func main() {
	// 连接MYSQL
	model.MysqlConnect()

	// 连接REDIS
	model.RedisConnect()

	// 获取或初始化数据库内的配置信息
	model.InitSysConfig()

	// 开启HTTP服务
	startServer()
}
