package router

import (
	"encoding/json"
	"goproject/src/common"
	"goproject/src/controller"
	"goproject/src/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InitRouter 绑定所有路由
func InitRouter(port string) {
	router := gin.New()
	router.POST("/maintenance", JSONParams(), controller.GetMaintenance) // 维护
	router.POST("/signin", JSONParams(), controller.SignIn)              // 登陆
	router.POST("/signup", JSONParams(), controller.SignUp)              // 注册
	router.Run(port)
}

// JSONParams 取得JSON参数
func JSONParams() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf := make([]byte, 1024)
		n, _ := c.Request.Body.Read(buf)
		info := make(map[string]string)
		err := json.Unmarshal(buf[0:n], &info)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"result": false, "msg": "参数结构错误"})
		} else {
			c.Set("MAP", info)
		}
		// 在线检查
		if !common.InSlice(c.Request.URL.Path, []string{"/signin", "/signup"}) {
			userid, err := model.CheckOnline(info["oid"])
			c.Set("userid", userid)
			if err != nil && userid == "" {
				c.AbortWithStatusJSON(http.StatusOK, gin.H{"result": false, "msg": "您已经被登出"})
			}
		}
	}
}
