package common

import (
	"regexp"
)

// ValidSignUp 验证注册信息
func ValidSignUp(user map[string]string) bool {
	match1, _ := regexp.MatchString("^[a-zA-Z0-9]{5,16}$", user["username"])
	match2, _ := regexp.MatchString("^[a-zA-Z0-9]{6,16}$", user["password"])
	match3, _ := regexp.MatchString("[\u4e00-\u9fa5]", user["realname"])
	if !match1 || !match2 || !match3 {
		return false
	}
	return true
}
