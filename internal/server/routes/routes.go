package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lahnasti/go-market/internal/server"
)

func ProductRoutes (r *gin.Engine, server *server.Server) {
	productGroup := r.Group("/products") 
	{
		productGroup.GET("/", server.GetAllProductsHandler)
		productGroup.POST("/add", server.AddProductHandler)
		productGroup.PUT("/:uid", server.UpdateProductHandler)
		productGroup.DELETE("/:uid", server.DeleteBookHandler)
	}
}