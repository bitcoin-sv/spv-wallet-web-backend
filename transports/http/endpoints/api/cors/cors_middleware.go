package cors

import (
	"net/http"

	"bux-wallet/config"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// CorsMiddleware is a middleware that handles CORS.
func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		for _, allowedOrigin := range viper.GetStringSlice(config.EnvHttpServerCorsAllowedDomains) {
			if allowedOrigin == origin {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cache-Control")
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
				break
			}
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
