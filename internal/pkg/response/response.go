package response

import (
	"game-server/internal/pkg/errors"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Result  interface{} `json:"result,omitempty"`
	Message string      `json:"message,omitempty"`
	Code    string      `json:"code,omitempty"`
}

func Success(context *gin.Context, data interface{}) {
	code := 200
	if context.Request.Method == "POST" {
		code = 201
	}

	context.JSON(code, Response{
		Result: data,
	})
}

func Error(context *gin.Context, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		context.JSON(appErr.StatusCode, gin.H{
			"code":    appErr.Code,
			"message": appErr.Message,
		})
		return
	}

	context.JSON(500, gin.H{
		"code":    "INTERNAL_SERVER_ERROR",
		"message": err.Error(),
	})
}
