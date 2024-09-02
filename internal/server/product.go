package server

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lahnasti/go-market/internal/models"
)

func (s *Server) deleter(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if len(s.deleteChan) == 5 {
				for i := 0; i < 5; i++ {
					<-s.deleteChan
				}
				if err := s.Db.DeleteProduct(); err != nil {
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

func (s *Server) AddProductHandler(ctx *gin.Context) {
	var product models.Product
    if err := ctx.ShouldBindJSON(&product); err!= nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := s.Valid.Struct(product); err!= nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
    if err := ctx.ShouldBindJSON(&product); err!= nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := s.Valid.Struct(product); err!= nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	uid := ctx.Param("uid")
	uIdInt, err := strconv.Atoi(uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

    err = s.Db.UpdateProduct(uIdInt, product)
    if err!= nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"message": "Product updated", "product": product})
}

func (s *Server) DeleteBookHandler(ctx *gin.Context) {
	uid := ctx.Param("uid")
    uIdInt, err := strconv.Atoi(uid)
    if err!= nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    s.deleteChan <- uIdInt
    ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted", "product_uid": uIdInt})
}