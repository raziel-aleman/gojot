package gojot

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims is the struct that will be signed and encoded in the JWT
type Claims struct {
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT token with the given user ID
func GenerateToken(userID uint64, secretKey []byte, expirationTime time.Duration) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(expirationTime)},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token from the request context
func ValidateToken(ctx context.Context, secretKey []byte) (*Claims, error) {
	tokenString, err := extractTokenFromHeader(ctx)
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// extractTokenFromHeader extracts the JWT token from the Authorization header
func extractTokenFromHeader(ctx context.Context) (string, error) {
	bearerToken := ctx.Value(ContextKeyToken).(string)

	tokenParts := strings.Split(bearerToken, "Bearer ")
	if len(tokenParts) != 2 {
		return "", errors.New("invalid authorization header format")
	}

	return tokenParts[1], nil

}
