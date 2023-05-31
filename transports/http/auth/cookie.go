package auth

import (
	"bux-wallet/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func SetAuthCookie(c *gin.Context, accessKeyId string) {
	domain := viper.GetString(config.EnvHttpServerCookieDomain)

	// If domain is localhost, set it to empty string.
	if domain == "localhost" {
		domain = ""
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("Authorization", accessKeyId, 1800, "/", domain, true, true)
	c.Header("Access-Control-Allow-Credentials", "true")
}
