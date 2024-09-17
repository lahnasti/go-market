package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lahnasti/go-market/internal/models"
)

func (s *Server) MakePurchaseHandler(ctx *gin.Context) {
	var purchase models.Purchase
	if err := ctx.ShouldBindJSON(&purchase); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data", "error": err.Error()})
		return
	}
	if err := s.Valid.Struct(purchase); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Not a valid purchase", "error": err.Error()})
		return
	}
	if purchase.Quantity <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Quantity must be greater than 0"})
		return
	}
	purchaseID, err := s.Db.MakePurchase(purchase)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Purchase successfully completed", "id": purchaseID})
}

func (s *Server) GetUserPurchasesHandler(ctx *gin.Context) {
	userID := ctx.Param("id")
	uIdInt, err := strconv.Atoi(userID)
	if err != nil || uIdInt <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID", "error": err.Error()})
		return
	}
	purchases, err := s.Db.GetUserPurchases(uIdInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "List purchase found", "purchases": purchases})
}

func (s *Server) GetProductPurchasesHandler(ctx *gin.Context) {
	productId := ctx.Param("id")
    uIdInt, err := strconv.Atoi(productId)
    if err!= nil || uIdInt <= 0 {
        ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid product ID", "error": err.Error()})
        return
    }
    purchases, err := s.Db.GetProductPurchases(uIdInt)
    if err!= nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"message": "List purchase found", "purchases": purchases})
}