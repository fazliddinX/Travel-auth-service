package middleware

import (
	t "Auth-service/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

func MiddleWareAcces() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token is empty"})
			c.Abort()
			return
		}
		claims, err := t.ExtractClaimAcces(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("id", claims.Id)
		c.Set("name", claims.Name)
		c.Set("age", claims.Age)
		c.Set("email", claims.Email)
		c.Next()
	}
}

func MiddleWareRefresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token is empty"})
			c.Abort()
			return
		}
		claims, err := t.ExtractClaimAcces(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("id", claims.Id)
		c.Set("name", claims.Name)
		c.Set("age", claims.Age)
		c.Set("email", claims.Email)
		c.Next()
	}
}
