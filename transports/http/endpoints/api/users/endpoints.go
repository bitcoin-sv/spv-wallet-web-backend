package users

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"bux-wallet/domain"
	"bux-wallet/domain/users"
	"bux-wallet/logging"
	"bux-wallet/transports/http/auth"
	router "bux-wallet/transports/http/endpoints/routes"
)

type handler struct {
	service users.UserService
	log     logging.Logger
}

// NewHandler creates new endpoint handler.
func NewHandler(s *domain.Services, lf logging.LoggerFactory) (router.RootEndpoints, router.ApiEndpoints) {
	h := &handler{
		service: *s.UsersService,
		log:     lf.NewLogger("users-handler"),
	}

	prefix := "/api/v1"

	// Register root endpoints.
	rootEndpoints := router.RootEndpointsFunc(func(router *gin.RouterGroup) {
		router.POST(prefix+"/user", h.register)
	})

	// Register api endpoints which are athorized by session token.
	apiEndpoints := router.ApiEndpointsFunc(func(router *gin.RouterGroup) {
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
		h.log.Errorf("Invalid payload: %s", err)
		c.JSON(http.StatusBadRequest, "Invalid request.")
		return
	}

	// Check if sended passwords match
	if reqUser.Password != reqUser.PasswordConfirmation {
		c.JSON(http.StatusBadRequest, "Passwords do not match.")
		return
	}

	newUser, err := h.service.CreateNewUser(reqUser.Email, reqUser.Password)

	// Check if user with this email already exists or there is another error
	if err != nil {
		switch err.(type) {
		case *users.UserError:
			c.JSON(http.StatusBadRequest, fmt.Sprintf("Error creating user: %v", err))
		case *users.PaymailError:
			c.JSON(http.StatusBadRequest, fmt.Sprintf("Error registering Paymail: %v", err))
		case *users.XPubError:
			c.JSON(http.StatusBadRequest, fmt.Sprintf("Error registering XPub: %v", err))
		default:
			c.JSON(http.StatusInternalServerError, "Something went wrong when creating new user. Please try again later.")
		}
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
	user, err := h.service.GetUserById(c.GetInt(auth.SessionUserId))
	if err != nil {
		h.log.Errorf("User not found: %s", err)
		c.JSON(http.StatusBadRequest, "An error occurred while getting user details")
		return
	}

	currentBalance, err := h.service.GetUserBalance(c.GetString(auth.SessionAccessKey))
	if err != nil {
		h.log.Errorf("Balance not found: %s", err)
		c.JSON(http.StatusBadRequest, "An error occurred while getting user details")
		return
	}

	response := UserResponse{
		UserId:  user.Id,
		Paymail: user.Paymail,
		Email:   user.Email,
		Balance: *currentBalance,
	}

	c.JSON(http.StatusOK, response)
}
