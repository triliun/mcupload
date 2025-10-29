package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"user_username"`
	Role     string    `json:"user_role"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}
