package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

// NewViperConfig creates and returns new viper config.
func NewViperConfig() *Config {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	setHTTPServerDefaults()
	setSpvWalletDefaults()
	setHashDefaults()
	setLoggingDefaults()
	setEndpointsDefaults()
	setWebsocketDefaults()
	setContactsDefaults()
	setCacheDefaults()
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

// setHTTPServerDefaults sets default values for http server.
func setHTTPServerDefaults() {
	viper.SetDefault(EnvHTTPServerReadTimeout, 40)
	viper.SetDefault(EnvHTTPServerWriteTimeout, 40)
	viper.SetDefault(EnvHTTPServerPort, 8180)
	viper.SetDefault(EnvHTTPServerCookieDomain, "localhost")
	viper.SetDefault(EnvHTTPServerCookieSecure, false)
	viper.SetDefault(EnvHTTPServerCorsAllowedDomains, []string{})
	viper.SetDefault(EnvHTTPServerSessionSecret, "secret")
}

// setSpvWalletDefaults sets default values for spv-wallet connection.
func setSpvWalletDefaults() {
	viper.SetDefault(EnvAdminXpriv, "xprv9s21ZrQH143K3CbJXirfrtpLvhT3Vgusdo8coBritQ3rcS7Jy7sxWhatuxG5h2y1Cqj8FKmPp69536gmjYRpfga2MJdsGyBsnB12E19CESK")
	viper.SetDefault(EnvServerURL, "http://localhost:3003/v1")
	viper.SetDefault(EnvPaymailDomain, "example.com")
	viper.SetDefault(EnvPaymailAvatar, "http://localhost:3003/static/paymail/avatar.jpg")
}

// setHashDefaults sets default values for hash.
func setHashDefaults() {
	viper.SetDefault(EnvHashSalt, "spv-wallet")
}

func setLoggingDefaults() {
	viper.SetDefault(EnvLoggingLevel, "Debug")
	viper.SetDefault(EnvLoggingInstanceName, "spv-wallet-web-backend")
	viper.SetDefault(EnvLoggingFormat, "console")
	viper.SetDefault(EnvLoggingLogOrigin, false)
}

// setEndpointsDefaults sets default values for endpoints used in app.
func setEndpointsDefaults() {
	viper.SetDefault(EnvEndpointsExchangeRate, "https://api.whatsonchain.com/v1/bsv/main/exchangerate")
}

// setWebhookDefaults sets default values for websocket.
func setWebsocketDefaults() {
	viper.SetDefault(EnvWebsocketHistoryMax, 300)
	viper.SetDefault(EnvWebsocketHistoryTTL, 10)
}

func setContactsDefaults() {
	viper.SetDefault(EnvContactsPasscodePeriod, uint(3600)) // 1h
	viper.SetDefault(EnvContactsPasscodeDigits, uint(2))
}

// setCacheDefaults sets default values for cache.
func setCacheDefaults() {
	viper.SetDefault(EnvCacheSettingsTTL, 60*time.Second)
}
