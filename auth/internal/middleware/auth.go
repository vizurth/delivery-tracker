package middleware

import (
	"delivery-tracker/common/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		c.Next()
	}
}

func AuthMiddlewareAdmin(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}
		claims, err := jwt.ParseToken(tokenString, []byte(secret))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		role := claims.Role
		if role != "admin" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Permission denied"})
			return
		}

		c.Next()
	}
}

func AuthMiddlewareСourier(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}
		claims, err := jwt.ParseToken(tokenString, []byte(secret))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		role := claims.Role
		if role != "courier" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Permission denied"})
			return
		}

		c.Next()
	}
}
