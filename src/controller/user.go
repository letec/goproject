package controller

import (
	"goproject/src/common"
	"goproject/src/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SignUp 注册用户
func SignUp(c *gin.Context) {
	params := []string{"username", "password", "repassword"}
	info, _ := c.Get("MAP")
	user := info.(map[string]string)
	ret := common.CheckParamsExist(params, user)
	if ret == false {
		c.JSON(http.StatusOK, gin.H{"result": false, "msg": "参数缺失"})

	}
	ret = common.ValidSignUp(user)
	if ret == false {
		c.JSON(http.StatusOK, gin.H{"result": false, "msg": "参数验证不通过"})
		return
	}
	if user["password"] != user["repassword"] {
		c.JSON(http.StatusOK, gin.H{"result": false, "msg": "两次输入的密码不一致"})
		return
	}
	flag, err := model.CheckUserExist(user["username"])
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"result": false, "msg": "网络错误"})
		return
	}
	if flag == true {
		c.JSON(http.StatusOK, gin.H{"result": false, "msg": "用户名已经存在"})
		return
	}
	flag, err = model.SignUpUser(user)
	if err != nil || flag == false {
		c.JSON(http.StatusOK, gin.H{"result": false, "msg": "网络错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": true, "msg": "注册成功"})
	return
}

// SignIn 用户登陆
func SignIn(c *gin.Context) {
	params := []string{"username", "password"}
	info, _ := c.Get("MAP")
	user := info.(map[string]string)
	ret := common.CheckParamsExist(params, user)
	status := 0
	if ret == false {
		c.JSON(http.StatusOK, gin.H{"result": false, "msg": "参数缺失"})
		return
	}
	where := map[string]string{"UserName=": user["username"]}
	result, err := model.GetRow("user", []string{}, where)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"result": false, "msg": "网络错误" + err.Error()})
		return
	}
	if result != nil {
		status, _ = strconv.Atoi(result["Status"].(string))
		if status != 0 {
			c.JSON(http.StatusOK, gin.H{"result": false, "msg": "账号已经被冻结,如有疑问请联系管理员!"})
			return
		}
		rpwd := result["Salt"].(string) + user["password"] + result["UserName"].(string)
		cpwd := common.MD5(rpwd)
		if cpwd == result["PassWord"] {
			oid, err := model.SetOid(string(result["id"].(string)))
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"result": false, "msg": "网络错误"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"result": true, "msg": "登陆成功", "oid": oid})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"result": false, "msg": "用户名或密码错误"})
	return
}
