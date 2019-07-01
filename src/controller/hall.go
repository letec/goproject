package controller

import (
	"goproject/src/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

var hallConfig = map[string]map[string]string{}

func init() {
	hallConfig["chineseChess"] = map[string]string{"tableNumbers": "35"}
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
	if _, exist := hallConfig[key]; !exist {
		c.JSON(http.StatusOK, gin.H{"result": false, "msg": "游戏代号错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": true, "msg": "", "data": hallConfig[key]})
}
