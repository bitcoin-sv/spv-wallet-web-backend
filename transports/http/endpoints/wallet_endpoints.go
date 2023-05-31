package endpoints

import (
	"errors"

	"bux-wallet/domain"
	"bux-wallet/transports/http/endpoints/status"
	"bux-wallet/transports/http/endpoints/swagger"

	"bux-wallet/transports/http/endpoints/api/access"
	"bux-wallet/transports/http/endpoints/api/users"
	router "bux-wallet/transports/http/endpoints/routes"
	httpserver "bux-wallet/transports/http/server"

	"github.com/gin-gonic/gin"
)

// SetupWalletRoutes main point where we're registering endpoints registrars (handlers that will register endpoints in gin engine)
//
//	and middlewares. It's returning function that can be used to setup engine of httpserver.HttpServer
func SetupWalletRoutes(s *domain.Services) httpserver.GinEngineOpt {
	routes := []interface{}{
		swagger.NewHandler(),
		status.NewHandler(),
		users.NewHandler(s),
		access.NewHandler(s),
	}

	return func(engine *gin.Engine) {
		rootRouter := engine.Group("")
		apiRouter := engine.Group("/api/v1")
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
