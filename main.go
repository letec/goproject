package main

import (
	"common"
	"model"
	"net/http"
	"os"
	"router"
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
	err = http.ListenAndServe("www.xr.com:"+port, nil)
	if err != nil {
		common.WriteLog(sysLog, "监听"+port+"端口失败,退出程序!")
		os.Exit(0)
	}
}

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
	startServer()
}
