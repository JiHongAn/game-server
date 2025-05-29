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
	redisClient = redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Redis client initialized successfully")
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
	return redisClient.Set(ctx, key, value, expiration).Err()
}

func Get(key string) (string, error) {
	return redisClient.Get(ctx, key).Result()
}

func Del(key string) error {
	return redisClient.Del(ctx, key).Err()
}

func HSet(key string, field string, value string) error {
	return redisClient.HSet(ctx, key, field, value).Err()
}

func HDel(key string, field string) error {
	return redisClient.HDel(ctx, key, field).Err()
}

func HGet(key string, field string) (string, error) {
	return redisClient.HGet(ctx, key, field).Result()
}

func HGetAll(key string) (map[string]string, error) {
	return redisClient.HGetAll(ctx, key).Result()
}

func TTL(key string) (time.Duration, error) {
	return redisClient.TTL(ctx, key).Result()
}

// Exists 키 존재 여부 확인
func Exists(key string) (bool, error) {
	result, err := redisClient.Exists(ctx, key).Result()
	return result > 0, err
}

// Incr 카운터 증가
func Incr(key string) (int64, error) {
	return redisClient.Incr(ctx, key).Result()
}

// Decr 카운터 감소
func Decr(key string) (int64, error) {
	return redisClient.Decr(ctx, key).Result()
}

// SAdd Set에 멤버 추가
func SAdd(key string, member string) error {
	return redisClient.SAdd(ctx, key, member).Err()
}

// SRem Set에서 멤버 제거
func SRem(key string, member string) error {
	return redisClient.SRem(ctx, key, member).Err()
}

// SCard Set의 멤버 개수 조회
func SCard(key string) (int64, error) {
	return redisClient.SCard(ctx, key).Result()
}

// RPush List에 멤버 추가
func RPush(key string, value string) error {
	return redisClient.RPush(ctx, key, value).Err()
}

// Keys 키 패턴 조회
func Keys(pattern string) ([]string, error) {
	return redisClient.Keys(ctx, pattern).Result()
}

// SMembers Set의 모든 멤버 조회
func SMembers(key string) ([]string, error) {
	return redisClient.SMembers(ctx, key).Result()
}

// SIsMember Set에 멤버가 있는지 확인
func SIsMember(key string, member interface{}) (bool, error) {
	return redisClient.SIsMember(ctx, key, member).Result()
}
