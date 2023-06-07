package access

import (
	"bux-wallet/domain"
	"bux-wallet/domain/users"
	"net/http"

	"bux-wallet/transports/http/auth"
	"bux-wallet/transports/http/endpoints/api"
	router "bux-wallet/transports/http/endpoints/routes"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service users.UserService
}

// NewHandler creates new endpoint handler.
func NewHandler(s *domain.Services) (router.RootEndpoints, router.ApiEndpoints) {
	h := &handler{service: *s.UsersService}

	prefix := "/api/v1"

	// Register root endpoints which are athorized by admin token.
	rootEndpoints := router.RootEndpointsFunc(func(router *gin.RouterGroup) {
		router.POST(prefix+"/sign-in", h.signIn)
	})

	// Register api endpoints which are athorized by session token.
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
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	signInUser, err := h.service.SignInUser(reqUser.Email, reqUser.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewErrorResponseFromError(err))
		return
	}

	err = auth.UpdateSession(c, signInUser.AccessKey, signInUser.User.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewErrorResponseFromError(err))
	}

	response := SignInResponse{
		Paymail: signInUser.User.Paymail,
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

	err := h.service.SignOutUser(c.GetString(auth.SessionToken))
	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewErrorResponseFromError(err))
		return
	}

	c.Status(http.StatusOK)
}
