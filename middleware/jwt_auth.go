package middleware

import (
	"game-server/pkg/auth"
	"game-server/pkg/errors"
	"game-server/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JwtAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(context, errors.Unauthorized())
			context.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := auth.ValidateAccessToken(tokenString)

		if err != nil {
			response.Error(context, errors.Unauthorized())
			context.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			context.Set("userId", claims["userId"])
			context.Next()
		} else {
			response.Error(context, errors.Unauthorized())
			context.Abort()
		}
	}
}
