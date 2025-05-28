package auth

import (
	"testing"
	"time"
)

func TestValidateJWT_ValidToken(t *testing.T) {
	userID := int32(123)
	username := "testuser"
	expiresIn := time.Minute * 5

	// Generate a valid token using the same key
	tokenString, err := MakeJWT(userID, username, expiresIn)
	if err != nil {
		t.Fatalf("failed to create JWT: %v", err)
	}

	gotUsername, gotUserID, err := ValidateJWT(tokenString, key)
	if err != nil {
		t.Fatalf("ValidateJWT returned error: %v", err)
	}
	if gotUsername != username {
		t.Errorf("expected username %q, got %q", username, gotUsername)
	}
	if gotUserID != userID {
		t.Errorf("expected userID %d, got %d", userID, gotUserID)
	}
}

func TestValidateJWT_InvalidSignature(t *testing.T) {
	userID := int32(456)
	username := "anotheruser"
	expiresIn := time.Minute * 5

	tokenString, err := MakeJWT(userID, username, expiresIn)
	if err != nil {
		t.Fatalf("failed to create JWT: %v", err)
	}

	// Use a wrong secret
	_, _, err = ValidateJWT(tokenString, "wrong_secret")
	if err == nil {
		t.Error("expected error for invalid signature, got nil")
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := int32(789)
	username := "expireduser"
	expiresIn := -time.Minute // already expired

	tokenString, err := MakeJWT(userID, username, expiresIn)
	if err != nil {
		t.Fatalf("failed to create JWT: %v", err)
	}

	_, _, err = ValidateJWT(tokenString, key)
	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
}

func TestValidateJWT_MalformedToken(t *testing.T) {
	_, _, err := ValidateJWT("not.a.jwt.token", key)
	if err == nil {
		t.Error("expected error for malformed token, got nil")
	}
}
