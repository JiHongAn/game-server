package database

import (
	"fmt"
	"game-server/internal/config"
	"game-server/internal/domain"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 데이터베이스 초기화 및 연결
func InitDB(cfg *config.Config) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Seoul",
		cfg.Database.Host,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Database,
		cfg.Database.Port,
	)

	// GORM 로거 설정 (개발 환경에서만 상세 로그)
	var gormLogger logger.Interface
	if cfg.Env == "development" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto Migration 실행
	if err := autoMigrate(db); err != nil {
		return fmt.Errorf("failed to run auto migration: %w", err)
	}

	// 개발 환경에서 샘플 데이터 생성
	if cfg.Env == "development" {
		if err := seedData(db); err != nil {
			log.Printf("Warning: failed to seed data: %v", err)
		}
	}

	DB = db
	log.Printf("Database connected successfully (env: %s)", cfg.Env)
	return nil
}

// autoMigrate 모든 모델에 대해 Auto Migration 실행
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.Game{},
		&domain.Match{},
	)
}

// seedData 개발용 샘플 데이터 생성
func seedData(db *gorm.DB) error {
	// 게임 데이터 생성
	games := []domain.Game{
		{Name: "Rock Paper Scissors", Description: "가위바위보 게임", MaxPlayers: 2},
		{Name: "Tic Tac Toe", Description: "틱택토 게임", MaxPlayers: 2},
		{Name: "Chess", Description: "체스 게임", MaxPlayers: 2},
	}

	for _, game := range games {
		var existingGame domain.Game
		if err := db.Where("name = ?", game.Name).First(&existingGame).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&game).Error; err != nil {
					return fmt.Errorf("failed to create game %s: %w", game.Name, err)
				}
				log.Printf("Created sample game: %s", game.Name)
			}
		}
	}

	return nil
}

// GetDB 데이터베이스 인스턴스 반환
func GetDB() *gorm.DB {
	return DB
}
