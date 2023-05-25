package registration

import (
	"bux-wallet/domain"

	router "bux-wallet/transports/http/endpoints/routes"

	"github.com/gin-gonic/gin"
)

// NewHandler creates new endpoint handler.
func NewHandler(s *domain.Services) router.RootEndpoints {
	return router.RootEndpointsFunc(func(router *gin.RouterGroup) {
		router.POST("register", registerUser)
	})
}

func registerUser(c *gin.Context) {
	c.Status(200)
}
