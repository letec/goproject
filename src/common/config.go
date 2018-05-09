package common

import (
	"encoding/json"
	"io/ioutil"
)

func getJSONInFile(filePath string) (map[string]string, error) {
	config := make(map[string]string)
	data, err := ioutil.ReadFile(filePath)
	if err == nil {
		datajson := []byte(data)
		err = json.Unmarshal(datajson, &config)
	}
	return config, err
}

// LoadRedisConfig 加载redis配置文件
func LoadRedisConfig() (map[string]string, error) {
	filePath := "config/redis.json"
	ret, err := getJSONInFile(filePath)
	return ret, err
}

// LoadMysqlConfig 加载mysql配置文件
func LoadMysqlConfig() (map[string]string, error) {
	filePath := "config/mysql.json"
	ret, err := getJSONInFile(filePath)
	return ret, err
}

// LoadServerConfig 加载server配置文件
func LoadServerConfig() (map[string]string, error) {
	filePath := "config/server.json"
	ret, err := getJSONInFile(filePath)
	return ret, err
}
