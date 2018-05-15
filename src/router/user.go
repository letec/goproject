package router

import (
	"controller"
	"net/http"
)

// SignUp 注册入口
func SignUp(w http.ResponseWriter, req *http.Request) {
	info, err := AllParams(w, req)
	if err == nil {
		_ = controller.SignUp(w, req, info)
	}
}

// SignIn 登陆入口
func SignIn(w http.ResponseWriter, req *http.Request) {
	info, err := AllParams(w, req)
	if err == nil {
		_ = controller.SignIn(w, req, info)
	}
}
