package database

import (
	"fmt"
	"game-server/internal/config"
	"game-server/internal/domain"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	writerDB *gorm.DB
	readerDB *gorm.DB
)

// InitMySQL MySQL 데이터베이스 초기화 및 연결 (리더/라이터 분리)
func InitMySQL(cfg *config.Config) error {
	writerDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.WriterDB.Username,
		cfg.WriterDB.Password,
		cfg.WriterDB.Host,
		cfg.WriterDB.Port,
		cfg.WriterDB.Database,
	)

	readerDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.ReaderDB.Username,
		cfg.ReaderDB.Password,
		cfg.ReaderDB.Host,
		cfg.ReaderDB.Port,
		cfg.ReaderDB.Database,
	)

	var gormLogger logger.Interface
	if cfg.Env == "development" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	writer, err := gorm.Open(mysql.Open(writerDSN), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to writer MySQL database: %w", err)
	}

	reader, err := gorm.Open(mysql.Open(readerDSN), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to reader MySQL database: %w", err)
	}

	if err := autoMigrateMySQL(writer); err != nil {
		return fmt.Errorf("failed to run auto migration: %w", err)
	}

	writerDB = writer
	readerDB = reader
	log.Printf("MySQL connected successfully (env: %s, writer/reader separated)", cfg.Env)
	return nil
}

// autoMigrateMySQL 모든 모델에 대해 Auto Migration 실행
func autoMigrateMySQL(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.Game{},
		&domain.Match{},
	)
}

// GetWriterDB MySQL Writer 인스턴스 반환 (쓰기 작업용)
func GetWriterDB() *gorm.DB {
	return writerDB
}

// GetReaderDB MySQL Reader 인스턴스 반환 (읽기 작업용)
func GetReaderDB() *gorm.DB {
	return readerDB
}

// GetDB 기존 호환성을 위한 Writer DB 반환 (deprecated)
// 새로운 코드에서는 GetWriterDB() 또는 GetReaderDB()를 사용하세요
func GetDB() *gorm.DB {
	return writerDB
}

// CloseMySQL MySQL 연결 종료
func CloseMySQL() error {
	var errors []error

	if writerDB != nil {
		if sqlDB, err := writerDB.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				errors = append(errors, fmt.Errorf("failed to close writer DB: %w", err))
			}
		}
	}

	if readerDB != nil {
		if sqlDB, err := readerDB.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				errors = append(errors, fmt.Errorf("failed to close reader DB: %w", err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors closing MySQL connections: %v", errors)
	}

	log.Println("MySQL connections closed successfully")
	return nil
}
