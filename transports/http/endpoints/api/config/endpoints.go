package config

import (
	"github.com/bitcoin-sv/spv-wallet-web-backend/domain"
	router "github.com/bitcoin-sv/spv-wallet-web-backend/transports/http/endpoints/routes"
	"github.com/gin-gonic/gin"
)

type handler struct {
	services *domain.Services
}

func NewHandler(s *domain.Services) router.RootEndpoints {
	handler := &handler{
		services: s,
	}
	prefix := "/api/v1"

	rootEndpoints := router.RootEndpointsFunc(func(router *gin.RouterGroup) {
		router.GET(prefix+"/config", handler.getPublicConfig)
	})

	return rootEndpoints
}

// getConfig returns config fields exposed to clients.
//
//	@Summary Get config returns config fields exposed to clients
//	@Tags sharedconfig
//	@Produce json
//	@Success 200 {object} PublicConfig
//	@Router /api/v1/config [get]
func (h *handler) getPublicConfig(c *gin.Context) {
	pubConf := h.services.ConfigService.GetPublicConfig()
	if pubConf == nil {
		c.JSON(500, "Failed to get public config")
		return
	}
	c.JSON(200, pubConf)
}
