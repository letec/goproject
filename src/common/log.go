package common

import (
	"os"
	"strings"
)

// WriteLog 追加写入日志文件
func WriteLog(path string, text string) bool {
	result := false
	ret, _ := PathExists(path)
	if ret == false {
		_, err2 := os.Create(path)
		if err2 != nil {
			return false
		}
	}
	fd, err3 := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err3 != nil {
		return false
	}
	time := YmdHis()
	content := strings.Join([]string{"======", time, "=====\n", text, "\n\n"}, "")
	buf := []byte(content)
	_, err4 := fd.Write(buf)
	if err4 == nil {
		result = true
	}
	fd.Close()
	return result
}
