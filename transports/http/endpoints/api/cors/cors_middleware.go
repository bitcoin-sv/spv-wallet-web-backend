package cors

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet-web-backend/config"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Middleware is a middleware that handles CORS.
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		for _, allowedOrigin := range viper.GetStringSlice(config.EnvHTTPServerCorsAllowedDomains) {
			if allowedOrigin == origin {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
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
