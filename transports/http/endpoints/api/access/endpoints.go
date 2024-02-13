package access

import (
	"errors"
	"net/http"
	"web-backend/domain"
	"web-backend/domain/users"

	"github.com/rs/zerolog"

	"web-backend/transports/http/auth"
	router "web-backend/transports/http/endpoints/routes"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service users.UserService
	log     *zerolog.Logger
}

// NewHandler creates new endpoint handler.
func NewHandler(s *domain.Services, log *zerolog.Logger) (router.RootEndpoints, router.ApiEndpoints) {
	h := &handler{
		service: *s.UsersService,
		log:     log,
	}

	prefix := "/api/v1"

	// Register root endpoints which are authorized by admin token.
	rootEndpoints := router.RootEndpointsFunc(func(router *gin.RouterGroup) {
		router.POST(prefix+"/sign-in", h.signIn)
	})

	// Register api endpoints which are authorized by session token.
	apiEndpoints := router.ApiEndpointsFunc(func(router *gin.RouterGroup) {
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
	err := c.Bind(&reqUser)

	// Check if request body is valid JSON
	if err != nil {
		h.log.Error().Msgf("Invalid payload. Error: %s", err)
		c.JSON(http.StatusBadRequest, "Invalid request.")
		return
	}

	signInUser, err := h.service.SignInUser(reqUser.Email, reqUser.Password)
	if err != nil {
		if errors.Is(err, users.ErrInvalidCredentials) {
			c.JSON(http.StatusBadRequest, "Sorry, your username or password is incorrect. Please try again.")
			return
		}

		h.log.Error().Msgf("Sign-in error: %s", err)
		c.JSON(http.StatusInternalServerError, "Something went wrong. Please try again later.")
		return
	}

	err = auth.UpdateSession(c, signInUser)
	if err != nil {
		h.log.Error().Msgf("Sign-in error. Session wasn't saved: %s", err)
		c.JSON(http.StatusBadRequest, "Something went wrong. Please try again later.")
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
	err := h.service.SignOutUser(c.GetString(auth.SessionAccessKeyId), c.GetString(auth.SessionAccessKey))
	if err != nil {
		h.log.Error().Msgf("Sign-out error: %s", err)
		c.JSON(http.StatusInternalServerError, "An error occurred during the logout process.")
		return
	}

	err = auth.TerminateSession(c)
	if err != nil {
		h.log.Error().Msgf("Sign-out error. Session wasn't terminated: %s", err)
		c.JSON(http.StatusInternalServerError, "An error occurred during the logout process.")
		return
	}

	c.Status(http.StatusOK)
}
