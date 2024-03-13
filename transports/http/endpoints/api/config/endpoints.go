package config

import (
	"github.com/bitcoin-sv/spv-wallet-web-backend/domain"
	router "github.com/bitcoin-sv/spv-wallet-web-backend/transports/http/endpoints/routes"
	"github.com/gin-gonic/gin"
)

type handler struct {
	services *domain.Services
}

// NewHandler creates new endpoint handler.
func NewHandler(s *domain.Services) router.RootEndpoints {
	return &handler{
		services: s,
	}
}

// RegisterEndpoints registers endpoints in root context of application.
func (h *handler) RegisterEndpoints(router *gin.RouterGroup) {
	prefix := "/api/v1"
	router.GET(prefix+"/config", h.getPublicConfig)
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
