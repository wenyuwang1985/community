package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wenyuwang1985/community/pkg/jwt"
)

func AuthRequired(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, Response{Code: 401, Msg: "缺少认证信息"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, Response{Code: 401, Msg: "认证格式错误"})
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(parts[1], secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, Response{Code: 401, Msg: "Token 无效或已过期"})
			c.Abort()
			return
		}

		if claims.TokenType != "access" {
			c.JSON(http.StatusUnauthorized, Response{Code: 401, Msg: "Token 类型错误"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
