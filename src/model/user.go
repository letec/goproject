package model

// GetUserInfoByID 查询用户名
func GetUserInfoByID(userid string) (map[string]interface{}, error) {
	userDesc := []string{"id", "username", "password", "salt", "realname", "phone", "bankCode"}
	where := map[string]string{
		"id": "=" + userid,
	}
	userInfo, err := GetRow("user", userDesc, where)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

// GetUserInfoByIDs 查询用户名
func GetUserInfoByIDs(userid string) (map[int]map[string]interface{}, error) {
	userDesc := []string{"id", "username", "realname", "phone", "bankCode"}
	where := map[string]string{
		"id": "=" + userid,
	}
	userInfo, err := GetRows("user", userDesc, where, 0)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}
