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
	rows, err := GetRow("user", userDesc, where)
	// 报错了就返回空和错误
	if err != nil {
		return nil, err
	}
	// 取得所有列
	columns, _ := rows.Columns()
	// 计算列的数量
	length := len(columns)
	// 定义一个存放值的map
	value := make([]interface{}, length)
	// 定义一个存放指针的map
	columnPointers := make([]interface{}, length)
	// 循环把变量地址拿到 &代表地址
	for i := 0; i < length; i++ {
		columnPointers[i] = &value[i]
	}
	// 这个是最经典的 把一个interface的切片作为可变参数传进去scan 有多少列都在columnPointers里面了
	rows.Scan(columnPointers...)
	// 创建结果变量存返回值
	result := make(map[string]interface{})
	// 循环把结果放到result里面
	for i := 0; i < length; i++ {
		columnName := columns[i]
		columnValue := columnPointers[i].(*interface{})
		result[columnName] = string((*columnValue).([]uint8))
	}
	return result, nil
}
