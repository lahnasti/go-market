package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lahnasti/go-market/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) GetUserProfileHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil || userID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID", "error": err.Error()})
		return
	}
	user, err := s.Db.GetUserProfile(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User profile found", "user": user})
}

func (s *Server) RegisterUserHandler(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data", "error": err.Error()})
		return
	}
	if err := s.Valid.Struct(user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Not a valid user", "error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password", "error": err.Error()})
		return
	}
	user.Password = string(hash)
	id, err := s.Db.RegisterUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to register user", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "id": id})
}

func (s *Server) LoginUserHandler(ctx *gin.Context) {
	var credentials struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data", "error": err.Error()})
		return
	}
	if err := s.Valid.Struct(credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Not a valid user", "error": err.Error()})
		return
	}
	userID, err := s.Db.LoginUser(credentials.Username, credentials.Password)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		} else if err.Error() == "invalid password" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal server error", "error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User was login", "user_id": userID})
}
