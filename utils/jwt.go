package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"time"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func GenerateAccessToken(userID int64, email, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"type":    "access",
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func GenerateRefreshToken(userID int64, email, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"type":    "refresh",
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	log.Println("Starting token parse step 5")
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		log.Println("checking signing method...")

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("invalid signing method: %v\n", t.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		log.Println("signing method verified")
		return jwtKey, nil
	})
	if err != nil {
		log.Printf("failed to parse token: %v\n", err)
		return nil, nil, err
	}

	log.Println("token parsed successfully")
	return token, claims, nil
}

func IsTokenType(claims jwt.MapClaims, expectedType string) bool {
	log.Println("IsTokenType step 6 ")
	if typ, ok := claims["type"].(string); ok {
		return typ == expectedType
	}
	return false
}
