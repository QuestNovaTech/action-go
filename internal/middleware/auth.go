package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"actiondelta/internal/auth"
)

// AuthMiddleware 校验请求头中的JWT，并将用户ID注入到上下文。
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 优先从 Authorization: Bearer <token> 获取
		header := c.GetHeader("Authorization")
		token := ""
		if header != "" && strings.HasPrefix(header, "Bearer ") {
			token = strings.TrimPrefix(header, "Bearer ")
		}
		// 兼容 Authentication 头，既支持 Bearer 前缀，也支持裸 token
		if token == "" {
			alt := c.GetHeader("Authentication")
			if alt != "" {
				if strings.HasPrefix(alt, "Bearer ") {
					token = strings.TrimPrefix(alt, "Bearer ")
				} else {
					token = alt
				}
			}
		}
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "missing or invalid token"})
			return
		}

		claims, err := auth.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "invalid token"})
			return
		}
		// 即将过期则自动刷新，避免刷新风暴可加入最小间隔控制（此处简化）
		if claims.ExpiresAt != nil {
			if time.Until(claims.ExpiresAt.Time) < 5*time.Minute && time.Until(claims.ExpiresAt.Time) > 0 {
				if at, rt, err := auth.GenerateTokens(claims.UserId); err == nil {
					c.Header("New-Access-Token", at)
					c.Header("New-Refresh-Token", rt)
				}
			}
		}
		c.Set("userId", claims.UserId)
		c.Next()
	}
}
