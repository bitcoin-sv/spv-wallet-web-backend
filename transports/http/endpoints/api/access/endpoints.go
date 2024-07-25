package access

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet-web-backend/domain"
	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/users"
	"github.com/bitcoin-sv/spv-wallet-web-backend/spverrors"
	"github.com/bitcoin-sv/spv-wallet-web-backend/transports/http/auth"
	router "github.com/bitcoin-sv/spv-wallet-web-backend/transports/http/endpoints/routes"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type handler struct {
	service *users.UserService
	log     *zerolog.Logger
}

// NewHandler creates new endpoint handler.
func NewHandler(s *domain.Services, log *zerolog.Logger) (router.RootEndpoints, router.APIEndpoints) {
	h := &handler{
		service: s.UsersService,
		log:     log,
	}

	prefix := "/api/v1"

	// Register root endpoints which are authorized by admin token.
	rootEndpoints := router.RootEndpointsFunc(func(router *gin.RouterGroup) {
		router.POST(prefix+"/sign-in", h.signIn)
	})

	// Register api endpoints which are authorized by session token.
	apiEndpoints := router.APIEndpointsFunc(func(router *gin.RouterGroup) {
		router.POST("/sign-out", h.signOut)
	})

	return rootEndpoints, apiEndpoints
}

// Sign in user.
//
//	@Summary Sign in user
//	@Tags user
//	@Accept json
//	@Produce json
//	@Success 200 {object} SignInResponse
//	@Router /api/v1/sign-in [post]
//	@Param data body SignInUser true "User sign in data"
func (h *handler) signIn(c *gin.Context) {
	var reqUser SignInUser
	if err := c.Bind(&reqUser); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, h.log)
		return
	}

	signInUser, err := h.service.SignInUser(reqUser.Email, reqUser.Password)
	if err != nil {
		spverrors.ErrorResponse(c, err, h.log)
		return
	}

	err = auth.UpdateSession(c, signInUser)
	if err != nil {
		h.log.Error().Msgf("Sign-in error. Session wasn't saved: %s", err)
		spverrors.ErrorResponse(c, spverrors.ErrSessionUpdate, h.log)
		return
	}

	response := SignInResponse{
		Paymail: signInUser.User.Paymail,
		Balance: signInUser.Balance,
	}
	c.JSON(http.StatusOK, response)
}

// Sign out user.
//
//	@Summary Sign out user
//	@Tags user
//	@Accept */*
//	@Produce json
//	@Success 200
//	@Router /api/v1/sign-out [post]
func (h *handler) signOut(c *gin.Context) {
	// Right now we cannot revoke access key without authentication with XPriv
	// All we can do is to terminate session

	err := auth.TerminateSession(c)
	if err != nil {
		h.log.Error().Msgf("Sign-out error. Session wasn't terminated: %s", err)
		spverrors.ErrorResponse(c, spverrors.ErrSessionTerminate, h.log)
		return
	}

	c.Status(http.StatusOK)
}
