package users

import (
	"net/http"

	"bux-wallet/domain"
	"bux-wallet/domain/users"

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
	user := router.Group("/user")
	{
		user.POST("", h.register)
	}
}

// register registers new user.
// @Description Register new user with given data, paymail is created based on username from sended email.
//
//	@Summary Register new user
//	@Tags user
//	@Accept */*
//	@Produce json
//	@Success 200 {object} RegisterResponse
//	@Router /api/v1/user [post]
//	@Param data body RegisterUser true "User data"
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
		c.JSON(http.StatusBadRequest, api.NewErrorResponseFromString("passwords do not match"))
		return
	}

	newUser, err := h.service.CreateNewUser(reqUser.Email, reqUser.Password)

	// Check if user with this email already exists or there is another error
	if err != nil {
		c.JSON(http.StatusBadRequest, api.NewErrorResponseFromError(err))
		return
	}

	// Create response
	response := RegisterResponse{
		Mnemonic: newUser.Mnemonic,
		Paymail:  newUser.User.Paymail,
	}

	c.JSON(http.StatusOK, response)
}
