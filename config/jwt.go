package config

import (
	"loan/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	JWTSecret = []byte(getEnvOrDefault("JWT_SECRET", "your-256-bit-secret"))
)

type JWTClaim struct {
	UserID string        `json:"user_id"`
	Roles  []models.Role `json:"roles"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string, roles []models.Role) (string, error) {
	claims := JWTClaim{
		UserID: userID,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

func ValidateToken(tokenString string) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(t *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaim); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
