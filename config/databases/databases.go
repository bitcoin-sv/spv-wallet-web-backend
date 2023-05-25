package databases

import (
	"bux-wallet/config"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // nolint: golint
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

// SetUpDatabase is used to set up database connection.
func SetUpDatabase() *sql.DB {
	// Load config.
	host := viper.GetString(config.EnvDbHost)
	port := viper.GetInt(config.EnvDbPort)
	user := viper.GetString(config.EnvDbUser)
	password := viper.GetString(config.EnvDbPassword)
	dbname := viper.GetString(config.EnvDbName)
	sslMode := viper.GetString(config.EnvDbSslMode)

	// Build connection string.
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslMode)

	fmt.Println(psqlInfo)

	// Open database connection.
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	// Check database connection.
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// defer db.Close()

	fmt.Println("Successfully connected!")

	runMigration(db)

	return db
}

// runMigration is used to run database migrations.
func runMigration(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+viper.GetString(config.EnvDbMigrationsPath),
		"postgres", driver)

	if err != nil {
		panic(err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
}
