package config

import (
	"github.com/bitcoin-sv/spv-wallet-web-backend/domain"
	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/config"
	"github.com/bitcoin-sv/spv-wallet-web-backend/spverrors"
	router "github.com/bitcoin-sv/spv-wallet-web-backend/transports/http/endpoints/routes"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type handler struct {
	services *domain.Services
	log      *zerolog.Logger
}

// PublicConfig is used for swagger generation
type PublicConfig = config.PublicConfig

// NewHandler creates new endpoint handler.
func NewHandler(s *domain.Services, log *zerolog.Logger) router.RootEndpoints {
	return &handler{
		services: s,
		log:      log,
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
		h.log.Error().Msg("Failed to get public config")
		spverrors.ErrorResponse(c, spverrors.ErrGetConfig, h.log)
		return
	}
	c.JSON(200, pubConf)
}
