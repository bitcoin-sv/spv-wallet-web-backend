package endpoints

import (
	"bux-wallet/config"
	"bux-wallet/domain"
	"errors"

	"bux-wallet/transports/http/endpoints/registration"
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
		registration.NewHandler(s),
	}

	return func(engine *gin.Engine) {
		rootRouter := engine.Group("")
		prefix := viper.GetString(config.EnvHttpServerUrlPrefix)
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
