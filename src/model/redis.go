package model

import (
	"encoding/json"
	"goproject/src/common"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
)

// RedisClient 连接池
var RedisClient *redis.Pool

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
				// 建立连接池
				RedisClient = &redis.Pool{
					MaxIdle:   5000,
					MaxActive: 1000,
					Wait:      true,
					Dial: func() (redis.Conn, error) {
						con, err := redis.Dial("tcp", config["host"]+":"+config["port"])
						redis.DialConnectTimeout(5 * time.Second)
						if err != nil {
							errInfo = "REDIS连接失败!"
							return nil, err
						}
						return con, nil
					},
				}
			}
			if err == nil {
				flag = true
			}
			errInfo = "REDIS连接失败!"
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
	Rds := RedisClient.Get()
	defer Rds.Close()
	var val string
	result := make(map[string]string)
	isExit, err := redis.Bool(Rds.Do("EXISTS", key))
	if err == nil && isExit {
		ret, err := redis.String(Rds.Do("GET", key))
		if err != nil {
			val = ret
		}
	}
	err = json.Unmarshal([]byte(val), &result)
	return result, err
}

// RdsSetJSON 设置数据
func RdsSetJSON(key string, val map[string]string, exp ...string) error {
	Rds := RedisClient.Get()
	defer Rds.Close()
	ret, err := json.Marshal(val)
	if err != nil {
		return err
	}
	if len(exp) == 1 {
		_, err = Rds.Do("SET", key, ret, "EX", exp[0])
	} else {
		_, err = Rds.Do("SET", key, ret)
	}
	return err
}

// RedisSet 写入一个key/value值
func RedisSet(key string, val string, exp ...string) error {
	Rds := RedisClient.Get()
	defer Rds.Close()
	var err error
	if len(exp) == 1 {
		_, err = Rds.Do("SET", key, val, "EX", exp[0])
	} else {
		_, err = Rds.Do("SET", key, val)
	}
	return err
}

// RedisGet 拿到一个key/value值
func RedisGet(key string) string {
	Rds := RedisClient.Get()
	defer Rds.Close()
	ret, err := redis.String(Rds.Do("GET", key))
	if err != nil {
		return ""
	}
	return ret
}

// RedisSetHash 设置一个Hash值
func RedisSetHash(key string, offset string, value string) error {
	Rds := RedisClient.Get()
	defer Rds.Close()
	var err error
	_, err = Rds.Do("HSET", key, offset, value)
	Rds.Do("EXPIRE", key, 86400)
	return err
}

// RedisGetHash 设置一个Hash值
func RedisGetHash(key string, offset string) string {
	Rds := RedisClient.Get()
	defer Rds.Close()
	ret, err := redis.String(Rds.Do("HGET", key, offset))
	if err != nil {
		return ""
	}
	return ret
}

// RedisDel 删除
func RedisDel(key string) bool {
	Rds := RedisClient.Get()
	defer Rds.Close()
	_, err := Rds.Do("DEL", key)
	return err == nil
}
