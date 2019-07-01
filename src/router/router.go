package router

import (
	"encoding/json"
	"fmt"
	"goproject/src/common"
	"goproject/src/controller"
	"goproject/src/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// InitRouter 绑定所有路由
func InitRouter(port string) {
	router := gin.New()
	router.Use(Cors())
	router.GET("/captcha", controller.CreateCaptcha)                     // 验证码
	router.POST("/maintenance", JSONParams(), controller.GetMaintenance) // 维护
	router.POST("/signin", JSONParams(), controller.SignIn)              // 登陆
	router.POST("/signup", JSONParams(), controller.SignUp)              // 注册
	router.POST("/hall", JSONParams(), controller.GetSeatList)           // 注册
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

// Cors 跨域设置
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			// 允许访问所有域
			c.Header("Access-Control-Allow-Origin", "*")
			//服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			// 允许跨域设置 可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}
