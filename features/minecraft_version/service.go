package minecraft_version

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateVersion(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	err = s.repo.Create(ctx, &Version{
		ID:        uuid,
		Edition:   req.Edition,
		Version:   req.Version,
		CreatedAt: now,
	})
	if err != nil {
		return nil, err
	}

	return &CreateResponse{
		ID:        uuid,
		Edition:   req.Edition,
		Version:   req.Version,
		CreatedAt: now,
	}, nil
}

func (s *Service) GetAllVersions(ctx context.Context) ([]GetResponse, error) {
	resp, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Service) UpdateVersion(ctx context.Context, req *UpdateRequest) (*UpdateResponse, error) {
	existingVersion, err := s.repo.GetOneByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	existingVersion.UpdatedAt = now

	if req.Edition != "" {
		existingVersion.Edition = req.Edition
	}
	if req.Version != "" {
		existingVersion.Version = req.Version
	}

	err = s.repo.Update(ctx, existingVersion)
	if err != nil {
		return nil, err
	}

	return &UpdateResponse{
		ID:        req.ID,
		Edition:   req.Edition,
		Version:   req.Version,
		UpdatedAt: now,
	}, nil
}

func (s *Service) DeleteVersion(ctx context.Context, minecraftVersionID uuid.UUID) (*DeleteResponse, error) {
	err := s.repo.Delete(ctx, minecraftVersionID)
	if err != nil {
		return nil, err
	}

	return &DeleteResponse{ID: minecraftVersionID}, nil
}
