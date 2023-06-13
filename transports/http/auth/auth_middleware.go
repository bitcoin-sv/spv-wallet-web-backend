package auth

import (
	"bux-wallet/domain"
	"bux-wallet/domain/users"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware middleware that is checking the variables set in session.
type AuthMiddleware struct {
	adminBuxClient   users.AdmBuxClient
	buxClientFactory users.BuxClientFactory
	services         *domain.Services
}

// NewAuthMiddleware create middleware that is checking the variables in session.
func NewAuthMiddleware(s *domain.Services) *AuthMiddleware {
	adminBuxClient, err := s.BuxClientFactory.CreateAdminBuxClient()
	if err != nil {
		panic(fmt.Errorf("error during creating admin bux client: %w", err))
	}
	return &AuthMiddleware{
		adminBuxClient:   adminBuxClient,
		buxClientFactory: s.BuxClientFactory,
		services:         s,
	}
}

// ApplyToApi is a middleware which checks if the validity of variables in session.
func (h *AuthMiddleware) ApplyToApi(c *gin.Context) {
	session := sessions.Default(c)

	// Try to retrieve session access key id.
	accessKeyId := session.Get(SessionAccessKeyId)
	if accessKeyId == nil || accessKeyId == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	// Try to retrieve session access key.
	accessKey := session.Get(SessionAccessKey)
	if accessKey == nil || accessKey == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	// Try to retrieve session user id.
	userId := session.Get(SessionUserId)
	if userId == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	err := h.checkAccessKey(accessKey.(string), accessKeyId.(string))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(SessionAccessKeyId, accessKeyId)
	c.Set(SessionAccessKey, accessKey)
	c.Set(SessionUserId, userId)
}

// checkAccessKey checks if access key is valid by getting it from BUX.
func (h *AuthMiddleware) checkAccessKey(accessKey, accessKeyId string) error {
	// Create bux client with keys from session
	buxClient, err := h.buxClientFactory.CreateWithAccessKey(accessKey)
	if err != nil {
		return fmt.Errorf("unauthorized, error during checking access key in BUX")
	}

	_, err = buxClient.GetAccessKey(accessKeyId)
	if err != nil {
		return fmt.Errorf("unauthorized")
	}

	return nil
}
