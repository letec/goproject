package model

import (
	"encoding/json"
	"goproject/src/common"
	"os"

	"github.com/garyburd/redigo/redis"
)

var rds redis.Conn

// RedisConnect 连接redis
func RedisConnect() {
	flag := false
	errInfo := ""
	var err error
	config := make(map[string]string)
	config, err = common.LoadRedisConfig()
	if err != nil {
		errInfo = "REDIS配置文件读取失败!"
	} else {
		if common.InMap("host", config) && common.InMap("port", config) {
			if err != nil {
				errInfo = "REDIS配置文件读取失败!"
			} else {
				rds, err = redis.Dial("tcp", "127.0.0.1:6379")
				if err == nil {
					flag = true
				}
				errInfo = "REDIS连接失败!"
			}
		} else {
			errInfo = "REDIS配置文件没有填写正确!"
		}
	}
	if flag == false {
		common.WriteLog(sysLogPath, errInfo)
		os.Exit(0)
	}
}

// RdsGetJSON 取得数据
func RdsGetJSON(key string) (map[string]string, error) {
	var val string
	result := make(map[string]string)
	isExit, err := redis.Bool(rds.Do("EXISTS", key))
	if err == nil && isExit {
		ret, err := redis.String(rds.Do("GET", key))
		if err != nil {
			val = ret
		}
	}
	err = json.Unmarshal([]byte(val), &result)
	return result, err
}

// RdsSetJSON 设置数据
func RdsSetJSON(key string, val map[string]string, exp ...string) error {
	ret, err := json.Marshal(val)
	if err != nil {
		return err
	}
	if len(exp) == 1 {
		_, err = rds.Do("SET", key, ret, "EX", exp[0])
	} else {
		_, err = rds.Do("SET", key, ret)
	}
	return err
}

// RedisSet 写入一个key/value值
func RedisSet(key string, val string, exp ...string) error {
	var err error
	if len(exp) == 1 {
		_, err = rds.Do("SET", key, val, "EX", exp[0])
	} else {
		_, err = rds.Do("SET", key, val)
	}
	return err
}

// RedisGet 拿到一个key/value值
func RedisGet(key string) string {
	ret, err := redis.String(rds.Do("GET", key))
	if err != nil {
		return ""
	}
	return ret
}
