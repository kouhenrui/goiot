package response

import "github.com/gin-gonic/gin"

// Response 统一 HTTP 返回结构。
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功返回。
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// SuccessMessage 自定义成功消息。
func SuccessMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(200, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// Error 错误返回。
func Error(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// 交给全局错误处理中间件处理
	_ = c.Error(err)
}
