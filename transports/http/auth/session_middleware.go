package auth

import (
	"database/sql"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet-web-backend/config"
	router "github.com/bitcoin-sv/spv-wallet-web-backend/transports/http/endpoints/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Session variables.
const (
	SessionAccessKeyId = "accessKeyId"
	SessionAccessKey   = "accessKey"
	SessionUserId      = "userId"
	SessionUserPaymail = "paymail"
	SessionXPriv       = "xPriv"
)

// NewSessionMiddleware create Session middleware that is retrieving auth token from cookie.
func NewSessionMiddleware(db *sql.DB, engine *gin.Engine) router.ApiMiddlewareFunc {
	secret := viper.GetString(config.EnvHttpServerSessionSecret)
	store, err := postgres.NewStore(db, []byte(secret))
	if err != nil {
		panic(err)
	}

	// If we're running on localhost, we need to set domain to empty string.
	domain := viper.GetString(config.EnvHttpServerCookieDomain)
	if domain == "localhost" {
		domain = ""
	}

	secure := viper.GetBool(config.EnvHttpServerCookieSecure)

	options := sessions.Options{
		MaxAge:   1800,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
		Domain:   domain,
	}

	store.Options(options)
	engine.Use(sessions.Sessions("Authorization", store))

	return router.ApiMiddlewareFunc(sessions.Sessions("Authorization", store))
}
