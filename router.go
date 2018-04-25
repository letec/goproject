package main

import (
	"controller"
	"net/http"
)

func getAllRouters() map[string]func(w http.ResponseWriter, req *http.Request) {
	return map[string]func(w http.ResponseWriter, req *http.Request){
		"/hello": controller.HelloServer,
		"/demo":  controller.Test,
		"/Test2": controller.Test2,
	}
}

// 绑定所有路由
func initRouter() {
	routeList := getAllRouters()

	for k, v := range routeList {
		http.HandleFunc(k, v)
	}
}
