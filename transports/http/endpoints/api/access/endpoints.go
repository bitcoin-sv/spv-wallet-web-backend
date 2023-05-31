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
func NewHandler(s *domain.Services) router.ApiEndpoints {
	return &handler{service: *s.UsersService}
}

// RegisterApiEndpoints registers routes that are part of service API.
func (h *handler) RegisterApiEndpoints(router *gin.RouterGroup) {
	access := router.Group("")
	{
		access.POST("sign-in", h.signIn)
	}
}

// Sign in user.
//
//	@Summary Sign in user
//	@Tags user
//	@Accept */*
//	@Produce json
//	@Success 200 {object} SignInResponse
//	@Router /sign-in [post]
//	@Param data body SignInUser true "User sign in data"
func (h *handler) signIn(c *gin.Context) {
	var reqUser SignInUser
	err := c.Bind(&reqUser)

	// Check if request body is valid JSON
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	accessKeyId, err := h.service.SignInUser(reqUser.Email, reqUser.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewErrorResponseFromError(err))
		return
	}

	auth.SetAuthCookie(c, accessKeyId)
	c.Status(http.StatusOK)
}
