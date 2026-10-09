package response

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	nanoid "github.com/matoous/go-nanoid/v2"
)

const (
	TraceIDKey    = "traceID"
	TraceIDHeader = "X-Trace-ID"
)

// ErrorHandler 全局错误处理中间件。
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 优先读取客户端传入的 TraceID

		traceID := c.GetHeader(TraceIDHeader)
		// 2. 没有则生成
		if traceID == "" {
			traceID, _ = nanoid.New(8)
		}
		// 3. 保存到 Gin Context
		c.Set(TraceIDKey, traceID)
		// 4. 写入响应 Header
		c.Header(TraceIDHeader, traceID)
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		handleError(c, err)
	}
}

func handleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var bizErr *BizError

	// 业务错误
	if errors.As(err, &bizErr) {
		if !c.Writer.Written() {
			c.AbortWithStatusJSON(
				bizErr.HTTPStatus,
				Response{
					Code:    bizErr.Code,
					Message: bizErr.Message,
				},
			)
		}

		return
	}

	// 系统未知错误
	log.Printf(
		"internal server error: %v",
		err,
	)

	if !c.Writer.Written() {
		c.AbortWithStatusJSON(
			500,
			Response{
				Code:    CodeInternalError,
				Message: "internal server error",
			},
		)
	}
}

// Recovery 捕获 panic。
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf(
					"panic recovered: %v",
					err,
				)

				if !c.Writer.Written() {
					c.AbortWithStatusJSON(
						500,
						Response{
							Code:    CodeInternalError,
							Message: "internal server error",
						},
					)
				}
			}
		}()

		c.Next()
	}
}

// NoRoute 统一处理 404。
func NoRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatusJSON(
			404,
			Response{
				Code:    CodeNotFound,
				Message: "route not found",
			},
		)
	}
}
