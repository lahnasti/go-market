package server

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lahnasti/go-market/internal/models"
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

func (s *Server) GetAllProductsHandler(ctx *gin.Context) {
	products, err := s.Db.GetAllProducts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "List of products", "products": products})
}

func (s *Server) GetProductByIDHandler(ctx *gin.Context) {
	uid := ctx.Param("id")
	uIdInt, err := strconv.Atoi(uid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid id", "error": err.Error()})
	}
	product, err := s.Db.GetProductByID(uIdInt)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Product found", "product": product})
}

func (s *Server) AddProductHandler(ctx *gin.Context) {
	var product models.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data", "error": err.Error()})
		return
	}
	if err := s.Valid.Struct(product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Not a valid product", "error": err.Error()})
		return
	}

	productUID, err := s.Db.AddProduct(product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Product added", "uid": productUID})
}

func (s *Server) UpdateProductHandler(ctx *gin.Context) {
	var product models.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := s.Valid.Struct(product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid := ctx.Param("id")
	uIdInt, err := strconv.Atoi(uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = s.Db.UpdateProduct(uIdInt, product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	product.UID = uIdInt

	ctx.JSON(http.StatusOK, gin.H{"message": "Product updated", "product": product})
}

func (s *Server) DeleteProductHandler(ctx *gin.Context) {
	uid := ctx.Param("id")
	uIdInt, err := strconv.Atoi(uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = s.Db.SetDeleteStatus(uIdInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	s.deleteChan <- uIdInt
	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully", "uid": uIdInt})
}
