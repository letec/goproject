package controller

import (
	"encoding/json"
	"io/ioutil"
	"model"
	"net/http"
)

// HelloServer 输出
func HelloServer(w http.ResponseWriter, req *http.Request) {
	var user map[string]interface{}
	body, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(body, &user)
	w.Write(body)
}

// Test 路由
func Test(w http.ResponseWriter, req *http.Request) {
	rows, err := model.GetUserInfoByIDs("0")
	result := map[string]interface{}{"code": "4001"}
	if err == nil {
		result["code"] = "2006"
		result["user"] = rows
	}
	b, _ := json.Marshal(result)
	w.Write(b)
}
