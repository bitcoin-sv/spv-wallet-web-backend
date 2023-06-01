package auth

import (
	"bux-wallet/config"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	authorizationHeader = "Authorization"
)

// TokenMiddleware middleware that is retrieving token from Authorization header.
type TokenMiddleware struct{}

// NewTokenMiddleware create Token middleware that is retrieving token from Authorization header.
func NewTokenMiddleware() *TokenMiddleware {
	return &TokenMiddleware{}
}

// ApplyToApi is a middleware which checks if the request has a valid token.
func (h *TokenMiddleware) ApplyToApi(c *gin.Context) {
	rawToken, err := h.parseAuthHeader(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		return
	}

	err = h.getToken(rawToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		return
	}
}

func (h *TokenMiddleware) parseAuthHeader(c *gin.Context) (string, error) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	return headerParts[1], nil
}

func (h *TokenMiddleware) getToken(token string) error {
	adminToken := viper.GetString(config.EnvHttpServerAuthToken)

	if token != adminToken {
		return errors.New("invalid access token")
	}

	return nil
}
