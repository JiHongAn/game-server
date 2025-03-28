package handler

import (
	"game-server/middleware"
	"game-server/pkg/response"
	"game-server/service"

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
	private := router.Group("/games")
	{
		private.Use(middleware.JwtAuth())
		private.GET("/", handler.GetGames)
	}
}

/**
 * @GET /games 유저 조회
 */
func (handler *GameHandler) GetGames(context *gin.Context) {
	result, err := handler.gameService.GetGames()
	if err != nil {
		response.Error(context, err)
		return
	}
	response.Success(context, result)
}
