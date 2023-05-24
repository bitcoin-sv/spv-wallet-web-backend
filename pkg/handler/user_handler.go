package handler

import (
	"bux-wallet/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initUsersRoutes(router *gin.RouterGroup) {
	bins := router.Group("/users")
	{
		bins.POST("/sign-in", h.signIn)
		bins.POST("/sign-up", h.signUp)
	}
}

// DeleteBin godoc
// @Summary Sign In.
// @Tags users
// @Accept json
// @Produce json
// @Success 200
// @Router /users/sign-in [post]
// @Security Bearer
func (h *Handler) signIn(c *gin.Context) {
	user, err := h.services.Users.SignIn(service.UserSignInInput{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}

// DeleteBin godoc
// @Summary Sign Up.
// @Tags users
// @Accept json
// @Produce json
// @Success 200
// @Router /users/sign-up [post]
// @Security Bearer
func (h *Handler) signUp(c *gin.Context) {
	user, err := h.services.Users.SignUp(service.UserSignUpInput{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}
