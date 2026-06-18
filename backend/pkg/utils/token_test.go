package utils

import (
	"testing"
)

func TestTokenLifecycle(t *testing.T) {
	// Initialize jwtSecret for the test
	jwtSecret = []byte("my-test-secret-key")

	userID := "user-123"
	role := "admin"

	// 1. Generate Token
	token, err := GenerateToken(userID, role)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	if token == "" {
		t.Fatal("generated token is empty")
	}

	// 2. Verify Token
	claims, err := VerifyJWT(token)
	if err != nil {
		t.Fatalf("failed to verify token: %v", err)
	}

	if claims.UserId != userID {
		t.Errorf("expected userID %q, got %q", userID, claims.UserId)
	}
	if claims.Role != role {
		t.Errorf("expected role %q, got %q", role, claims.Role)
	}
}

func TestVerifyJWT_InvalidToken(t *testing.T) {
	jwtSecret = []byte("my-test-secret-key")

	_, err := VerifyJWT("invalid-token-string")
	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
}
