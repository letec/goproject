package controller

import (
	"common"
	"encoding/json"
	"model"
	"net/http"
)

// SignUp 注册用户
func SignUp(w http.ResponseWriter, req *http.Request, user map[string]string) bool {
	params := []string{"username", "password", "realname"}
	info := map[string]string{}
	ret := common.CheckParamsExist(params, user)
	if ret == false {
		info["code"] = "10002"
		info["msg"] = "参数缺失"
		b, _ := json.Marshal(info)
		w.Write(b)
		return false
	}
	ret = common.ValidSignUp(user)
	if ret == false {
		info["code"] = "10003"
		info["msg"] = "参数验证不通过"
		b, _ := json.Marshal(info)
		w.Write(b)
		return false
	}
	flag, err := model.CheckUserExist(user["username"])
	if err != nil {
		info["code"] = "10004"
		info["msg"] = "网络错误"
		b, _ := json.Marshal(info)
		w.Write(b)
		return false
	}
	if flag == true {
		info["code"] = "20001"
		info["msg"] = "用户名已经存在"
		b, _ := json.Marshal(info)
		w.Write(b)
		return false
	}
	flag, err = model.SignUpUser(user)
	if err != nil || flag == false {
		info["code"] = "10004"
		info["msg"] = "网络错误"
		b, _ := json.Marshal(info)
		w.Write(b)
		return false
	}
	info["code"] = "10000"
	info["msg"] = "注册成功"
	b, _ := json.Marshal(info)
	w.Write(b)
	return true
}
