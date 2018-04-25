package model

import (
	"common"
	"database/sql"
	"os"
	"strings"

	// 加载MYSQL数据库
	_ "github.com/go-sql-driver/mysql"
)

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

// SCAN ROW PARAMS
func scanAllParams(rows *sql.Rows) map[string]interface{} {
	// 取得所有列
	columns, _ := rows.Columns()
	// 计算列的数量
	length := len(columns)
	// 定义一个存放值的map
	value := make([]interface{}, length)
	// 定义一个存放指针的map
	columnPointers := make([]interface{}, length)
	// 循环把变量地址拿到 &代表地址
	for i := 0; i < length; i++ {
		columnPointers[i] = &value[i]
	}
	// 这个是最经典的 把一个interface的切片作为可变参数传进去scan 有多少列都在columnPointers里面了
	rows.Scan(columnPointers...)
	// 创建结果变量存返回值
	result := make(map[string]interface{})
	// 循环把结果放到result里面
	for i := 0; i < length; i++ {
		columnName := columns[i]
		columnValue := columnPointers[i].(*interface{})
		result[columnName] = string((*columnValue).([]uint8))
	}
	return result
}

// CREATE SELECT SQL
func getSelectSQL(table string, userDesc []string, where map[string]string) string {
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
	return sql
}

// GetRow 取得单条数据
func GetRow(table string, userDesc []string, where map[string]string) (map[string]interface{}, error) {
	sql := getSelectSQL(table, userDesc, where) + " LIMIT 1"
	rows, err := db.Query(sql)
	if err != nil {
		common.WriteLog(dbLogPath, sql)
		return nil, err
	}
	rows.Next()
	result := scanAllParams(rows)
	return result, err
}

// GetRows 取得多条数据
func GetRows(table string, userDesc []string, where map[string]string) (map[int]map[string]interface{}, error) {
	sql := getSelectSQL(table, userDesc, where)
	rows, err := db.Query(sql)
	if err != nil {
		common.WriteLog(dbLogPath, sql)
		return nil, err
	}
	var result = make(map[int]map[string]interface{})
	index := 0
	for rows.Next() {
		result[index] = scanAllParams(rows)
	}
	return result, err
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
