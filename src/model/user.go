package model

import (
	"fmt"
	"goproject/src/common"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

// CheckOnline 根据oid查询用户是否在线 重置过期时间 并返回userid
func CheckOnline(oid string) (string, error) {
	Rds := RedisClient.Get()
	defer Rds.Close()
	userid, err := redis.String(Rds.Do("GET", oid))
	if err == nil {
		dOid, err := redis.String(Rds.Do("GET", "SESSION_"+userid))
		if err == nil && dOid == oid {
			Rds.Do("EXPIRE", oid, 1200)
			Rds.Do("EXPIRE", "SESSION_"+userid, 1200)
			return userid, err
		}
	}
	Rds.Do("DELETE", oid)
	Rds.Do("DELETE", "SESSION_"+userid)
	return "", err
}

// GetUserInfoByID 根据userid查询用户资料
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
	where := map[string]string{"UserName=": username}
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
	insert := make(map[string]string)
	salt := strconv.FormatInt(common.RandInt64(10000, 99999), 10)
	insert["UserName"] = user["username"]
	insert["PassWord"] = common.MD5(salt + user["password"] + user["username"])
	insert["Salt"] = salt
	insert["CreateTime"] = strconv.FormatInt(time.Now().Unix(), 10)
	userSQL := createInRowSQL("user", insert)
	trans, err := db.Begin()
	if err != nil {
		return false, err
	}
	r, err := trans.Exec(userSQL)
	if err != nil {
		fmt.Println(err)
		trans.Rollback()
		return false, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		fmt.Println(err)
		trans.Rollback()
		return false, err
	}
	Info := make(map[string]string)
	Info["UserID"] = strconv.FormatInt(id, 10)
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

// SetOid 生成并向REDIS写入在线OID
func SetOid(userid string) (string, error) {
	Rds := RedisClient.Get()
	time := strconv.FormatInt(time.Now().UnixNano(), 10)
	oid := common.MD5(time + userid)
	_, err := Rds.Do("SET", oid, userid, "EX", 1200)
	_, err = Rds.Do("SET", "SESSION_"+userid, oid, "EX", 1200)
	defer Rds.Close()
	return oid, err
}

// GetUserInfoInHall 获取大厅中用户信息
func GetUserInfoInHall(allID []string) (map[int]map[string]interface{}, error) {
	var result = make(map[int]map[string]interface{})
	if len(allID) == 0 {
		return result, nil
	}
	ids := strings.Join(allID, ",")
	sql := fmt.Sprintf("SELECT a.*,b.Avatar FROM `user` LEFT JOIN `account` ON a.id = b.UserID WHERE a.id IN (%v)", ids)
	rows, err := db.Query(sql)
	if err != nil {
		common.WriteLog(dbLogPath, sql)
		return nil, err
	}
	index := 0
	for rows.Next() {
		result[index] = scanAllParams(rows)
		index++
	}
	return result, err
}
