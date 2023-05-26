package users

import (
	"bux-wallet/domain"
	"bux-wallet/domain/users"
	"net/http"

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
	tokens := router.Group("")
	{
		tokens.POST("/register", h.register)
	}
}

// register registers new user.
func (h *handler) register(c *gin.Context) {
	var reqUser RegisterUser
	err := c.Bind(&reqUser)

	// Check if request body is valid JSON
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// Check if sended passwords match
	if reqUser.Password != reqUser.PasswordConfirmation {
		c.JSON(http.StatusBadRequest, "passwords do not match")
		return
	}

	user, err := h.service.CreateNewUser(reqUser.Email, reqUser.Password)

	// Check if user with this email already exists or there is another error
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}
