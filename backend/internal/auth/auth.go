package auth

import (
	"golang.org/x/crypto/bcrypt"
	"time"
	"strconv"
	"fmt"
	"net/http"
	"github.com/golang-jwt/jwt/v5"
)

const key = "kp7Nbll8hKieGO2L1EQyOphkwJDVH0xD/tOp+DssAJ0Ll1t+jFnqdxK2BOmrAlMzQehX/2v4nt1xdDyXuU0CQQ=="

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID int32, username string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "takeaway-dine_in-system",
		Subject:   username,
		ID:        strconv.Itoa(int(userID)),
		ExpiresAt: &jwt.NumericDate{
			Time: time.Now().Add(expiresIn),
		},
		IssuedAt:  &jwt.NumericDate{
			Time: time.Now(),
		},
	})
	tokenString, err := token.SignedString([]byte(key))
	return tokenString, err
}

func ValidateJWT(tokenString, tokenSecret string) (string, int32, error) {
	var token *jwt.Token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return "", 0, err
	}

	id, err := strconv.ParseInt(token.Claims.(*jwt.RegisteredClaims).ID, 10, 32)
	return token.Claims.(*jwt.RegisteredClaims).Subject, int32(id), nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("No Authorization Header")
	}
	const prefix = "Bearer "
	if len(authHeader) < len(prefix) || authHeader[:len(prefix)] != prefix {
		return "", fmt.Errorf("Invalid Authorization Header")
	}
	return authHeader[len(prefix):], nil
}