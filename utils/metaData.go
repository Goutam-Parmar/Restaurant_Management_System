package utils

import (
	"fmt"
	"strings"
)

type AuthClaims struct {
	UserID int64
	Email  string
	Role   string
}

func ExtractAuthClaims(header string) (*AuthClaims, error) {
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		return nil, fmt.Errorf("missing or malformed Authorization header")
	}

	tokenStr := strings.TrimPrefix(header, "Bearer ")
	_, claims, err := ParseToken(tokenStr)
	if err != nil {
		return nil, err
	}

	uidFloat, ok1 := claims["user_id"].(float64)
	email, ok2 := claims["email"].(string)
	role, ok3 := claims["role"].(string)

	if !ok1 || !ok2 || !ok3 {
		return nil, fmt.Errorf("invalid claims in token")
	}

	return &AuthClaims{
		UserID: int64(uidFloat),
		Email:  email,
		Role:   role,
	}, nil
}
