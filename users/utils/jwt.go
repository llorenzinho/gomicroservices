package utils

import (
	"errors"
	"time"
	"users/config"
	"users/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID   uint   `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJWT(user *models.User, secret string, expirationSeconds int) (string, error) {
	claims := CustomClaims{
		UserID:   user.Id,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expirationSeconds) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "app",
			Subject:   "user-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func ValidateJwt(tokenString string, jwtCfg config.JWTConfig) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString, &CustomClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(jwtCfg.Secret), nil
		},
	)

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, jwt.ErrTokenMalformed
	}

	return claims, nil
}

func ValidateJwtHelper(c *gin.Context, jwtCfg config.JWTConfig) (*CustomClaims, int, error) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		return nil, 401, errors.New("authorization header is required")
	}

	tokenString = tokenString[len("Bearer "):]
	if tokenString == "" {
		return nil, 401, errors.New("token is required")
	}

	claims, err := ValidateJwt(tokenString, jwtCfg)
	if err != nil {
		return nil, 401, err
	}

	return claims, 200, nil
}
