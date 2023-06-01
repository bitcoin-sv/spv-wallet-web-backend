package auth

import (
	"bux-wallet/config"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// SetupSessionStore setup session store.
func SetupSessionStore(store sessions.Store, engine *gin.Engine) {
	options := sessions.Options{
		MaxAge:   1800,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	// If we're running on localhost, we need to set domain to empty string.
	domain := viper.GetString(config.EnvHttpServerCookieDomain)
	if domain == "localhost" {
		domain = ""
	}

	store.Options(options)
	engine.Use(sessions.Sessions("Authorization", store))

}

// GetSessionToken updates session with accessKeyId and userId
func UpdateSession(c *gin.Context, accessKeyId string, userId int) {
	session := sessions.Default(c)
	session.Set(sessionToken, accessKeyId)
	session.Set(sessionUserId, userId)
	session.Save()
	c.Header("Access-Control-Allow-Credentials", "true")
}
