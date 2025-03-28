package main

import (
	"game-server/handler"
	"game-server/pkg/errors"
	"game-server/pkg/response"
	"game-server/service"
	"game-server/socket"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func setupRouter() *gin.Engine {
	router := gin.Default()

	// 404 에러 처리
	router.NoRoute(func(context *gin.Context) {
		response.Error(context, errors.NotFound())
	})

	gameService := service.NewGameService()
	gameHandler := handler.NewGameHandler(gameService)
	gameHandler.RegisterRoutes(router)

	// 기본 엔드포인트
	router.GET("/", func(context *gin.Context) {
		context.JSON(200, gin.H{"result": time.Now().Format(time.RFC3339)})
	})
	return router
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Match 서버 시작
	matchServer := socket.NewMatchServer()
	matchPort := os.Getenv("MATCH_PORT")
	if matchPort == "" {
		matchPort = "8081" // Match 서버 기본 포트
	}
	go func() {
		if err := matchServer.Start(matchPort); err != nil {
			log.Printf("Match server error: %v", err)
		}
	}()

	// HTTP 서버 시작
	router := setupRouter()
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080" // HTTP 서버 기본 포트
	}
	router.Run(":" + httpPort)
}
