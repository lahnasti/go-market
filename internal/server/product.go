package server

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lahnasti/go-market/internal/models"
	"github.com/lahnasti/go-market/internal/server/responses"
)

func (s *Server) deleter(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
			if len(s.deleteChan) == 5 {
				for i := 0; i < 5; i++ {
					<-s.deleteChan
				}
				if err := s.Db.DeleteProducts(); err != nil {
					s.ErrorChan <- err
					return
				}
			}
		}
	}
}

// GetAllProductsHandler godoc
// @Summary Get all products
// @Description Получить список всех продуктов
// @Tags products
// @Produce json
// @Success 200 {object} responses.Success
// @Failure 500 {object} responses.Error
// @Router /products [get]
func (s *Server) GetAllProductsHandler(ctx *gin.Context) {
	products, err := s.Db.GetAllProducts()
	if err != nil {
		//responses.SendError(ctx, http.StatusInternalServerError, "message", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//responses.SendSuccess(ctx, http.StatusOK, "List of products", products)
	ctx.JSON(http.StatusOK, gin.H{"products": products})
}

// GetProductByIDHandler godoc
// @Summary Get a product by ID
// @Description Получить продукт по ID
// @Tags products
// @Param id path int true "Product ID"
// @Produce json
// @Success 200 {object} responses.Success
// @Failure 400 {object} responses.Error
// @Failure 404 {object} responses.Error
// @Router /products/{id} [get]
func (s *Server) GetProductByIDHandler(ctx *gin.Context) {
	uid := ctx.Param("id")
	uIdInt, err := strconv.Atoi(uid)
	if err != nil {
		responses.SendError(ctx, http.StatusBadRequest, "Invalid id", err)
		return
	}
	product, err := s.Db.GetProductByID(uIdInt)
	if err != nil {
		responses.SendError(ctx, http.StatusNotFound, "Product not found", err)
		return
	}
	responses.SendSuccess(ctx, http.StatusOK, "Product found", product)
}

// AddProductHandler godoc
// @Summary Add a new product
// @Description Добавить новый продукт
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product data"
// @Success 201 {object} responses.Success
// @Failure 400 {object} responses.Error
// @Failure 500 {object} responses.Error
// @Router /products [post]
func (s *Server) AddProductHandler(ctx *gin.Context) {
	var product models.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		responses.SendError(ctx, http.StatusBadRequest, "Invalid request data", err)
		return
	}
	if err := s.Valid.Struct(product); err != nil {
		responses.SendError(ctx, http.StatusBadRequest, "Not a valid product", err)
		return
	}

	productUID, err := s.Db.AddProduct(product)
	if err != nil {
		responses.SendError(ctx, http.StatusInternalServerError, "error", err)
		return
	}
	responses.SendSuccess(ctx, http.StatusCreated, "Product added", productUID)
}

// UpdateProductHandler godoc
// @Summary Update a product
// @Description Обновить данные продукта
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body models.Product true "Product data"
// @Success 200 {object} responses.Success
// @Failure 400 {object} responses.Error
// @Failure 500 {object} responses.Error
// @Router /products/{id} [put]
func (s *Server) UpdateProductHandler(ctx *gin.Context) {
	var product models.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		responses.SendError(ctx, http.StatusBadRequest, "Invalid request data", err)
		return
	}
	if err := s.Valid.Struct(product); err != nil {
		responses.SendError(ctx, http.StatusBadRequest, "Not a valid product", err)
		return
	}

	uid := ctx.Param("id")
	uIdInt, err := strconv.Atoi(uid)
	if err != nil {
		responses.SendError(ctx, http.StatusInternalServerError, "error", err)
		return
	}

	err = s.Db.UpdateProduct(uIdInt, product)
	if err != nil {
		responses.SendError(ctx, http.StatusInternalServerError, "error", err)
		return
	}

	product.UID = uIdInt
	responses.SendSuccess(ctx, http.StatusOK, "Product updated", product)
}

// / DeleteProductHandler godoc
// @Summary Delete a product
// @Description Удалить продукт по ID
// @Tags products
// @Param id path int true "Product ID"
// @Produce json
// @Success 200 {object} responses.Success
// @Failure 400 {object} responses.Error
// @Failure 500 {object} responses.Error
// @Router /products/{id} [delete]
func (s *Server) DeleteProductHandler(ctx *gin.Context) {
	uid := ctx.Param("id")
	uIdInt, err := strconv.Atoi(uid)
	if err != nil {
		responses.SendError(ctx, http.StatusBadRequest, "Invalid id", err)
		return
	}
	err = s.Db.SetDeleteStatus(uIdInt)
	if err != nil {
		responses.SendError(ctx, http.StatusInternalServerError, "error", err)
		return
	}
	s.deleteChan <- uIdInt
	responses.SendSuccess(ctx, http.StatusOK, "Product deleted", uIdInt)
}
