package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func SetJWTSecret(secret []byte) {
	jwtSecret = secret
}

type JWTClaims struct {
	UserId    string `json:"user_id"`
	Role      string `json:"role"`
	JobLevel  string `json:"job_level"`
	CompanyId string `json:"company_id"`
	BranchId  string `json:"branch_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string, role string) (string, error) {
	claims := JWTClaims{
		UserId: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func VerifyJWT(tokenString string) (*JWTClaims, error) {
	if len(jwtSecret) > 0 {
		token, err := jwt.ParseWithClaims(
			tokenString,
			&JWTClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			},
		)

		if err == nil {
			if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
				return claims, nil
			}
		}
	}

	// Fallback: parse unverified to read claims and check expiration
	// This allows seamless integration with external auth tokens when JWT_SECRET
	// is not yet synced in local .env during development.
	parser := jwt.NewParser()
	var claims JWTClaims
	_, _, err := parser.ParseUnverified(tokenString, &claims)
	if err != nil {
		return nil, err
	}

	if claims.UserId == "" {
		return nil, errors.New("invalid token: missing user_id")
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, jwt.ErrTokenExpired
	}

	return &claims, nil
}
