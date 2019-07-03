package controller

import (
	"fmt"
	"goproject/src/common"
	"goproject/src/model"
	"net/http"
	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

// HallConfig 房间配置
var HallConfig = map[string]map[string]string{}

func init() {
	HallConfig["chineseChess"] = map[string]string{"tableNumbers": "35"}
}

// GetSeatList 获取房间列表
func GetSeatList(c *gin.Context) {
	info, _ := c.Get("MAP")
	params := info.(map[string]string)
	if !common.CheckParamsExist([]string{"gameCode"}, params) {
		c.JSON(http.StatusOK, gin.H{"result": false, "msg": "参数缺失"})
		return
	}
	key := params["gameCode"]
	if _, exist := HallConfig[key]; !exist {
		c.JSON(http.StatusOK, gin.H{"result": false, "msg": "游戏代号错误"})
		return
	}
	nums, _ := strconv.Atoi(HallConfig["chineseChess"]["tableNumbers"])
	var list = []map[string]string{}
	for i := 1; i <= nums; i++ {
		temp := map[string]string{}
		index := strconv.Itoa(i)
		temp["PlayerA"] = model.RedisGetHash("HALL_chineseHall_"+index, "PlayerA")
		temp["PlayerB"] = model.RedisGetHash("HALL_chineseHall_"+index, "PlayerB")
		temp["PlayerAStatus"] = model.RedisGetHash("HALL_chineseHall_"+index, "PlayerAStatus")
		temp["PlayerBStatus"] = model.RedisGetHash("HALL_chineseHall_"+index, "PlayerBStatus")
		temp["Status"] = model.RedisGetHash("HALL_chineseHall_"+index, "Status")
		list = append(list, temp)
	}
	Rds := model.RedisClient.Get()
	defer Rds.Close()
	result, _ := redis.Strings(Rds.Do("SMEMBERS", "HALL_chineseHall_User_List"))
	var users = []string{}
	for i := 0; i < len(result); i++ {
		users = append(users, result[i])
	}
	userList, _ := model.GetUserInfoInHall(users)
	c.JSON(http.StatusOK, gin.H{"result": true, "msg": "", "data": gin.H{"tables": list, "users": userList}})
}

// GetIntoHall 进入房间
func GetIntoHall(c *gin.Context) {
	userid, _ := c.Get("userid")
	Rds := model.RedisClient.Get()
	defer Rds.Close()
	result, _ := redis.Int(Rds.Do("SISMEMBER", "HALL_chineseHall_User_List", userid))
	if result == 1 {
		c.JSON(http.StatusOK, gin.H{"result": false, "msg": "已经在房间中", "data": "EXIST"})
		return
	}
	reply, err := Rds.Do("SADD", "HALL_chineseHall_User_List", userid)
	fmt.Println(reply, err)
}
