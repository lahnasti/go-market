package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lahnasti/go-market/internal/server"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(s *server.Server) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Настройка маршрутов аутентификации
	SetupAuthRoutes(s)

	// Настройка маршрутов продуктов
	SetupMarketRoutes(s)

	return r
}