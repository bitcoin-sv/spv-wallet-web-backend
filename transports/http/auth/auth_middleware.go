package auth

import (
	"bux-wallet/domain"
	"bux-wallet/domain/users"
	buxclient "bux-wallet/transports/bux/client"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware middleware that is checking the variables set in session.
type AuthMiddleware struct {
	adminBuxClient *buxclient.AdminBuxClient
	services       *domain.Services
}

// NewMiddleware create Cookie middleware that is checking the variables in session.
func NewAuthMiddleware(s *domain.Services) *AuthMiddleware {
	return &AuthMiddleware{
		adminBuxClient: s.UsersService.BuxClient,
		services:       s,
	}
}

// ApplyToApi is a middleware which checks if the validity of variables in session.
func (h *AuthMiddleware) ApplyToApi(c *gin.Context) {
	token := c.GetString(sessionToken)
	userId := c.GetInt(sessionUserId)

	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("missing auth cookie in gin context"))
		return
	}

	err := h.checkAccessKey(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		return
	}

	user, err := h.checkUser(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		return
	}

	c.Set("paymail", user.Paymail)
	c.Set("xpriv", user.Xpriv)
}

// checkAccessKey checks if access key is valid by getting it from BUX.
func (h *AuthMiddleware) checkAccessKey(token string) error {
	// TODO: access token validation
	//
	// err := h.adminBuxClient.GetAccessKey(token)
	// if err != nil {
	// 	return fmt.Errorf("error during checking access key in BUX: %w", err)
	// 	return nil
	// }
	return nil
}

// checkUser checks if user exists in database.
func (h *AuthMiddleware) checkUser(userId int) (*users.User, error) {
	user, err := h.services.UsersService.GetUserById(userId)
	if err != nil {
		return nil, fmt.Errorf("error during getting user by id: %w", err)
	}
	return user, nil
}
