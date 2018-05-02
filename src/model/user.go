package model

import (
	"common"
	"crypto/md5"
	"fmt"
	"io"
)

// GetUserInfoByID 查询用户名
func GetUserInfoByID(userid string) (map[string]interface{}, error) {
	userDesc := []string{"id", "username", "password", "salt", "realname", "phone", "bankCode"}
	where := map[string]string{
		"id": "=" + userid,
	}
	userInfo, err := GetRow("user", userDesc, where)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

// CheckUserExist 查询用户名是否存在
func CheckUserExist(username string) (bool, error) {
	userDesc := []string{"id"}
	where := map[string]string{
		"username=": username,
	}
	userInfo, err := GetRow("user", userDesc, where)
	if err != nil {
		return false, err
	}
	if len(userInfo) > 0 {
		return true, nil
	}
	return false, nil
}

// SignUpUser 注册新用户
func SignUpUser(user map[string]string) (bool, error) {
	user["salt"] = string(common.RandInt64(10000, 99999))
	rpwd := string(user["salt"] + user["password"] + user["username"])
	w := md5.New()
	io.WriteString(w, rpwd)
	user["password"] = fmt.Sprintf("%x", w.Sum(nil)) //w.Sum(nil)将w的hash转成[]byte格式
	ret, err := InsertRow("user", user)
	if err != nil {
		return false, err
	}
	if ret != 1 {
		return false, nil
	}
	return true, nil
}
