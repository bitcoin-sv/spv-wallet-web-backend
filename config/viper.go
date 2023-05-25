package config

import (
	"strings"

	"github.com/spf13/viper"
)

// NewViperConfig creates and returns new viper config.
func NewViperConfig(appname string) *Config {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	setHttpServerDefaults()
	return &Config{}
}

// WithDb edits and returns database-based viper configuration.
func (c *Config) WithDb() *Config {
	viper.SetDefault(EnvDbHost, "localhost")
	viper.SetDefault(EnvDbPort, "5432")
	viper.SetDefault(EnvDbUser, "postgres")
	viper.SetDefault(EnvDbPassword, "postgres")
	viper.SetDefault(EnvDbName, "postgres")
	viper.SetDefault(EnvDbSslMode, "disable")
	viper.SetDefault(EnvDbMigrationsPath, "data/sql/migrations")

	c.Db = &Db{
		Host:           viper.GetString(EnvDbHost),
		Port:           viper.GetInt(EnvDbPort),
		User:           viper.GetString(EnvDbUser),
		Password:       viper.GetString(EnvDbPassword),
		Name:           viper.GetString(EnvDbName),
		SslMode:        viper.GetString(EnvDbSslMode),
		MigrationsPath: viper.GetString(EnvDbMigrationsPath),
	}
	return c
}

// setHttpServerDefaults sets default values for http server.
func setHttpServerDefaults() {
	viper.SetDefault(EnvHttpServerReadTimeout, 10)
	viper.SetDefault(EnvHttpServerWriteTimeout, 10)
	viper.SetDefault(EnvHttpServerPort, 8080)
	viper.SetDefault(EnvHttpServerUrlPrefix, "/api/v1")
}
