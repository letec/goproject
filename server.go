package main

import (
	"common"
	"net/http"
	"os"
)

const sysLog = "log/sysLog.log"

// startServer 运行服务器
func startServer() {
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		common.WriteLog(sysLog, "监听8080端口失败,退出程序!")
		os.Exit(0)
	}
}
