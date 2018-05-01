package common

import "regexp"

// ValidSignUp 验证注册信息
func ValidSignUp(user map[string]string) bool {
	match1, _ := regexp.MatchString("p([a-z]+)ch", user["username"])
	match2, _ := regexp.MatchString("p([a-z]+)ch", user["password"])
	match3, _ := regexp.MatchString("p([a-z]+)ch", user["realname"])
	if !match1 || !match2 || !match3 {
		return false
	}
	return true
}
