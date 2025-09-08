package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

const (
	ScopeAuth = "authentication"
	ScopeRefresh = "refresh"
)

type Token struct {
	PlainText string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    string    `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

func GenerateToken(userID string, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	emptyBytes := make([]byte, 32)
	_, err := rand.Read(emptyBytes)
	if err != nil {
		return nil, err
	}

	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(emptyBytes)
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]
	return token, nil
}

// HashToken returns a sha256 hash of the provided token string
func HashToken(t string) []byte {
	hash := sha256.Sum256([]byte(t))
	return hash[:]
}

type JWTPair struct {
	AccessToken   string    `json:"access_token"`
	RefreshToken  string    `json:"refresh_token"`
	AccessExpiry  time.Time `json:"access_expires_at"`
	RefreshExpiry time.Time `json:"refresh_expires_at"`
}

// GenerateJWT creates a signed JWT with the provided subject, ttl and scope using HMAC SHA256
func GenerateJWT(subject string, ttl time.Duration, scope string) (string, time.Time, error) {
	err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
	
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", time.Time{}, errors.New("JWT_SECRET not set")
	}

	now := time.Now()
	expiresAt := now.Add(ttl)
	claims := jwt.MapClaims{
		"sub":   subject,
		"scope": scope,
		"exp":   expiresAt.Unix(),
		"iat":   now.Unix(),
		"iss":   "bookstore",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expiresAt, nil
}

// GenerateJWTPair creates both access and refresh JWTs using the same secret.
// accessTTL is typically short (e.g., 15m), refreshTTL longer (e.g., 7d).
func GenerateJWTPair(subject string, accessTTL, refreshTTL time.Duration) (*JWTPair, error) {
	accessToken, accessExp, err := GenerateJWT(subject, accessTTL, ScopeAuth)
	if err != nil {
		return nil, err
	}
	refreshToken, refreshExp, err := GenerateJWT(subject, refreshTTL, ScopeRefresh)
	if err != nil {
		return nil, err
	}
	return &JWTPair{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		AccessExpiry:  accessExp,
		RefreshExpiry: refreshExp,
	}, nil
}
