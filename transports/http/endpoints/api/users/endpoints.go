package users

import (
	"bux-wallet/domain"
	"bux-wallet/domain/users"
	"net/http"

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
	users := router.Group("/user")
	{
		users.POST("", h.register)
	}
}

// register registers new user.
//
//	@Summary Register new user
//	@Tags user
//	@Accept */*
//	@Produce json
//	@Success 200 {object} RegisterResponse
//	@Router /user [post]
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
		c.JSON(http.StatusBadRequest, api.CreateErrorResponse("passwords do not match"))
		return
	}

	mnemonic, paymail, err := h.service.CreateNewUser(reqUser.Email, reqUser.Password)

	// Check if user with this email already exists or there is another error
	if err != nil {
		c.JSON(http.StatusBadRequest, api.CreateErrorResponse(err.Error()))
		return
	}

	// Create response
	response := RegisterReposne{
		Mnemonic: mnemonic,
		Paymail:  paymail,
	}

	c.JSON(http.StatusOK, response)
}
