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
	info["code"] = "20000"
	info["msg"] = "注册成功"
	b, _ := json.Marshal(info)
	w.Write(b)
	return true
}

// SignIn 用户登陆
func SignIn(w http.ResponseWriter, req *http.Request, user map[string]string) bool {
	params := []string{"username", "password"}
	info := map[string]string{}
	ret := common.CheckParamsExist(params, user)
	if ret == false {
		info["code"] = "10002"
		info["msg"] = "参数缺失"
		b, _ := json.Marshal(info)
		w.Write(b)
		return false
	}
	userDesc := []string{"id", "username", "password", "salt", "status"}
	where := map[string]string{"username=": user["username"]}
	result, err := model.GetRow("user", userDesc, where)
	if err != nil {
		info["code"] = "10004"
		info["msg"] = "网络错误"
		b, _ := json.Marshal(info)
		w.Write(b)
		return false
	}
	if result != nil {
		if result["status"].(int) != 0 {
			info["code"] = "10005"
			info["msg"] = "账号已经被冻结,如有疑问请联系管理员!"
			b, _ := json.Marshal(info)
			w.Write(b)
			return false
		}
		rpwd := result["salt"].(string) + user["password"] + result["username"].(string)
		cpwd := common.MD5(rpwd)
		if cpwd == result["password"] {
			oid, err := model.SetOid(result["id"].(string))
			info["code"] = "10004"
			info["msg"] = "网络错误"
			if err == nil {
				info["code"] = "20000"
				info["msg"] = "登陆成功"
				info["oid"] = oid
			}
			b, _ := json.Marshal(info)
			w.Write(b)
			return true
		}
	}
	info["code"] = "20002"
	info["msg"] = "用户名或密码错误"
	b, _ := json.Marshal(info)
	w.Write(b)
	return false
}
