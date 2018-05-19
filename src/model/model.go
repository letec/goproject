package model

import (
	"common"
	"os"
	"strconv"
	"time"
)

const sysLogPath = "log/sysLog.log"
const dbLogPath = "log/dbLog.log"

// SysConfig 这个值储存了一些必须的配置默认值 用作在启动程序时自动补足数据库中没有的一些配置 随后被程序调用改变状态
var SysConfig = map[string]string{
	"maintenance":  "1", // 全站维护
	"ag_open":      "1", // AG游戏维护
	"loterry_open": "1", // 彩票维护
	"pocker_open":  "1", // 扑克维护
	"chess_open":   "1", // 棋类维护
}

// InitSysConfig 获取或补齐系统配置
func InitSysConfig() {
	configs := SystemConfigs()
	tmp := []string{}
	inserts := []string{}
	for i := 0; i < len(configs); i++ {
		tmp = append(tmp, configs[i]["cfgName"].(string))
	}
	for k := range SysConfig {
		if common.InSlice(k, tmp) == false {
			inserts = append(inserts, k)
		}
	}
	if len(inserts) > 0 {
		trans, err := db.Begin()
		if err != nil {
			common.WriteLog(dbLogPath, "事务开启失败,程序退出!func InitSysConfig()")
			os.Exit(0)
		}
		sql := "INSERT INTO sys_config(cfgName,cfgValue,cfgTime,cfgAdmin,cfgIP) VALUES(?,?,?,?,?)"
		stmt, _ := trans.Prepare(sql)
		for _, v := range inserts {
			_, err := stmt.Exec(v, SysConfig[v], time.Now().Unix(), "SYSTEM", "127.0.0.1")
			if err != nil {
				trans.Rollback()
				common.WriteLog(dbLogPath, sql)
				os.Exit(0)
			}
		}
		trans.Commit()
	}
}

// SystemConfigs 获取数据库内的配置信息
func SystemConfigs() map[int]map[string]interface{} {
	configList, err := GetRows("sys_config", []string{}, map[string]string{}, 0)
	if err != nil {
		errInfo := "读取sys配置表出错,程序退出!"
		common.WriteLog(sysLogPath, errInfo)
		os.Exit(0)
	}
	return configList
}

// GetMaintenance 获取维护信息
func GetMaintenance() map[string]string {
	rdi, err := RdsGetJSON("GetMaintenance")
	if err != nil && len(rdi) > 0 {
		return rdi
	}
	sql := "SELECT cfgName,cfgValue FROM sys_config WHERE cfgName IN ('maintenance','ag_open','loterry_open','pocker_open','chess_open') "
	rs, err := db.Query(sql)
	result := map[string]string{}
	if err != nil {
		return SysConfig
	}
	for rs.Next() {
		tmp := scanAllParams(rs)
		result[tmp["cfgName"].(string)] = strconv.Itoa(tmp["cfgValue"].(int))
	}
	RdsSetJSON("GetMaintenance", result, "5")
	return result
}
