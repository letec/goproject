package controller

import (
	"encoding/json"
	"model"
	"net/http"
)

// GetMaintenance 得到维护情况
func GetMaintenance(w http.ResponseWriter) bool {
	result := model.GetMaintenance()
	b, _ := json.Marshal(result)
	w.Write(b)
	return true
}
