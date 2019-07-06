package controller

import (
	"goproject/src/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

// AllGame 所有游戏
var AllGame = []string{
	"ChineseChess",
}

// GetMaintenance 得到维护情况
func GetMaintenance(c *gin.Context) {
	result := model.GetMaintenance()
	c.JSON(http.StatusOK, result)
}

// CreateCaptcha 生成验证码方法
func CreateCaptcha(c *gin.Context) {
	//config struct for Character
	//字符,公式,验证码配置
	var configC = base64Captcha.ConfigCharacter{
		Height: 42,
		Width:  120,
		// CaptchaModeNumber:数字,CaptchaModeAlphabet:字母,CaptchaModeArithmetic:算术,CaptchaModeNumberAlphabet:数字字母混合.
		Mode:               base64Captcha.CaptchaModeNumber,
		ComplexOfNoiseText: base64Captcha.CaptchaComplexLower,
		ComplexOfNoiseDot:  base64Captcha.CaptchaComplexLower,
		IsShowHollowLine:   false,
		IsShowNoiseDot:     false,
		IsShowNoiseText:    false,
		IsShowSlimeLine:    false,
		IsShowSineLine:     false,
		CaptchaLen:         4,
	}
	idKeyC, capC := base64Captcha.GenerateCaptcha("", configC)
	base64stringC := base64Captcha.CaptchaWriteToBase64Encoding(capC)
	c.JSON(http.StatusOK, gin.H{"result": true, "msg": idKeyC, "img": base64stringC})
}
