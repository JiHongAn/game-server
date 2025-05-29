package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config 애플리케이션 설정
type Config struct {
	Env    string
	Server struct {
		HTTPPort  string
		MatchPort string
	}
	WriterDB MySQLConfig
	ReaderDB MySQLConfig
	Redis    RedisConfig
}

// MySQLConfig MySQL 설정
type MySQLConfig struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

// DSN 데이터베이스 연결 문자열 생성
func (c *MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)
}

// RedisConfig Redis 설정
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// Load 환경에 따라 설정을 로드
func Load() (*Config, error) {
	cfg := &Config{
		Env: getEnv("APP_ENV"),
	}

	// 서버 설정
	cfg.Server.HTTPPort = getEnv("HTTP_PORT")
	cfg.Server.MatchPort = getEnv("MATCH_PORT")

	// Writer DB 설정 (필수)
	cfg.WriterDB.Host = getEnv("WRITER_DB_HOST")
	cfg.WriterDB.Port = getEnvAsInt("WRITER_DB_PORT")
	cfg.WriterDB.Database = getEnv("WRITER_DB_NAME")
	cfg.WriterDB.Username = getEnv("WRITER_DB_USER")
	cfg.WriterDB.Password = getEnv("WRITER_DB_PASSWORD")

	// Reader DB 설정 (필수)
	cfg.ReaderDB.Host = getEnv("READER_DB_HOST")
	cfg.ReaderDB.Port = getEnvAsInt("READER_DB_PORT")
	cfg.ReaderDB.Database = getEnv("READER_DB_NAME")
	cfg.ReaderDB.Username = getEnv("READER_DB_USER")
	cfg.ReaderDB.Password = getEnv("READER_DB_PASSWORD")

	// Redis 설정
	cfg.Redis.Host = getEnv("REDIS_HOST")
	cfg.Redis.Port = getEnv("REDIS_PORT")
	cfg.Redis.Password = getEnv("REDIS_PASSWORD")
	cfg.Redis.DB = getEnvAsInt("REDIS_DB")

	return cfg, nil
}

// getEnv 필수 환경변수 (없으면 에러 반환)
func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("%s is required", key))
	}
	return value
}

// getEnvAsInt 필수 정수형 환경변수 (없거나 잘못된 값이면 에러 반환)
func getEnvAsInt(key string) int {
	value := getEnv(key)
	intValue, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("invalid %s value: %v", key, err))
	}
	return intValue
}

// GetEnvFile 환경에 따른 설정 파일 경로 반환
func GetEnvFile() string {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development" // ENV만 기본값 허용
	}
	return fmt.Sprintf("configs/.env.%s", env)
}
