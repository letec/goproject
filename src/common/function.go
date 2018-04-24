package common

import (
	"os"
	"time"
)

// InMap 判断一个键在不在MAP里面
func InMap(key string, list map[string]string) bool {
	for k := range list {
		if k == key {
			return true
		}
	}
	return false
}

// CheckAllParam 判断参数是否缺失
func CheckAllParam(allParam []string, list map[string]string) bool {
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

// YmdHis 获取 Y-m-d H:i:s 字符串
func YmdHis() string {
	time := time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
	return time
}
