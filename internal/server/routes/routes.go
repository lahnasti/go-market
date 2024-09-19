package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lahnasti/go-market/internal/server"
)

func SetupRoutes(s *server.Server) *gin.Engine {
	r := gin.Default()

	productGroup := r.Group("/products")
	{
		productGroup.GET("/", s.GetAllProductsHandler)
		productGroup.GET("/:id", s.GetProductByIDHandler)
		productGroup.POST("/add", s.AddProductHandler)
		productGroup.PUT("/:id", s.UpdateProductHandler)
		productGroup.DELETE("/:id", s.DeleteProductHandler)
	}
	purchaseGroup := r.Group("/purchases")
	{
		purchaseGroup.GET("/user/:id", s.GetUserPurchasesHandler)
		purchaseGroup.GET("/product/:id", s.GetProductPurchasesHandler)
		purchaseGroup.POST("/add", s.MakePurchaseHandler)
	}
	userGroup := r.Group("/users")
	{
		userGroup.GET(":id", s.GetUserProfileHandler)
		userGroup.POST("/register", s.RegisterUserHandler)
		userGroup.POST("/login", s.LoginUserHandler)
	}
	return r
}
