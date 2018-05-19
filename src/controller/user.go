package controller

import (
	"common"
	"model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SignUp 注册用户
func SignUp(c *gin.Context) {
	params := []string{"username", "password", "realname"}
	v, b := c.Get("isNext")
	if !(b && v.(bool)) {
		return
	}
	info, b := c.Get("MAP")
	user := info.(map[string]string)
	ret := common.CheckParamsExist(params, user)
	if ret == false {
		c.JSON(http.StatusOK, gin.H{"code": "10002", "msg": "参数缺失"})
		return
	}
	ret = common.ValidSignUp(user)
	if ret == false {
		c.JSON(http.StatusOK, gin.H{"code": "10003", "msg": "参数验证不通过"})
		return
	}
	flag, err := model.CheckUserExist(user["username"])
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "10004", "msg": "网络错误"})
		return
	}
	if flag == true {
		c.JSON(http.StatusOK, gin.H{"code": "20001", "msg": "用户名已经存在"})
		return
	}
	flag, err = model.SignUpUser(user)
	if err != nil || flag == false {
		c.JSON(http.StatusOK, gin.H{"code": "10004", "msg": "网络错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "20000", "msg": "注册成功"})
	return
}

// SignIn 用户登陆
func SignIn(c *gin.Context) {
	params := []string{"username", "password"}
	v, b := c.Get("isNext")
	if !(b && v.(bool)) {
		return
	}
	info, b := c.Get("MAP")
	user := info.(map[string]string)
	ret := common.CheckParamsExist(params, user)
	status := 0
	if ret == false {
		c.JSON(http.StatusOK, gin.H{"code": "10002", "msg": "参数缺失"})
		return
	}
	userDesc := []string{"id", "username", "password", "salt", "status"}
	where := map[string]string{"username=": user["username"]}
	result, err := model.GetRow("user", userDesc, where)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": "10004", "msg": "网络错误"})
		return
	}
	if result != nil {
		status, _ = strconv.Atoi(result["status"].(string))
		if status != 0 {
			c.JSON(http.StatusOK, gin.H{"code": "10005", "msg": "账号已经被冻结,如有疑问请联系管理员!"})
			return
		}
		rpwd := result["salt"].(string) + user["password"] + result["username"].(string)
		cpwd := common.MD5(rpwd)
		if cpwd == result["password"] {
			oid, err := model.SetOid(string(result["id"].(string)))
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"code": "10004", "msg": "网络错误"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": "20000", "msg": "登陆成功", "oid": oid})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": "20002", "msg": "用户名或密码错误"})
	return
}
