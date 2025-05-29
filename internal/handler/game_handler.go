package handler

import (
	"game-server/internal/middleware"
	"game-server/internal/service"

	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	gameService *service.GameService
}

func NewGameHandler(gameService *service.GameService) *GameHandler {
	return &GameHandler{
		gameService: gameService,
	}
}

func (handler *GameHandler) RegisterRoutes(router *gin.Engine) {
	// 게임 관련 API
	games := router.Group("/games")
	{
		games.Use(middleware.JwtAuth())
	}
}
