package main

import (
	"bux-wallet/config"
	"bux-wallet/config/databases"
	db_users "bux-wallet/data/users"
	"bux-wallet/domain"
	"bux-wallet/logging"
	"bux-wallet/transports/http/endpoints"
	httpserver "bux-wallet/transports/http/server"
	"bux-wallet/transports/websocket"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
)

const appname = "bux-wallet-backend"

// nolint: godot
// @title           Bux Wallet API
// @version			1.0
// @description     This is an API for bux wallet.
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
	defer db.Close() // nolint: all

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

	go func() {
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Msgf("cannot start server because of an error: %v", err)
			os.Exit(1)
		}
	}()

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
