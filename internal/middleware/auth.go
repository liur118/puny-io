package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liur/puny-io/internal/service"
)

func AuthMiddleware(jwtService *service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		username, err := jwtService.ValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		c.Set("username", username)
		c.Next()
	}
}
