// Package handler provides HTTP handlers.
package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
)

func (h *Handler) tokenIdentity(c *gin.Context) {
	err := h.parseAuthHeader(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
	}
}

func (h *Handler) parseAuthHeader(c *gin.Context) error {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		return errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return errors.New("invalid auth header")
	}

	if headerParts[1] != h.cfg.Middleware.AuthToken {
		return errors.New("invalid auth token")
	}

	return nil
}
