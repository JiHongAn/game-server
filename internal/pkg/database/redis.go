package database

import (
	"context"
	"fmt"
	"game-server/internal/config"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	ctx         = context.Background()
)

// InitRedis Redis 클라이언트 초기화
func InitRedis(cfg *config.Config) error {
	// Redis 클라이언트 옵션 설정
	options := &redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		DB:           0,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	}

	// Redis 클라이언트 생성
	client := redis.NewClient(options)

	// 연결 테스트
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	redisClient = client
	log.Printf("Redis connected successfully (host: %s:%d)", cfg.Redis.Host, cfg.Redis.Port)
	return nil
}

// GetRedisClient Redis 클라이언트 인스턴스 반환
func GetRedisClient() *redis.Client {
	return redisClient
}

// CloseRedis Redis 연결 종료
func CloseRedis() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}

func Set(key string, value string, expiration time.Duration) error {
	if redisClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}
	return redisClient.Set(ctx, key, value, expiration).Err()
}

func Get(key string) (string, error) {
	if redisClient == nil {
		return "", fmt.Errorf("Redis client not initialized")
	}
	return redisClient.Get(ctx, key).Result()
}

func Del(key string) error {
	if redisClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}
	return redisClient.Del(ctx, key).Err()
}

func HSet(key string, field string, value interface{}) error {
	if redisClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}
	return redisClient.HSet(ctx, key, field, value).Err()
}

func HGet(key string, field string) (string, error) {
	if redisClient == nil {
		return "", fmt.Errorf("Redis client not initialized")
	}
	return redisClient.HGet(ctx, key, field).Result()
}

func HGetAll(key string) (map[string]string, error) {
	if redisClient == nil {
		return nil, fmt.Errorf("Redis client not initialized")
	}
	return redisClient.HGetAll(ctx, key).Result()
}

func TTL(key string, expiration time.Duration) error {
	if redisClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}
	return redisClient.Expire(ctx, key, expiration).Err()
}

// Exists 키 존재 여부 확인
func Exists(key string) (bool, error) {
	if redisClient == nil {
		return false, fmt.Errorf("Redis client not initialized")
	}
	result, err := redisClient.Exists(ctx, key).Result()
	return result > 0, err
}

// Incr 카운터 증가
func Incr(key string) (int64, error) {
	if redisClient == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}
	return redisClient.Incr(ctx, key).Result()
}

// Decr 카운터 감소
func Decr(key string) (int64, error) {
	if redisClient == nil {
		return 0, fmt.Errorf("Redis client not initialized")
	}
	return redisClient.Decr(ctx, key).Result()
}

// SAdd Set에 멤버 추가
func SAdd(key string, members ...interface{}) error {
	if redisClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}
	return redisClient.SAdd(ctx, key, members...).Err()
}

// SRem Set에서 멤버 제거
func SRem(key string, members ...interface{}) error {
	if redisClient == nil {
		return fmt.Errorf("Redis client not initialized")
	}
	return redisClient.SRem(ctx, key, members...).Err()
}

// SMembers Set의 모든 멤버 조회
func SMembers(key string) ([]string, error) {
	if redisClient == nil {
		return nil, fmt.Errorf("Redis client not initialized")
	}
	return redisClient.SMembers(ctx, key).Result()
}

// SIsMember Set에 멤버가 있는지 확인
func SIsMember(key string, member interface{}) (bool, error) {
	if redisClient == nil {
		return false, fmt.Errorf("Redis client not initialized")
	}
	return redisClient.SIsMember(ctx, key, member).Result()
}
