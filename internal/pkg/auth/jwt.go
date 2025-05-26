package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var publicKey *rsa.PublicKey

// InitJWT JWT 공개키를 초기화합니다
func InitJWT() error {
	publicKeyPEM := os.Getenv("JWT_PUBLIC_KEY")
	if publicKeyPEM == "" {
		return fmt.Errorf("JWT_PUBLIC_KEY environment variable is required")
	}

	key, err := parsePublicKey(publicKeyPEM)
	if err != nil {
		return fmt.Errorf("failed to parse JWT public key: %w", err)
	}

	publicKey = key
	return nil
}

func parsePublicKey(publicKeyPEM string) (*rsa.PublicKey, error) {
	// Base64 디코딩
	decodedKey, err := base64.StdEncoding.DecodeString(publicKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 public key: %w", err)
	}

	block, _ := pem.Decode(decodedKey)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key is not RSA public key")
	}

	return rsaKey, nil
}

// ValidateAccessToken Auth 서버에서 발급된 JWT 토큰을 공개키로 검증합니다
func ValidateAccessToken(tokenString string) (*jwt.Token, error) {
	if publicKey == nil {
		return nil, fmt.Errorf("JWT public key not initialized")
	}

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// RSA 서명 방식인지 확인
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
}

// GetUserIDFromToken 토큰에서 사용자 ID를 추출합니다
func GetUserIDFromToken(token *jwt.Token) (string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	userID, ok := claims["userId"].(string)
	if !ok {
		return "", fmt.Errorf("userId not found in token")
	}

	return userID, nil
}
