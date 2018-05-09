package main

import (
	"net/http"
	"router"
)

func getAllRouters() map[string]func(w http.ResponseWriter, req *http.Request) {
	return map[string]func(w http.ResponseWriter, req *http.Request){
		"/SignUp": router.SignUp,
		"/SignIn": router.SignIn,
	}
}

// 绑定所有路由
func initRouter() {
	routeList := getAllRouters()
	for k, v := range routeList {
		http.HandleFunc(k, v)
	}
}
