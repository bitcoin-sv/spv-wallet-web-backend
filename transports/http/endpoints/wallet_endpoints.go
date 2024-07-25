package endpoints

import (
	"database/sql"
	"errors"

	"github.com/bitcoin-sv/spv-wallet-web-backend/domain"
	"github.com/bitcoin-sv/spv-wallet-web-backend/transports/http/auth"
	"github.com/bitcoin-sv/spv-wallet-web-backend/transports/http/endpoints/api/access"
	"github.com/bitcoin-sv/spv-wallet-web-backend/transports/http/endpoints/api/config"
	"github.com/bitcoin-sv/spv-wallet-web-backend/transports/http/endpoints/api/contacts"
	"github.com/bitcoin-sv/spv-wallet-web-backend/transports/http/endpoints/api/transactions"
	"github.com/bitcoin-sv/spv-wallet-web-backend/transports/http/endpoints/api/users"
	router "github.com/bitcoin-sv/spv-wallet-web-backend/transports/http/endpoints/routes"
	"github.com/bitcoin-sv/spv-wallet-web-backend/transports/http/endpoints/status"
	"github.com/bitcoin-sv/spv-wallet-web-backend/transports/http/endpoints/swagger"
	httpserver "github.com/bitcoin-sv/spv-wallet-web-backend/transports/http/server"
	"github.com/bitcoin-sv/spv-wallet-web-backend/transports/websocket"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// SetupWalletRoutes main point where we're registering endpoints registrars (handlers that will register endpoints in gin engine)
//
//	and middlewares. It's returning function that can be used to setup engine of httpserver.HTTPServer
func SetupWalletRoutes(s *domain.Services, db *sql.DB, log *zerolog.Logger, ws websocket.Server) httpserver.GinEngineOpt {
	accessRootEndpoints, accessAPIEndpoints := access.NewHandler(s, log)
	usersRootEndpoints, usersAPIEndpoints := users.NewHandler(s, log)

	routes := []interface{}{
		swagger.NewHandler(),
		status.NewHandler(),
		config.NewHandler(s, log),
		usersRootEndpoints,
		usersAPIEndpoints,
		accessRootEndpoints,
		accessAPIEndpoints,
		transactions.NewHandler(s, log, ws),
		contacts.NewHandler(s, log),
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
			case router.APIEndpoints:
				r.RegisterAPIEndpoints(apiRouter)
			default:
				panic(errors.New("unexpected router endpoints registrar"))
			}
		}
	}
}
