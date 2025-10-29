package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken         = errors.New("invalid token")
	ErrTokenRevoked         = errors.New("invalid token revoked")
	ErrAuthHeaderIsRequired = errors.New("authorization header is required")
	ErrTokenExpired         = errors.New("token expired")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
)

const issuerClaim = "mcupload.com"

type Service struct {
	Repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) ValidateAccessToken(ctx context.Context, tokenString string, jwtSecret string) (*AccessClaims, error) {
	token, err := ParseWithClaims(tokenString, &AccessClaims{}, jwtSecret)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	var isRevoked bool
	isRevoked, err = s.Repo.IsTokenRevoked(ctx, claims.ID)
	if err != nil {
		return nil, err
	}
	if isRevoked {
		return nil, ErrTokenRevoked
	}

	return claims, nil
}

func (s *Service) ValidateRefreshToken(ctx context.Context, tokenString string, refreshSecret string) (*RefreshClaims, error) {
	token, err := ParseWithClaims(tokenString, &RefreshClaims{}, refreshSecret)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	var isRevoked bool
	isRevoked, err = s.Repo.IsTokenRevoked(ctx, claims.ID)
	if err != nil {
		return nil, err
	}
	if isRevoked {
		return nil, ErrTokenRevoked
	}

	return claims, nil
}

func (s *Service) GenerateTokenPair(userID uuid.UUID, username string, role string, jwtSecret string, refreshSecret string) (*TokenPair, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	jti := uuid.String()

	// Generate access token
	var accessToken string
	accessToken, err = s.generateAccessToken(jti, userID, username, role, jwtSecret)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	var refreshToken string
	refreshToken, err = s.generateRefreshToken(jti, userID, refreshSecret)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(24 * time.Hour).Unix(),
	}, nil
}

func (s *Service) generateAccessToken(jti string, userID uuid.UUID, username string, role string, jwtSecret string) (string, error) {
	claims := AccessClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // Access token expires in 15 minutes
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    issuerClaim,
			Subject:   userID.String(),
			ID:        jti,
			Audience:  []string{"mcupload-api"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func (s *Service) generateRefreshToken(jti string, userID uuid.UUID, refreshSecret string) (string, error) {
	claims := RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // Refresh token expires in 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    issuerClaim,
			Subject:   userID.String(),
			ID:        jti,
			Audience:  []string{"mcupload-api-refresh"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(refreshSecret))
}
