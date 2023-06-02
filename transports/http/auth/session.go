package auth

import (
	"bux-wallet/config"
	"database/sql"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// SetupSessionStore setup session store.
func SetupSessionStore(db *sql.DB, engine *gin.Engine) error {
	store, err := postgres.NewStore(db, []byte("secret"))
	if err != nil {
		return err
	}

	// If we're running on localhost, we need to set domain to empty string.
	domain := viper.GetString(config.EnvHttpServerCookieDomain)
	if domain == "localhost" {
		domain = ""
	}

	options := sessions.Options{
		MaxAge:   1800,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Domain:   domain,
	}

	store.Options(options)
	engine.Use(sessions.Sessions("Authorization", store))

	return nil
}

// UpdateSession updates session with accessKeyId and userId.
func UpdateSession(c *gin.Context, accessKeyId string, userId int) error {
	session := sessions.Default(c)
	session.Set(sessionToken, accessKeyId)
	session.Set(sessionUserId, userId)
	err := session.Save()
	if err != nil {
		return err
	}
	c.Header("Access-Control-Allow-Credentials", "true")
	return nil
}
