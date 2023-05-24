// Package config provides configuration for the application.
package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config contains all the configuration for the application.
type (
	// Config contains all the base configuration for the application.
	Config struct {
		HTTP     HTTPConfig
		DB       DBConfig
		HashSalt string `envconfig:"HASH_SALT" required:"true"`
	}

	// HTTPConfig contains configuration for the HTTP server.
	HTTPConfig struct {
		Port       int `envconfig:"PORT" required:"true"`
		Swagger    SwaggerConfig
		Middleware MiddlewareConfig
	}

	// MongoConfig contains configuration for the MongoDB connection.
	DBConfig struct {
		ConnectionString string `envconfig:"DB_CONNECTION_STRING" required:"true"`
	}

	// MiddlewareConfig contains configuration for the middleware.
	MiddlewareConfig struct {
		CORSAllowedDomains []string `envconfig:"CORS_ALLOW_DOMAINS" required:"true"`
		AuthToken          string   `envconfig:"AUTH_TOKEN" required:"true"`
	}

	// SwaggerConfig contains configuration for swagger.
	SwaggerConfig struct {
		Login    string `envconfig:"SWAGGER_LOGIN" required:"true"`
		Password string `envconfig:"SWAGGER_PASSWORD" required:"true"`
	}
)

// Init initializes the configuration.
func Init() (*Config, error) {
	var c Config
	err := envconfig.Process("", &c)
	return &c, err
}
