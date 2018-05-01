package router

import (
	"common"
	"controller"
	"io/ioutil"
	"net/http"
)

// SignUp 注册入口
func SignUp(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	info, err := common.AllParams(body, w)
	if err == nil {
		_ = controller.SignUp(w, req, info)
	}
}

// SignIn 登陆入口
func SignIn(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	info, err := common.AllParams(body, w)
	if err == nil {
		_ = controller.SignIn(w, req, info)
	}
}
