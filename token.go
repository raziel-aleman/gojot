package gojot

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	ID    uint64 `json:"user_id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Claims is the struct that will be signed and encoded in the JWT.
type Claims struct {
	User User `json:"user"`
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT token with the given user ID.
func GenerateToken(user User, secretKey []byte, expirationTime time.Duration) (string, error) {
	claims := Claims{
		User: user,
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

// ValidateToken validates a JWT token from the request.
func ValidateToken(tokenString string, secretKey []byte) (*Claims, error) {

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
