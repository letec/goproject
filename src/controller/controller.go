package controller

import (
	"encoding/json"
	"fmt"
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
	rows, err := model.GetUserInfoByID("1")
	fmt.Println(err)
	fmt.Println(rows)
}
