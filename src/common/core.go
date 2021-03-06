package common

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"os"
	"time"
)

// InSlice 判断一个值是否在切片中
func InSlice(str string, slice []string) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == str {
			return true
		}
	}
	return false
}

// MapKeys 将一个map的键排列为切片
func MapKeys(data map[string]string) []string {
	result := []string{}
	for k := range data {
		result = append(result, k)
	}
	return result
}

// MapValues 将一个map的值排列为切片
func MapValues(data map[string]string) []string {
	result := []string{}
	for _, v := range data {
		result = append(result, v)
	}
	return result
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

// MD5 加密一个字符串为MD5
func MD5(context string) string {
	newMD5 := md5.New()
	io.WriteString(newMD5, context)
	result := fmt.Sprintf("%x", newMD5.Sum(nil))
	return result
}
