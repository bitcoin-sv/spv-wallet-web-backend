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

// ErrorUnauthorized is thrown if authorization failed.
var ErrorUnauthorized = errors.New("unauthorized")

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

	accessKeyId, accessKey, userId, paymail, err := h.authorizeSession(session)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(SessionAccessKeyId, accessKeyId)
	c.Set(SessionAccessKey, accessKey)
	c.Set(SessionUserId, userId)
	c.Set(SessionUserPaymail, paymail)
}

func (h *AuthMiddleware) authorizeSession(s sessions.Session) (accessKeyId, accessKey, userId, paymail interface{}, err error) {
	accessKeyId = s.Get(SessionAccessKeyId)
	accessKey = s.Get(SessionAccessKey)
	userId = s.Get(SessionUserId)
	paymail = s.Get(SessionUserPaymail)

	if isNilOrEmpty(accessKeyId) ||
		isNilOrEmpty(accessKey) ||
		userId == nil ||
		paymail == nil {
		return nil, nil, nil, nil, ErrorUnauthorized
	}

	err = h.checkAccessKey(accessKey.(string), accessKeyId.(string))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%w: %w", ErrorUnauthorized, err)
	}

	return
}

func isNilOrEmpty(s interface{}) bool {
	return s == nil || s == ""
}

// checkAccessKey checks if access key is valid by getting it from BUX.
func (h *AuthMiddleware) checkAccessKey(accessKey, accessKeyId string) error {
	// Create bux client with keys from session
	buxClient, err := h.buxClientFactory.CreateWithAccessKey(accessKey)
	if err != nil {
		return err
	}

	_, err = buxClient.GetAccessKey(accessKeyId)
	if err != nil {
		return err
	}

	return nil
}
