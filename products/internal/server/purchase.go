package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lahnasti/go-market/products/internal/models"
	"github.com/lahnasti/go-market/products/internal/server/responses"
	"github.com/lahnasti/go-market/auth/internal/server"
)

// MakePurchaseHandler обрабатывает создание новой покупки
// @Summary Создание покупки
// @Description Создает новую покупку для указанного продукта
// @Tags Покупки
// @Accept json
// @Produce json
// @Param purchase body models.Purchase true "Данные о покупке"
// @Success 200 {object} responses.Success
// @Failure 400 {object} responses.Error
// @Failure 500 {object} responses.Error
// @Router /purchases/add [post]
func (s *Server) MakePurchaseHandler(ctx *gin.Context) {
	var purchase models.Purchase
	if err := ctx.ShouldBindJSON(&purchase); err != nil {
		responses.SendError(ctx, http.StatusBadRequest, "Invalid request data", err)
		return
	}
	if err := s.Valid.Struct(purchase); err != nil {
		responses.SendError(ctx, http.StatusBadRequest, "Not a valid purchase", err)
		return
	}
	userCheckMes, err := json.Marshal(map[string]int{
		"userID": purchase.UserID,
	})
	if err != nil {
		responses.SendError(ctx, http.StatusInternalServerError, "Failed to marshal request", err)
		return
	}
	err = s.Rabbit.PublishMessage("user_check_queue", userCheckMes)
	if err != nil {
		responses.SendError(ctx, http.StatusInternalServerError, "Failed to send user check request", err)
        return
	}
	userValid, err := s.WaitForUserCheckResponse()
	if err != nil || !userValid {
		responses.SendError(ctx, http.StatusNotFound, "User not found or invalid", err)
		return
	}
	_, err = s.Db.GetProductByID(purchase.ProductID)
	if err != nil {
		responses.SendError(ctx, http.StatusNotFound, "Product not found", err)
		return
	}

	if purchase.Quantity <= 0 {
		responses.SendError(ctx, http.StatusBadRequest, "Quantity must be greater than 0", nil)
		return
	}
	purchaseID, err := s.Db.MakePurchase(purchase)
	if err != nil {
		responses.SendError(ctx, http.StatusInternalServerError, "error", err)
		return
	}
	responses.SendSuccess(ctx, http.StatusOK, "Purchase successfully completed", purchaseID)
}

// GetUserPurchasesHandler получает покупки пользователя
// @Summary Получение списка покупок пользователя
// @Description Возвращает список покупок для указанного пользователя
// @Tags Покупки
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} responses.Success
// @Failure 400 {object} responses.Error
// @Failure 500 {object} responses.Error
// @Router /purchases/user/{id} [get]
func (s *Server) GetUserPurchasesHandler(ctx *gin.Context) {
	userID := ctx.Param("id")
	uIdInt, err := strconv.Atoi(userID)
	if err != nil || uIdInt <= 0 {
		responses.SendError(ctx, http.StatusBadRequest, "Invalid user id", err)
		return
	}
	purchases, err := s.Db.GetUserPurchases(uIdInt)
	if err != nil {
		responses.SendError(ctx, http.StatusBadRequest, "Invalid user id", err)
		return
	}
	responses.SendSuccess(ctx, http.StatusOK, "List purchase found", purchases)
}

// GetProductPurchasesHandler получает покупки по продукту
// @Summary Получение списка покупок по продукту
// @Description Возвращает список покупок для указанного продукта
// @Tags Покупки
// @Produce json
// @Param id path int true "ID продукта"
// @Success 200 {object} responses.Success
// @Failure 400 {object} responses.Error
// @Failure 500 {object} responses.Error
// @Router /purchases/product/{id} [get]
func (s *Server) GetProductPurchasesHandler(ctx *gin.Context) {
	productId := ctx.Param("id")
	uIdInt, err := strconv.Atoi(productId)
	if err != nil || uIdInt <= 0 {
		responses.SendError(ctx, http.StatusBadRequest, "Invalid product id", err)
		return
	}
	purchases, err := s.Db.GetProductPurchases(uIdInt)
	if err != nil {
		responses.SendError(ctx, http.StatusInternalServerError, "error", err)
		return
	}
	responses.SendSuccess(ctx, http.StatusOK, "List purchase found", purchases)
}
