package model

import (
	"common"
	"database/sql"
	"errors"
	"os"
	"strconv"
	"strings"

	// 加载MYSQL驱动
	_ "github.com/go-sql-driver/mysql"
)

const dbLog = ""

var db *sql.DB

// GetDB 获取数据库链接资源
func GetDB() *sql.DB {
	return db
}

// MysqlConnect 连接MYSQL
func MysqlConnect() {
	flag := false
	errInfo := ""
	config, err := common.LoadMysqlConfig()
	if err != nil {
		errInfo = "MYSQL配置文件读取失败!"
	} else {
		if common.CheckParamsExist([]string{"dbm_username", "dbm_password", "dbm_host", "dbm_dbname"}, config) {
			str := []string{config["dbm_username"], ":", config["dbm_password"], "@tcp(", config["dbm_host"], ")/", config["dbm_dbname"]}
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
func createSelectSQL(table string, userDesc []string, where map[string]string, limit int) string {
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
	sql := "SELECT " + set + " FROM " + table + " WHERE 1 = 1 "
	length = len(where)
	if len(where) > 0 {
		index := 0
		for k, v := range where {
			if index < length-1 {
				sql += " AND "
			}
			sql += k + "'" + v + "'"
			index++
		}
	}
	if limit > 0 {
		sql += " LIMIT " + strconv.Itoa(limit)
	}
	return sql
}

func createInRowSQL(table string, data map[string]string) string {
	length := len(data)
	keys, vals := "", ""
	index := 0
	for k, v := range data {
		keys += k
		vals += "'" + v + "'"
		if index < length-1 {
			keys += ","
			vals += ","
		}
		index++
	}
	sql := "INSERT INTO " + table + "(" + keys + ") VALUES (" + vals + ")"
	return sql
}

func createInRowsSQL(table string, datas map[int]map[string]interface{}, length int) string {
	keys, vals := "", ""
	index := 0
	colLength := len(datas[0])
	colName := make(map[int]string)
	for k := range datas[0] {
		colName[index] = k
		keys += k
		if index < colLength-1 {
			keys += ","
		}
		index++
	}
	for index := 0; index < length; index++ {
		temp := "("
		for in := 0; in < colLength; in++ {
			temp += "\"" + datas[index][colName[in]].(string) + "\""
			if in < colLength-1 {
				temp += ","
			}
		}
		temp += ")"
		if index < length-1 {
			temp += ","
		}
		vals += temp
	}
	sql := "INSERT INTO " + table + "(" + keys + ") VALUES " + vals
	return sql
}

func createUpSQL(table string, data map[string]string, where map[string]string) string {
	length := len(data)
	updateInfo, whereInfo, index := "", "", 0
	for k, v := range data {
		updateInfo += k + "=" + "\"" + v + "\""
		if index < length-1 {
			updateInfo += ","
		}
		index++
	}
	length, index = len(where), 0
	for k, v := range where {
		whereInfo += k + "\"" + v + "\""
		if index < length-1 {
			whereInfo += " AND "
		}
		index++
	}
	sql := "UPDATE " + table + " SET " + updateInfo + " WHERE " + whereInfo
	return sql
}

// GetRow 取得单条数据
func GetRow(table string, userDesc []string, where map[string]string) (map[string]interface{}, error) {
	sql := createSelectSQL(table, userDesc, where, 1)
	rows, err := db.Query(sql)
	result := make(map[string]interface{})
	if err != nil {
		common.WriteLog(dbLogPath, sql)
		return nil, err
	}
	if rows != nil {
		rows.Next()
		result = scanAllParams(rows)
		if len(result) < 1 {
			result = nil
		}
	}
	return result, err
}

// GetRows 取得多条数据
func GetRows(table string, userDesc []string, where map[string]string, limit int) (map[int]map[string]interface{}, error) {
	sql := createSelectSQL(table, userDesc, where, limit)
	rows, err := db.Query(sql)
	if err != nil {
		common.WriteLog(dbLogPath, sql)
		return nil, err
	}
	var result = make(map[int]map[string]interface{})
	index := 0
	for rows.Next() {
		result[index] = scanAllParams(rows)
		index++
	}
	return result, err
}

// InsertRow 插入单条数据
func InsertRow(table string, data map[string]string) (int64, error) {
	length := len(data)
	if length == 0 || table == "" {
		return 0, errors.New("missing args")
	}

	sql := createInRowSQL(table, data)
	stmt, err := db.Prepare(sql)
	var insertID int64
	if err != nil {
		common.WriteLog(dbLog, "prepare sql fail : "+sql)
	} else {
		rs, err := stmt.Exec()
		if err != nil {
			common.WriteLog(dbLog, "insert sql fail : "+sql)
		} else {
			insertID, err = rs.LastInsertId()
		}
	}
	return insertID, err
}

// InsertRows 插入多条数据
func InsertRows(table string, datas map[int]map[string]interface{}) (int64, error) {
	length := len(datas)
	if length == 0 || table == "" {
		return 0, errors.New("missing args")
	}
	sql := createInRowsSQL(table, datas, length)
	stmt, err := db.Prepare(sql)
	var insertID int64
	if err != nil {
		common.WriteLog(dbLog, "prepare sql fail : "+sql)
	} else {
		rs, err := stmt.Exec()
		if err != nil {
			common.WriteLog(dbLog, "insert sql fail : "+sql)
		} else {
			insertID, err = rs.LastInsertId()
		}
	}
	return insertID, err
}

// DoUpdate 执行更新命令
func DoUpdate(table string, data map[string]string, where map[string]string) (int64, error) {
	length := len(data)
	if length == 0 || table == "" {
		return 0, errors.New("missing args")
	}
	sql := createUpSQL(table, data, where)
	stmt, err := db.Prepare(sql)
	var affectedID int64
	if err != nil {
		common.WriteLog(dbLog, "prepare sql fail : "+sql)
	} else {
		rs, err := stmt.Exec()
		if err != nil {
			common.WriteLog(dbLog, "insert sql fail : "+sql)
		} else {
			affectedID, err = rs.RowsAffected()
		}
	}
	return affectedID, err
}

// DoDelete 执行删除命令
func DoDelete(table string, where map[string]string) (int64, error) {
	length, index := len(where), 0
	if length == 0 || table == "" {
		return 0, errors.New("missing args")
	}

	var affectedID int64
	sql := "DELETE FROM " + table + " WHERE "
	for k, v := range where {
		sql += k + "\"" + v + "\""
		if index < length-1 {
			sql += " AND "
		}
		index++
	}
	rs, err := db.Exec(sql)
	if err != nil {
		common.WriteLog(dbLog, "prepare sql fail : "+sql)
	} else {
		if err != nil {
			common.WriteLog(dbLog, "insert sql fail : "+sql)
		} else {
			affectedID, err = rs.RowsAffected()
		}
	}
	return affectedID, err
}
