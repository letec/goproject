package model

import (
	"common"
	"database/sql"
	"os"
	"strings"

	// 加载MYSQL数据库
	_ "github.com/go-sql-driver/mysql"
)

type Record interface{}

const dbLog = ""

var db *sql.DB

// MysqlConnect 连接MYSQL
func MysqlConnect() {
	flag := false
	errInfo := ""
	config, err := common.LoadMysqlConfig()
	if err != nil {
		errInfo = "MYSQL配置文件读取失败!"
	} else {
		if common.CheckAllParam([]string{"username", "password", "host", "dbname"}, config) {
			str := []string{config["username"], ":", config["password"], "@tcp(", config["host"], ")/", config["dbname"]}
			db, _ = sql.Open("mysql", strings.Join(str, ""))
			db.SetMaxOpenConns(200)
			db.SetMaxIdleConns(10)
			err = db.Ping()
			if err == nil {
				flag = true
			} else {
				errInfo = "数据库连接失败!"
			}
		} else {
			errInfo = "MYSQL配置文件没有填写正确!"
		}
	}
	if flag == false {
		common.WriteLog(sysLogPath, errInfo)
		os.Exit(0)
	}
}

// GetRow 取得单条数据
func GetRow(table string, userDesc []string, where map[string]string) (*sql.Rows, error) {
	set := ""
	length := len(userDesc)
	if length > 0 {
		index := 0
		for k := range userDesc {
			set += userDesc[k]
			if index < length-1 {
				set += ","
			}
			index++
		}
	} else {
		set += "*"
	}
	sql := "SELECT " + set + " FROM " + table + " WHERE "
	length = len(where)
	if length > 0 {
		index := 0
		for k, v := range where {
			sql += k + v
			if index < length-1 {
				sql += " AND "
			}
			index++
		}
	}
	sql += " LIMIT 1"
	rows, err := db.Query(sql)
	if err != nil {
		common.WriteLog(dbLogPath, sql)
		return nil, err
	}
	rows.Next()
	return rows, err
}

// GetRows 取得多条数据
func GetRows(table string, where map[string]string) {

}

// InsertRow 插入单条数据
func InsertRow(table string, data map[string]string) {

}

// InsertRows 插入多条数据
func InsertRows(table string, datas map[string]string) {

}

// DoUpdate 执行更新命令
func DoUpdate() {

}

// DoDelete 执行删除命令
func DoDelete() {

}
