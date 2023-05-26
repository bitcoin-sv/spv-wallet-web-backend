package endpoints

import (
	"bux-wallet/config"
	"bux-wallet/domain"
	"errors"

	"bux-wallet/transports/http/endpoints/api/users"
	router "bux-wallet/transports/http/endpoints/routes"
	httpserver "bux-wallet/transports/http/server"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// SetupWalletRoutes main point where we're registering endpoints registrars (handlers that will register endpoints in gin engine)
//
//	and middlewares. It's returning function that can be used to setup engine of httpserver.HttpServer
func SetupWalletRoutes(s *domain.Services) httpserver.GinEngineOpt {
	routes := []interface{}{
		users.NewHandler(s),
	}

	return func(engine *gin.Engine) {
		prefix := viper.GetString(config.EnvHttpServerUrlPrefix)
		rootRouter := engine.Group(prefix)
		apiRouter := engine.Group(prefix)
		for _, r := range routes {
			switch r := r.(type) {
			case router.RootEndpoints:
				r.RegisterEndpoints(rootRouter)
			case router.ApiEndpoints:
				r.RegisterApiEndpoints(apiRouter)
			default:
				panic(errors.New("unexpected router endpoints registrar"))
			}
		}
	}
}
