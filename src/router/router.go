package router

import "net/http"

var routers = map[string]func(w http.ResponseWriter, req *http.Request){
	"/signin": SignIn,
	"/signup": SignUp,
}

// InitRouter 绑定所有路由
func InitRouter() {
	for k, v := range routers {
		http.HandleFunc(k, v)
	}
}
