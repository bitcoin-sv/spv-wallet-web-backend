package endpoints

import (
	"database/sql"
	"errors"
	"web-backend/domain"
	"web-backend/transports/http/endpoints/status"
	"web-backend/transports/http/endpoints/swagger"
	"web-backend/transports/websocket"

	"github.com/rs/zerolog"

	"web-backend/transports/http/auth"
	"web-backend/transports/http/endpoints/api/access"
	"web-backend/transports/http/endpoints/api/transactions"
	"web-backend/transports/http/endpoints/api/users"
	router "web-backend/transports/http/endpoints/routes"
	httpserver "web-backend/transports/http/server"

	"github.com/gin-gonic/gin"
)

// SetupWalletRoutes main point where we're registering endpoints registrars (handlers that will register endpoints in gin engine)
//
//	and middlewares. It's returning function that can be used to setup engine of httpserver.HttpServer
func SetupWalletRoutes(s *domain.Services, db *sql.DB, log *zerolog.Logger, ws websocket.Server) httpserver.GinEngineOpt {
	accessRootEndpoints, accessApiEndpoints := access.NewHandler(s, log)
	usersRootEndpoints, usersApiEndpoints := users.NewHandler(s, log)

	routes := []interface{}{
		swagger.NewHandler(),
		status.NewHandler(),
		usersRootEndpoints,
		usersApiEndpoints,
		accessRootEndpoints,
		accessApiEndpoints,
		transactions.NewHandler(s, log, ws),
	}

	return func(engine *gin.Engine) {
		apiMiddlewares := router.ToHandlers(
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
