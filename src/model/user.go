package model

// GetUserInfoByID 查询用户名
func GetUserInfoByID(userid string) (map[string]interface{}, error) {
	// 定义一个(切片)类型,类似于PHP一维数组
	userDesc := []string{"id", "username", "password", "salt", "realname", "phone", "bankCode"}
	// WHERE条件写进map
	where := map[string]string{
		"id": "=" + userid,
	}
	// 这个方法我封装的 拼接SQL查询到结果
	userInfo, err := GetRow("user", userDesc, where)
	// 报错了就返回空和错误
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

// GetUserInfoByIDs 查询用户名
func GetUserInfoByIDs(userid string) (map[int]map[string]interface{}, error) {
	// 定义一个(切片)类型,类似于PHP一维数组
	userDesc := []string{"id", "username", "realname", "phone", "bankCode"}
	// WHERE条件写进map
	where := map[string]string{
		"id": "=" + userid,
	}
	// 这个方法我封装的 拼接SQL查询到结果
	userInfo, err := GetRows("user", userDesc, where, 0)
	// 报错了就返回空和错误
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}
