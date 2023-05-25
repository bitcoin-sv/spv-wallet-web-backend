package main

import (
	"bux-wallet/config"
	"bux-wallet/config/databases"
	db_users "bux-wallet/data/users"
	"bux-wallet/domain"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"bux-wallet/transports/http/endpoints"
	httpserver "bux-wallet/transports/http/server"

	"github.com/spf13/viper"
)

const appname = "bux-wallet-backend"

func main() {
	// Load config.
	c := config.NewViperConfig(appname).
		WithDb()

	fmt.Println(c)
	fmt.Println(c.Db.Host)

	db := databases.SetUpDatabase()
	defer db.Close() // nolint: all


	repo := db_users.NewUsersRepository(db)
	s := domain.NewServices(repo)

	server := httpserver.NewHttpServer(viper.GetInt(config.EnvHttpServerPort))
	server.ApplyConfiguration(endpoints.SetupWalletRoutes(s))

	go func() {
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("cannot start server because of an error: ", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	if err := server.Shutdown(); err != nil {
		fmt.Println("failed to stop http server: ", err)
	}
}
