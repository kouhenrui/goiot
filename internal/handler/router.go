package handler

import (
	"goiot/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Device *DeviceHandler
	// Product   *handler.ProductHandler
	// Telemetry *handler.TelemetryHandler
}

func New(h *Handlers) *gin.Engine {
	r := gin.New()
	r.Use(
		response.Recovery(),
		response.ErrorHandler(),
	)
	r.NoRoute(response.NoRoute())
	r.GET("/health", HealthCheck)
	rg := r.Group("/api/v1")

	rg.POST("/devices/mqtt", h.Device.RegisterMQTT)
	// for _, module := range modules {
	// 	module.Register(api)
	// }

	return r
}
