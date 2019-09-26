package api

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

var DefaultMiddleware = func(c *gin.Context) {
	c.Next()
}

func TokenAuthMiddleware() gin.HandlerFunc {
	requiredToken := viper.GetString("metrics.api_token")
	if requiredToken == "" {
		return DefaultMiddleware
	}

	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		token = strings.Replace(token, "Bearer ", "", 1)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, "API token required")
			return
		}

		if token != requiredToken {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, "Invalid API token")
			return
		}
		c.Next()
	}
}
