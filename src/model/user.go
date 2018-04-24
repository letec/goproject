package model

// GetUserInfoByID 查询用户名
func GetUserInfoByID(userid string) (map[string]string, error) {
	userDesc := []string{"id", "username", "password", "salt", "realname", "phone", "bankCode"}
	var id int
	var username string
	var password string
	var salt int
	var realname string
	var phone string
	var bankCode string
	where := map[string]string{
		"id": "=" + userid,
	}
	rows, err := GetRow("user", userDesc, where)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	err = rows.Scan(&id, &username, &password, &salt, &realname, &phone, &bankCode)
	if err != nil {
		return nil, err
	}
	result["id"] = string(id)
	result["username"] = username
	result["salt"] = string(salt)
	result["realname"] = realname
	result["phone"] = phone
	result["bankCode"] = bankCode
	return result, nil
}
