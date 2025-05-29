package main

import (
	"game-server/internal/config"
	"game-server/internal/handler"
	"game-server/internal/pkg/auth"
	"game-server/internal/pkg/database"
	"game-server/internal/pkg/errors"
	"game-server/internal/pkg/response"
	"game-server/internal/service"
	"game-server/internal/socket"
	"log"
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

	// 서비스 및 핸들러 초기화
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
	// 환경별 설정 파일 로드
	envFile := config.GetEnvFile()
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Warning: %s file not found, trying .env", envFile)
	}

	// 설정 로드
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Printf("Starting server in %s environment", cfg.Env)

	// MySQL 데이터베이스 초기화
	if err := database.InitMySQL(cfg); err != nil {
		log.Fatalf("Failed to initialize MySQL: %v", err)
	}

	// Redis 초기화
	if err := database.InitRedis(cfg); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	// JWT 공개키 초기화
	if err := auth.InitJWT(); err != nil {
		log.Fatalf("Failed to initialize JWT: %v", err)
	}
	log.Println("JWT public key initialized successfully")

	// Match 서비스 및 서버 초기화
	matchService := service.NewMatchService()
	matchServer := socket.NewMatchServer(matchService)
	go func() {
		if err := matchServer.Start(cfg.Server.MatchPort); err != nil {
			log.Printf("Match server error: %v", err)
		}
	}()

	// HTTP 서버 시작
	router := setupRouter()
	log.Printf("HTTP Server starting on port %s", cfg.Server.HTTPPort)
	router.Run(":" + cfg.Server.HTTPPort)
}
