package common

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"os"
	"time"
)

// AllParams 取得JSON参数
func AllParams(body []byte, w http.ResponseWriter) (map[string]string, error) {
	info := make(map[string]string)
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

// InMap 判断一个键在不在MAP里面
func InMap(key string, list map[string]string) bool {
	for k := range list {
		if k == key {
			return true
		}
	}
	return false
}

// CheckParamsExist 判断参数是否缺失
func CheckParamsExist(allParam []string, list map[string]string) bool {
	for k := range allParam {
		if InMap(allParam[k], list) == false {
			return false
		}
	}
	return true
}

// PathExists 判断文件路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	return false, err
}

// YmdHis 获取当前 Y-m-d H:i:s 字符串
func YmdHis() string {
	time := time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
	return time
}

// RandInt64 生成随机整型
func RandInt64(min, max int64) int64 {
	maxBigInt := big.NewInt(max)
	i, _ := rand.Int(rand.Reader, maxBigInt)
	iInt64 := i.Int64()
	if iInt64 < min {
		iInt64 = RandInt64(min, max)
	}
	return iInt64
}
