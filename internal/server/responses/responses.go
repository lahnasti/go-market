package responses

import (
	"github.com/gin-gonic/gin"
)

type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

type Success struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendError(ctx *gin.Context, status int, message string, err error) {
	ctx.JSON(status, Error{
		Status:  status,
		Message: message,
		Error:   err.Error(),
	})
}

func SendSuccess(ctx *gin.Context, status int, message string, data interface{}) {
	ctx.JSON(status, Success{
		Status:  status,
		Message: message,
		Data:    data,
	})
}
