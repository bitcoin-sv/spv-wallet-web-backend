package users

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

	// Register root endpoints.
	rootEndpoints := router.RootEndpointsFunc(func(router *gin.RouterGroup) {
		router.POST(prefix+"/user", h.register)
	})

	// Register api endpoints which are athorized by session token.
	apiEndpoints := router.APIEndpointsFunc(func(router *gin.RouterGroup) {
		router.GET("/user", h.getUser)
	})

	return rootEndpoints, apiEndpoints
}

// register registers new user.
// @Description Register new user with given data, paymail is created based on username from sended email.
//
//	@Summary Register new user
//	@Tags user
//	@Accept json
//	@Produce json
//	@Success 200 {object} RegisterResponse
//	@Router /api/v1/user [post]
//	@Param data body RegisterUser true "User data"
func (h *handler) register(c *gin.Context) {
	var reqUser RegisterUser
	// Check if request body is valid JSON
	if err := c.Bind(&reqUser); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, h.log)
		return
	}

	// Check if sended passwords match
	if reqUser.Password != reqUser.PasswordConfirmation {
		spverrors.ErrorResponse(c, spverrors.ErrPasswordMismatch, h.log)
		return
	}

	newUser, err := h.service.CreateNewUser(reqUser.Email, reqUser.Password)
	if err != nil {
		spverrors.ErrorResponse(c, err, h.log)
		return
	}

	// Create response
	response := RegisterResponse{
		Mnemonic: newUser.Mnemonic,
		Paymail:  newUser.User.Paymail,
	}

	c.JSON(http.StatusOK, response)
}

// getUser return information about user from context.
//
//	@Summary Get user information
//	@Tags user
//	@Accept */*
//	@Produce json
//	@Success 200 {object} UserResponse
//	@Router /user [get]
func (h *handler) getUser(c *gin.Context) {
	user, err := h.service.GetUserByID(c.GetInt(auth.SessionUserID))
	if err != nil {
		h.log.Error().Msgf("User not found: %s", err)
		spverrors.ErrorResponse(c, spverrors.ErrGetUser, h.log)
		return
	}

	currentBalance, err := h.service.GetUserBalance(c.GetString(auth.SessionAccessKey))
	if err != nil {
		h.log.Error().Msgf("Balance not found: %s", err)
		spverrors.ErrorResponse(c, spverrors.ErrGetBalance, h.log)
		return
	}

	response := UserResponse{
		UserID:  user.ID,
		Paymail: user.Paymail,
		Email:   user.Email,
		Balance: *currentBalance,
	}

	c.JSON(http.StatusOK, response)
}
