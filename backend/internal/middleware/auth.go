package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gomall-lite-api/config"
	"gomall-lite-api/internal/dto"
	"gomall-lite-api/internal/logger"
)

func Auth(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Default().Warn("auth failed: missing bearer token", "request_id", GetRequestID(c), "path", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, dto.Fail(401, "请先登录"))
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			logger.Default().Warn("auth failed: invalid token", "request_id", GetRequestID(c), "path", c.Request.URL.Path, "error", err)
			c.JSON(http.StatusUnauthorized, dto.Fail(401, "登录已过期，请重新登录"))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.Default().Warn("auth failed: invalid claims", "request_id", GetRequestID(c), "path", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, dto.Fail(401, "token 无效"))
			c.Abort()
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			logger.Default().Warn("auth failed: user id missing", "request_id", GetRequestID(c), "path", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, dto.Fail(401, "token 无效"))
			c.Abort()
			return
		}

		c.Set("userID", uint(userIDFloat))
		c.Next()
	}
}
