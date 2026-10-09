package handler

import (
	"goiot/pkg/response"

	"github.com/gin-gonic/gin"
)

// HealthCheck 健康检查接口。
func HealthCheck(c *gin.Context) {
	response.Success(c, nil)
}
