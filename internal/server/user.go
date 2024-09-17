package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lahnasti/go-market/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// GetUserProfileHandler godoc
// @Summary Get user profile
// @Description Get the profile of a user by ID
// @Tags Users
// @Param id path int true "User ID"
// @Success 200 {object} gin.H{"message": "User profile found", "user": models.User}
// @Failure 400 {object} gin.H{"message": "Invalid user ID", "error": "Error details"}
// @Failure 500 {object} gin.H{"error": "Error details"}
// @Router /users/{id} [get]
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
// RegisterUserHandler godoc
// @Summary Register a new user
// @Description Create a new user in the system
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "User data"
// @Success 201 {object} gin.H{"message": "User registered successfully", "id": int}
// @Failure 400 {object} gin.H{"message": "Invalid request data", "error": "Error details"}
// @Failure 500 {object} gin.H{"message": "Failed to register user", "error": "Error details"}
// @Router /users/register [post]
func (s *Server) RegisterUserHandler(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err!= nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data", "error": err.Error()})
        return
    }
	if err := s.Valid.Struct(user); err!= nil {
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
// LoginUserHandler godoc
// @Summary Login user
// @Description Login user with username and password
// @Tags Users
// @Accept json
// @Produce json
// @Param credentials body struct{Username string; Password string} true "User credentials"
// @Success 201 {object} gin.H{"message": "User was login", "user_id": int}
// @Failure 400 {object} gin.H{"message": "Invalid request data", "error": "Error details"}
// @Failure 401 {object} gin.H{"message": "Invalid username or password"}
// @Failure 500 {object} gin.H{"message": "Internal server error", "error": "Error details"}
// @Router /users/login [post]
func (s *Server) LoginUserHandler(ctx *gin.Context) {
	var credentials struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	if err := ctx.ShouldBindJSON(&credentials); err!= nil {
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

