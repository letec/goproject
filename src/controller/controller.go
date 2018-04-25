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

// Test2 路由
func Test2(w http.ResponseWriter, req *http.Request) {
	datas := make(map[int]map[string]interface{})
	datas[0] = make(map[string]interface{})
	datas[0]["username"] = "hugo111"
	datas[0]["password"] = "qqq111"
	datas[0]["salt"] = "1111"
	datas[0]["phone"] = "12121212121"
	datas[0]["bankCode"] = "1212121212121212"
	table := "user"
	res, err := model.InsertRows(table, datas)
	b, _ := json.Marshal(res)
	if err != nil {
		w.Write(b)
	}
}
