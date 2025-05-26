package middleware

import (
	"game-server/internal/pkg/auth"
	"game-server/internal/pkg/errors"
	"game-server/internal/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

func JwtAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(context, errors.Unauthorized())
			context.Abort()
			return
		}

		// Bearer 토큰 추출
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			// Bearer 접두사가 없는 경우
			response.Error(context, errors.Unauthorized())
			context.Abort()
			return
		}

		// 토큰 검증
		token, err := auth.ValidateAccessToken(tokenString)
		if err != nil {
			response.Error(context, errors.Unauthorized())
			context.Abort()
			return
		}

		// 토큰이 유효한지 확인
		if !token.Valid {
			response.Error(context, errors.Unauthorized())
			context.Abort()
			return
		}

		// 사용자 ID 추출
		userID, err := auth.GetUserIDFromToken(token)
		if err != nil {
			response.Error(context, errors.Unauthorized())
			context.Abort()
			return
		}

		// 컨텍스트에 사용자 ID 저장
		context.Set("userId", userID)
		context.Next()
	}
}
