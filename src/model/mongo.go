package model

import (
	"goproject/src/common"
	"os"

	"gopkg.in/mgo.v2"
)

var Mongo *mgo.Session

// MongoDBConnection 连接mongo
func MongoDBConnection() {
	flag := false
	errInfo := ""
	var err error
	config := make(map[string]string)
	config, err = common.LoadMongoConfig()
	if err != nil {
		errInfo = "REDIS配置文件读取失败!"
	} else {
		if common.InMap("dsn", config) {
			Mongo, err = mgo.Dial(config["dsn"])
			errInfo = "MongoDB连接失败!"
			if err == nil {
				flag = true
			}
		}
	}
	if flag == false {
		common.WriteLog(sysLogPath, errInfo)
		os.Exit(0)
	}
}

// MongoInsert 插入数据集
func MongoInsert(dbName string, Collection string, collect *mgo.Collection) bool {
	con := Mongo.Copy()
	defer con.Close()
	c := con.DB(dbName).C(Collection)
	err := c.Insert(collect)
	if err == nil {
		return true
	}
	return true
}
