package token

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	domain "github.com/FrancoPesenda/eventra/internal/domain"
)

const (
	tokenExpireDuration = 24 * time.Hour
	contextKey          = "authenticated_user"
)

type Claims struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

var secretKey = os.Getenv("JWT_SECRET_KEY")

func init() {
	if secretKey == "" {
		secretKey = "eventra-secret-key-change-in-production"
	}
}

// GenerateToken crea un JWT token para un usuario
func GenerateToken(user domain.User) (string, error) {
	claims := Claims{
		UserID:   user.ID,
		UserName: user.UserName,
		Email:    user.Email,
		Role:     string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// ValidateToken valida y parsea un JWT token
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// ExtractTokenFromHeader extrae el token del header Authorization
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return parts[1], nil
}

// ClaimsToUser convierte JWT claims a domain.User
func ClaimsToUser(claims *Claims) domain.User {
	return domain.User{
		ID:       claims.UserID,
		UserName: claims.UserName,
		Email:    claims.Email,
		Role:     domain.Role(claims.Role),
	}
}
