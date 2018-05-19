package router

import (
	"controller"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InitRouter 绑定所有路由
func InitRouter(port string) {
	router := gin.New()
	router.GET("/signin", controller.GetMaintenance)        // 登陆
	router.POST("/signin", JSONParams(), controller.SignIn) // 登陆
	router.POST("/signup", JSONParams(), controller.SignUp) // 注册
	router.Run(port)
}

// JSONParams 取得JSON参数
func JSONParams() gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf []byte
		c.Request.Body.Read(buf)
		info := make(map[string]string)
		err := json.Unmarshal(buf, &info)
		fmt.Println(info)
		if err != nil {
			c.Set("isNext", false)
			c.JSON(http.StatusOK, gin.H{"code": "10001", "msg": "参数结构错误"})
		} else {
			c.Set("isNext", true)
			c.Set("MAP", info)
		}
	}
}
