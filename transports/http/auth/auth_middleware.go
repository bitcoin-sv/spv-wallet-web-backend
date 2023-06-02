package auth

import (
	"bux-wallet/domain"
	"bux-wallet/domain/users"
	buxclient "bux-wallet/transports/bux/client"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware middleware that is checking the variables set in session.
type AuthMiddleware struct {
	adminBuxClient *buxclient.AdminBuxClient
	services       *domain.Services
}

// NewAuthMiddleware create middleware that is checking the variables in session.
func NewAuthMiddleware(s *domain.Services) *AuthMiddleware {
	return &AuthMiddleware{
		adminBuxClient: s.UsersService.BuxClient,
		services:       s,
	}
}

// ApplyToApi is a middleware which checks if the validity of variables in session.
func (h *AuthMiddleware) ApplyToApi(c *gin.Context) {
	session := sessions.Default(c)

	// Try to retrieve session token.
	token := session.Get(SessionToken)
	if token == nil || token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	// Try to retrieve session user id.
	userId := session.Get(SessionUserId)
	if userId == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	err := h.checkAccessKey(token.(string))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		return
	}

	user, err := h.checkUser(userId.(int))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(SessionToken, token)
	c.Set(SessionUserId, userId)
	c.Set(SessionPaymail, user.Paymail)
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
