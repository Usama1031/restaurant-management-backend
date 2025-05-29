package middleware

import (
	"golang-restaurant-management/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken, err := c.Cookie("token")

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Authorization Required!"})
			c.Abort()
			return
		}

		claims, msg := helpers.ValidateToken(clientToken)

		if msg != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("First_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("user_type", claims.User_type)
		c.Set("uid", claims.Uid)

		c.Next()
	}
}
