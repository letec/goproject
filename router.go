package main

import (
	"net/http"
	"router"
)

var routers = map[string]func(w http.ResponseWriter, req *http.Request){
	"/SignUp": router.SignUp,
	"/SignIn": router.SignIn,
}

// 绑定所有路由
func initRouter() {
	for k, v := range routers {
		http.HandleFunc(k, v)
	}
}
