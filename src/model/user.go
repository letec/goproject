package model

import (
	"common"
	"strconv"
)

// GetUserInfoByID 查询用户名
func GetUserInfoByID(userid string) (map[string]interface{}, error) {
	userDesc := []string{"id", "username", "password", "salt", "realname", "phone", "bankCode"}
	where := map[string]string{"id=": userid}
	userInfo, err := GetRow("user", userDesc, where)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

// CheckUserExist 查询用户名是否存在
func CheckUserExist(username string) (bool, error) {
	userDesc := []string{"id"}
	where := map[string]string{"username=": username}
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
	user["salt"] = strconv.FormatInt(common.RandInt64(10000, 99999), 10)
	user["password"] = common.MD5(string(user["salt"] + user["password"] + user["username"]))
	userSQL := createInRowSQL("user", user)

	trans, err := db.Begin()
	if err != nil {
		return false, err
	}
	r, err := trans.Exec(userSQL)
	if err != nil {
		trans.Rollback()
		return false, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		trans.Rollback()
		return false, err
	}
	Info := make(map[string]string)
	Info["userid"] = strconv.FormatInt(id, 10)
	accountSQL := createInRowSQL("account", Info)
	r, err = trans.Exec(accountSQL)
	if err != nil {
		trans.Rollback()
		return false, err
	}
	id, err = r.LastInsertId()
	if err != nil {
		trans.Rollback()
		return false, err
	}
	trans.Commit()
	return true, nil
}
