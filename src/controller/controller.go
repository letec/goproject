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
	userid := "1010"
	// 定义一个(切片)类型,类似于PHP一维数组
	userDesc := []string{"id", "username", "realname", "phone", "bankCode"}
	// WHERE条件写进map
	where := map[string]string{
		"id": "=" + userid,
	}
	userInfo, err := model.GetRows("user", userDesc, where, 0)
	// 报错了就返回空和错误
	if err == nil {
		b, _ := json.Marshal(userInfo)
		w.Write(b)
	}
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

// Test3 路由
func Test3(w http.ResponseWriter, req *http.Request) {
	datas := make(map[string]string)
	where := make(map[string]string)
	datas["username"] = "mixdran"
	datas["salt"] = "2222"
	where["username="] = "hugo111"
	table := "user"
	res, err := model.DoUpdate(table, datas, where)
	b, _ := json.Marshal(res)
	if err != nil {
		w.Write(b)
	}
}

// Test4 路由
func Test4(w http.ResponseWriter, req *http.Request) {
	where := make(map[string]string)
	where["username="] = "hugo111"
	where["password="] = "qqq111"
	table := "user"
	res, err := model.DoDelete(table, where)
	b, _ := json.Marshal(res)
	if err == nil {
		w.Write(b)
	}
}

// Test5 路由
func Test5(w http.ResponseWriter, req *http.Request) {
	fmt.Println(model.CheckUserExist("mixdr"))
}
