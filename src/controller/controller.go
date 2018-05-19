package controller

import (
	"model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetMaintenance 得到维护情况
func GetMaintenance(c *gin.Context) {
	result := model.GetMaintenance()
	c.JSON(http.StatusOK, result)
}
