package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/triliun/mcupload/backend/middleware/auth"
	"github.com/triliun/mcupload/backend/shared"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	var hashPassword string
	hashPassword, err = shared.Crypto.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	err = s.repo.Create(ctx, &User{
		ID:        uuid,
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashPassword,
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return nil, err
	}

	return &RegisterResponse{
		Username:  req.Username,
		Email:     req.Email,
		CreatedAt: now,
	}, nil
}

func (s *Service) Login(ctx context.Context, req *LoginRequest, middleware *auth.AuthMiddleware) (*auth.TokenPair, error) {
	user, err := s.repo.FindOneByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	ErrInvalidUsernameOrPassword := errors.New("invalid username or password")

	if user == nil {
		return nil, ErrInvalidUsernameOrPassword
	}

	valid := shared.Crypto.CheckPassword(user.Password, req.Password)
	if !valid {
		return nil, ErrInvalidUsernameOrPassword
	}

	var resp *auth.TokenPair
	resp, err = middleware.Service.GenerateTokenPair(user.ID, user.Username, user.Role, middleware.JWTSecret, middleware.RefreshSecret)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Service) RefreshToken(ctx context.Context, req *RefreshTokenRequest, middleware *auth.AuthMiddleware) (*auth.TokenPair, error) {
	claims, err := middleware.Service.ValidateRefreshToken(ctx, req.RefreshToken, middleware.RefreshSecret)
	if err != nil {
		return nil, err
	}

	var user *User
	user, err = s.repo.GetOneByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	var resp *auth.TokenPair
	resp, err = middleware.Service.GenerateTokenPair(user.ID, user.Username, user.Role, middleware.JWTSecret, middleware.RefreshSecret)
	if err != nil {
		return nil, err
	}

	var id uuid.UUID
	id, err = uuid.NewV7()
	if err != nil {
		return nil, err
	}

	err = s.repo.RevokeToken(ctx, &Logout{
		ID:        id,
		TokenID:   uuid.MustParse(claims.ID),
		UserID:    claims.UserID,
		Reason:    "Refresh Token",
		ExpiresAt: claims.ExpiresAt.Time,
		RevokedAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Service) Logout(ctx context.Context, jti string, userID uuid.UUID, expires time.Time) (*LogoutResponse, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	const reason string = "Logout"

	err = s.repo.RevokeToken(ctx, &Logout{
		ID:        id,
		TokenID:   uuid.MustParse(jti),
		UserID:    userID,
		Reason:    reason,
		ExpiresAt: expires,
		RevokedAt: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return &LogoutResponse{Reason: reason}, nil
}
