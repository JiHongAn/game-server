package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// Config 애플리케이션 설정
type Config struct {
	Server   ServerConfig
	WriterDB WriterDBConfig
	ReaderDB ReaderDBConfig
	Redis    RedisConfig
	JWT      JWTConfig
	AWS      AWSConfig
	Env      string
}

// ServerConfig 서버 설정
type ServerConfig struct {
	HTTPPort  string
	MatchPort string
}

// WriterDBConfig 데이터베이스 설정 (Writer)
type WriterDBConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

// ReaderDBConfig 데이터베이스 설정 (Reader)
type ReaderDBConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

// RedisConfig Redis 설정
type RedisConfig struct {
	Host string
	Port int
}

// JWTConfig JWT 설정
type JWTConfig struct {
	PublicKey string
}

// AWSConfig AWS 설정
type AWSConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	S3Bucket        string
}

// Load 환경에 따라 설정을 로드
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			HTTPPort:  getEnv("PORT"),
			MatchPort: getEnv("MATCH_PORT"),
		},
		WriterDB: WriterDBConfig{
			Host:     getEnv("WRITER_DB_HOST"),
			Port:     getEnvAsInt("WRITER_DB_PORT"),
			Username: getEnv("WRITER_DB_USERNAME"),
			Password: getEnv("WRITER_DB_PASSWORD"),
			Database: getEnv("WRITER_DB_DATABASE"),
		},
		ReaderDB: ReaderDBConfig{
			Host:     getEnv("READER_DB_HOST"),
			Port:     getEnvAsInt("READER_DB_PORT"),
			Username: getEnv("READER_DB_USERNAME"),
			Password: getEnv("READER_DB_PASSWORD"),
			Database: getEnv("READER_DB_DATABASE"),
		},
		Redis: RedisConfig{
			Host: getEnv("REDIS_HOST"),
			Port: getEnvAsInt("REDIS_PORT"),
		},
		JWT: JWTConfig{
			PublicKey: getEnv("JWT_PUBLIC_KEY"),
		},
		AWS: AWSConfig{
			Region:          getEnv("AWS_REGION"),
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID"),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY"),
			S3Bucket:        getEnv("AWS_S3_BUCKET"),
		},
		Env: getEnv("ENV"),
	}

	return cfg, nil
}

// GetEnvFile 환경에 따른 설정 파일 경로 반환
func GetEnvFile() string {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development" // ENV만 기본값 허용
	}
	return fmt.Sprintf("configs/.env.%s", env)
}

// getEnv 필수 환경변수 (없으면 서버 종료)
func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("❌ Required environment variable %s is not set. Server cannot start.", key)
	}
	return value
}

// getEnvAsInt 필수 정수형 환경변수 (없거나 잘못된 값이면 서버 종료)
func getEnvAsInt(key string) int {
	value := getEnv(key)
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("❌ Environment variable %s must be a valid integer, got: %s. Server cannot start.", key, value)
	}
	return intValue
}
