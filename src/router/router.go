package router

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

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

// AllParams 取得JSON参数
func AllParams(w http.ResponseWriter, req *http.Request) (map[string]string, error) {
	info := make(map[string]string)
	body, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(body, &info)
	if err == nil {
		return info, nil
	}
	info["code"] = "10001"
	info["msg"] = "参数结构错误"
	b, _ := json.Marshal(info)
	w.Write(b)
	return nil, err
}
