package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"gomall-lite-api/internal/dto"
	"gomall-lite-api/internal/logger"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.Default().Error(
			"panic recovered",
			"request_id", GetRequestID(c),
			"error", recovered,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP(),
			"stack", string(debug.Stack()),
		)

		c.JSON(http.StatusInternalServerError, dto.Fail(500, "服务器内部错误"))
	})
}
