package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"strings"
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
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 days
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	log.Println("Starting token parse")
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
	if typ, ok := claims["type"].(string); ok {
		return typ == expectedType
	}
	return false
}

func ExtractUserIDAndEmailFromHeader(header string) (int64, string, error) {
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		return 0, "", fmt.Errorf("missing or malformed Authorization header")
	}

	tokenStr := strings.TrimPrefix(header, "Bearer ")
	_, claims, err := ParseToken(tokenStr)
	if err != nil {
		return 0, "", err
	}

	return ExtractUserFromClaims(claims)
}

func ExtractUserFromClaims(claims jwt.MapClaims) (int64, string, error) {
	uidFloat, ok1 := claims["user_id"].(float64)
	email, ok2 := claims["email"].(string)
	if !ok1 || !ok2 {
		return 0, "", fmt.Errorf("invalid claims")
	}
	return int64(uidFloat), email, nil
}

func ExtractUserIDEmailAndRoleFromHeader(header string) (int64, string, string, error) {
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		return 0, "", "", fmt.Errorf("missing or malformed Authorization header")
	}

	tokenStr := strings.TrimPrefix(header, "Bearer ")
	_, claims, err := ParseToken(tokenStr)
	if err != nil {
		return 0, "", "", err
	}

	userID, email, err := ExtractUserFromClaims(claims)
	if err != nil {
		return 0, "", "", err
	}

	role, ok := claims["role"].(string)
	if !ok {
		return 0, "", "", fmt.Errorf("invalid claims: role missing")
	}

	return userID, email, role, nil
}
