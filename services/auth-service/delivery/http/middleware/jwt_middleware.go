package middleware

import (
	"net/http"
	"strings"
	"github.com/LBRT87/GolangBackend/services/auth-service/pkg/jwt"
	"github.com/LBRT87/GolangBackend/services/auth-service/pkg/response"
	"github.com/gin-gonic/gin"
)

func tokenFromRequest(c *gin.Context) string {
	if cookieToken, err := c.Cookie("access_token"); err == nil && cookieToken != "" {
		return cookieToken
	}
	authHeader := c.GetHeader("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}
	return ""
}

func JWTAuth(jwtMgr *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := tokenFromRequest(c)
		if tokenStr == "" {
			response.Error(c, http.StatusUnauthorized, "token not found")
			c.Abort()
			return
		}

		claims, err := jwtMgr.Verify(tokenStr)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "token not valid or expired")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
