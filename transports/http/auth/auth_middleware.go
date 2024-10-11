package auth

import (
	"errors"
	"fmt"

	"github.com/bitcoin-sv/spv-wallet-web-backend/domain"
	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/users"
	"github.com/bitcoin-sv/spv-wallet-web-backend/spverrors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// ErrorUnauthorized is thrown if authorization failed.
var ErrorUnauthorized = errors.New("unauthorized")

// Middleware middleware that is checking the variables set in session.
type Middleware struct {
	adminWalletClient   users.AdminWalletClient
	walletClientFactory users.WalletClientFactory
	services            *domain.Services
	log                 *zerolog.Logger
}

// NewAuthMiddleware create middleware that is checking the variables in session.
func NewAuthMiddleware(s *domain.Services, logger *zerolog.Logger) *Middleware {
	adminWalletClient, err := s.WalletClientFactory.CreateAdminClient()
	if err != nil {
		panic(fmt.Errorf("error during creating adminWalletClient: %w", err))
	}
	log := logger.With().Str("service", "auth-middleware").Logger()
	return &Middleware{
		adminWalletClient:   adminWalletClient,
		walletClientFactory: s.WalletClientFactory,
		services:            s,
		log:                 &log,
	}
}

// ApplyToAPI is a middleware which checks if the validity of variables in session.
func (h *Middleware) ApplyToAPI(c *gin.Context) {
	session := sessions.Default(c)

	accessKeyID, accessKey, userID, paymail, xPriv, err := h.authorizeSession(session)
	if err != nil {
		spverrors.AbortWithErrorResponse(c, spverrors.ErrUnauthorized, h.log)
		return
	}

	c.Set(SessionAccessKeyID, accessKeyID)
	c.Set(SessionAccessKey, accessKey)
	c.Set(SessionUserID, userID)
	c.Set(SessionUserPaymail, paymail)
	c.Set(SessionXPriv, xPriv)
}

func (h *Middleware) authorizeSession(s sessions.Session) (accessKeyID, accessKey, userID, paymail, xPriv interface{}, err error) {
	accessKeyID = s.Get(SessionAccessKeyID)
	accessKey = s.Get(SessionAccessKey)
	userID = s.Get(SessionUserID)
	paymail = s.Get(SessionUserPaymail)
	xPriv = s.Get(SessionXPriv)

	if isNilOrEmpty(accessKeyID) ||
		isNilOrEmpty(accessKey) ||
		userID == nil ||
		paymail == nil {
		return nil, nil, nil, nil, nil, ErrorUnauthorized
	}

	err = h.checkAccessKey(accessKey.(string), accessKeyID.(string))
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("%w: %w", ErrorUnauthorized, err)
	}

	return
}

func isNilOrEmpty(s interface{}) bool {
	return s == nil || s == ""
}

//nolint:wrapcheck // error wrapped higher
func (h *Middleware) checkAccessKey(accessKey, accessKeyID string) error {
	userWalletClient, err := h.walletClientFactory.CreateWithAccessKey(accessKey)
	if err != nil {
		return err
	}

	_, err = userWalletClient.GetAccessKey(accessKeyID)
	return err
}
