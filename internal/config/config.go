package config

import (
	"os"
	"strconv"
)

// Config 애플리케이션 설정
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	AWS      AWSConfig
}

// ServerConfig 서버 설정
type ServerConfig struct {
	HTTPPort  string
	MatchPort string
}

// DatabaseConfig 데이터베이스 설정
type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
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

// Load 환경변수에서 설정을 로드
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			HTTPPort:  getEnv("PORT", "8080"),
			MatchPort: getEnv("MATCH_PORT", "8081"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			Username: getEnv("DB_USERNAME", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Database: getEnv("DB_DATABASE", "gameserver"),
		},
		JWT: JWTConfig{
			PublicKey: getEnv("JWT_PUBLIC_KEY", ""),
		},
		AWS: AWSConfig{
			Region:          getEnv("AWS_REGION", "us-east-1"),
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
			S3Bucket:        getEnv("AWS_S3_BUCKET", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
