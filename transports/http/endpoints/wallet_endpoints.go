package endpoints

import (
	"database/sql"
	"errors"

	"bux-wallet/domain"
	"bux-wallet/logging"
	"bux-wallet/transports/http/endpoints/status"
	"bux-wallet/transports/http/endpoints/swagger"

	"bux-wallet/transports/http/auth"
	"bux-wallet/transports/http/endpoints/api/access"
	"bux-wallet/transports/http/endpoints/api/transactions"
	"bux-wallet/transports/http/endpoints/api/users"
	router "bux-wallet/transports/http/endpoints/routes"
	httpserver "bux-wallet/transports/http/server"

	"github.com/gin-gonic/gin"
)

// SetupWalletRoutes main point where we're registering endpoints registrars (handlers that will register endpoints in gin engine)
//
//	and middlewares. It's returning function that can be used to setup engine of httpserver.HttpServer
func SetupWalletRoutes(s *domain.Services, db *sql.DB, lf logging.LoggerFactory) httpserver.GinEngineOpt {
	accessRootEndpoints, accessApiEndpoints := access.NewHandler(s, lf)
	usersRootEndpoints, usersApiEndpoints := users.NewHandler(s, lf)

	routes := []interface{}{
		swagger.NewHandler(),
		status.NewHandler(),
		usersRootEndpoints,
		usersApiEndpoints,
		accessRootEndpoints,
		accessApiEndpoints,
		transactions.NewHandler(s, lf),
	}

	return func(engine *gin.Engine) {
		apiMiddlewares := toHandlers(
			auth.NewSessionMiddleware(db, engine),
			auth.NewAuthMiddleware(s),
		)

		rootRouter := engine.Group("")
		apiRouter := engine.Group("/api/v1", apiMiddlewares...)
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

func toHandlers(middlewares ...router.ApiMiddleware) []gin.HandlerFunc {
	result := make([]gin.HandlerFunc, 0)
	for _, m := range middlewares {
		result = append(result, m.ApplyToApi)
	}
	return result
}
