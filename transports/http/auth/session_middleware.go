package auth

import (
	"bux-wallet/config"
	router "bux-wallet/transports/http/endpoints/routes"
	"database/sql"
	"net/http"

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
)

// NewSessionMiddleware create Session middleware that is retrieving auth token from cookie.
func NewSessionMiddleware(db *sql.DB, engine *gin.Engine) router.ApiMiddlewareFunc {
	store, err := postgres.NewStore(db, []byte("secret"))
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
