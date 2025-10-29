package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/triliun/mcupload/backend/shared"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUser(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	resp, err := s.repo.GetOneByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	return &GetResponse{
		Username:  resp.Username,
		CreatedAt: resp.CreatedAt,
	}, nil
}

func (s *Service) GetProfile(ctx context.Context, userID uuid.UUID) (*GetProfileResponse, error) {
	resp, err := s.repo.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &GetProfileResponse{
		Username:  resp.Username,
		Email:     resp.Email,
		CreatedAt: resp.CreatedAt,
		UpdatedAt: resp.UpdatedAt,
	}, nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID uuid.UUID, req *UpdateRequest) (*UpdateResponse, error) {
	existingUser, err := s.repo.GetOneByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	existingUser.UpdatedAt = time.Now()

	if req.Username != "" {
		existingUser.Username = req.Username
	}
	if req.Email != "" {
		existingUser.Email = req.Email
	}
	if req.Password != "" {
		hashPassword, err := shared.Crypto.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}

		existingUser.Password = hashPassword
	}

	err = s.repo.Update(ctx, existingUser)
	if err != nil {
		return nil, err
	}

	return &UpdateResponse{
		Username: req.Username,
		Email:    req.Email,
	}, nil
}
