package main

import (
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bitcoin-sv/spv-wallet-web-backend/config"
	"github.com/bitcoin-sv/spv-wallet-web-backend/config/databases"
	db_users "github.com/bitcoin-sv/spv-wallet-web-backend/data/users"
	"github.com/bitcoin-sv/spv-wallet-web-backend/domain"
	"github.com/bitcoin-sv/spv-wallet-web-backend/logging"
	"github.com/bitcoin-sv/spv-wallet-web-backend/transports/http/endpoints"
	httpserver "github.com/bitcoin-sv/spv-wallet-web-backend/transports/http/server"
	"github.com/bitcoin-sv/spv-wallet-web-backend/transports/websocket"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const appname = "spv-wallet-web-backend"

// nolint: godot
// @title           SPV Wallet WEB Backend
// @version			1.0
// @description     This is an API for the spv-wallet-web-frontend.
func main() {
	defaultLogger := logging.GetDefaultLogger()

	// Load config.
	config.NewViperConfig(appname).
		WithDb()

	log, err := logging.CreateLogger()
	if err != nil {
		defaultLogger.Error().Msg("cannot create logger")
		os.Exit(1)
	}

	db := databases.SetUpDatabase(log)
	defer db.Close() //nolint: all

	repo := db_users.NewUsersRepository(db)

	s, err := domain.NewServices(repo, log)
	if err != nil {
		log.Error().Msgf("cannot create services because of an error: %v", err)
		os.Exit(1)
	}

	ws, err := websocket.NewServer(log, s, db)
	if err != nil {
		log.Error().Msgf("failed to init a new websocket server: %v", err)
		os.Exit(1)
	}
	err = ws.Start()
	if err != nil {
		log.Error().Msgf("failed to start websocket server: %v", err)
		os.Exit(1)
	}

	server := httpserver.NewHttpServer(viper.GetInt(config.EnvHttpServerPort), log)
	server.ApplyConfiguration(endpoints.SetupWalletRoutes(s, db, log, ws))
	server.ApplyConfiguration(ws.SetupEntrypoint)

	go startServer(server)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	if err = server.Shutdown(); err != nil {
		log.Error().Msgf("failed to stop http server: %v", err)
	}
	if err = ws.Shutdown(); err != nil {
		log.Error().Msgf("failed to stop websocket server: %v", err)
	}
}

func startServer(server *httpserver.HttpServer) {
	if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error().Msgf("cannot start server because of an error: %v", err)
		os.Exit(1)
	}
}
