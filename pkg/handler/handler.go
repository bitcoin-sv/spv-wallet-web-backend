package handler

import (
	"bux-wallet/config"
	"bux-wallet/docs"
	"bux-wallet/pkg/service"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Services
	cfg      *config.HTTPConfig
}

func NewHandler(services *service.Services, c *config.HTTPConfig) *Handler {
	return &Handler{
		services: services,
		cfg:      c,
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()
	router.Use(h.CorsMiddleware())

	docs.SwaggerInfo.BasePath = "/api/v1"
	swagger := router.Group("/swagger", gin.BasicAuth(gin.Accounts{
		h.cfg.Swagger.Login: h.cfg.Swagger.Password,
	}))
	swagger.GET("/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	v1 := router.Group("/api/v1", h.tokenIdentity)

	h.initUsersRoutes(v1)

	return router
}
