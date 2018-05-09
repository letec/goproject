package router

import (
	"common"
	"net/http"
	"os"
)

const sysLog = "log/sysLog.log"

// StartServer 运行服务器
func StartServer() {
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
