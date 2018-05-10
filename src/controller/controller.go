package controller

import (
	"model"
)

// GetMaintenance 得到维护情况
func GetMaintenance() (map[string]string, error) {
	return model.GetMaintenance()
}
