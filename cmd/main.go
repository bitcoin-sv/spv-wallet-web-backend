package main

import (
	"bux-wallet/config"
	"bux-wallet/pkg/handler"
	"bux-wallet/pkg/logger"
	"bux-wallet/pkg/repository"
	"bux-wallet/pkg/server"
	"bux-wallet/pkg/service"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
)

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	log := logger.NewAppLogger()
	cfg, err := config.Init()
	if err != nil {
		log.Err(err).Msg("Invalid config")
		return
	}
	pool, err := pgxpool.New(context.Background(), cfg.DB.ConnectionString)
	if err != nil {
		log.Err(err).Msg("Unable to connect to database")
		os.Exit(1)
	}

	defer pool.Close()

	repositories := repository.NewRepositories(pool)

	services := service.NewServices(repositories)
	handlers := handler.NewHandler(services, &cfg.HTTP)

	srv := server.NewServer(cfg.HTTP.Port, handlers.Init())

	go func() {
		if err := srv.Run(); err != nil {
			log.Err(err).Msg("error occurred while running http server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info().Msg("Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Err(err).Msg("error occured on server shutting down")
	}
}
