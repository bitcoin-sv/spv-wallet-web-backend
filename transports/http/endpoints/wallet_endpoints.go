package endpoints

import (
	"database/sql"
	"errors"

	"bux-wallet/domain"

	"bux-wallet/transports/http/auth"
	"bux-wallet/transports/http/endpoints/api/access"
	"bux-wallet/transports/http/endpoints/api/users"
	router "bux-wallet/transports/http/endpoints/routes"
	httpserver "bux-wallet/transports/http/server"

	"github.com/gin-gonic/gin"
)

// SetupWalletRoutes main point where we're registering endpoints registrars (handlers that will register endpoints in gin engine)
//
//	and middlewares. It's returning function that can be used to setup engine of httpserver.HttpServer
func SetupWalletRoutes(s *domain.Services, db *sql.DB) httpserver.GinEngineOpt {
	accessRootEndpoints, accessApiEndpoints := access.NewHandler(s)
	usersRootEndpoints, usersApiEndpoints := users.NewHandler(s)

	routes := []interface{}{
		usersRootEndpoints,
		usersApiEndpoints,
		accessRootEndpoints,
		accessApiEndpoints,
	}

	// rootMiddlewares := toHandlers(auth.NewTokenMiddleware())
	apiMiddlewares := toHandlers(auth.NewSessionMiddleware(), auth.NewAuthMiddleware(s))

	return func(engine *gin.Engine) {
		// Setup session middleware.
		err := auth.SetupSessionStore(db, engine)
		if err != nil {
			panic(err)
		}

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
