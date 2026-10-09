package handler

import (
	"goiot/internal/app/device"
	"goiot/pkg/response"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	auth *device.AuthService
}

func NewDeviceHandler(auth *device.AuthService) *DeviceHandler {
	return &DeviceHandler{auth: auth}
}

type RegisterMQTTDeviceRequest struct {
	ProductID uint   `json:"product_id" binding:"required"`
	Name      string `json:"name" binding:"required"`
}

func (h *DeviceHandler) RegisterMQTT(c *gin.Context) {
	var req RegisterMQTTDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.auth.RegisterMQTTDevice(
		c.Request.Context(),
		req.ProductID,
		req.Name,
	)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}
