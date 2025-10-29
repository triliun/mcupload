package resource_pack

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/triliun/mcupload/backend/shared"
)

var (
	ErrResourcePackNotFound = errors.New("resource pack not found")
	ErrNotResourcePackOwner = errors.New("you don't have permission to access this resource pack")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateResourcePack(ctx context.Context, userID uuid.UUID) (*CreateResponse, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	err = s.repo.Create(ctx, &ResourcePack{
		ID:        uuid,
		Slug:      uuid.String(),
		Title:     "New, Untitled",
		AuthorID:  userID,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return nil, err
	}

	return &CreateResponse{
		ID:        uuid,
		CreatedAt: now,
	}, nil
}

func (s *Service) GetResourcePack(ctx context.Context, req GetRequest) (*GetResponse, error) {
	resp, err := s.repo.GetOneBySlug(ctx, req.Slug)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Service) GetAllMyResourcePack(ctx context.Context, authorID uuid.UUID) ([]GetAllMyResponse, error) {
	packs, err := s.repo.GetAllByAuthorID(ctx, authorID)
	if err != nil {
		return nil, err
	}

	return packs, nil
}

func (s *Service) GetMyResourcePack(ctx context.Context, authorID uuid.UUID, packID uuid.UUID) (*GetMyResponse, error) {
	pack, err := s.repo.GetOneByID(ctx, packID)
	if err != nil {
		return nil, err
	}

	if pack.AuthorID != authorID {
		return nil, ErrNotResourcePackOwner
	}

	return &GetMyResponse{
		Slug:        pack.Slug,
		Title:       pack.Title,
		Content:     shared.SQL.NullStringToPointer(pack.Content),
		Status:      pack.Status,
		CreatedAt:   pack.CreatedAt,
		UpdatedAt:   pack.UpdatedAt,
		PublishedAt: shared.SQL.NullTimeToPointer(pack.PublishedAt),
	}, nil
}

func (s *Service) ListResourcePacks(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	// Set default values
	req.Status = "published"
	// if req.Limit == 0 {
	// 	req.Limit = 20
	// }
	// if req.Limit > 100 {
	// 	req.Limit = 100
	// }

	resp, err := s.repo.List(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Service) UpdateResourcePack(ctx context.Context, req *UpdateRequest) (*UpdateResponse, error) {
	existingPack, err := s.repo.GetOneByID(ctx, req.PackID)
	if err != nil {
		return nil, err
	}

	if existingPack.AuthorID != req.AuthorID {
		return nil, ErrNotResourcePackOwner
	}

	now := time.Now()

	if req.Slug != "" {
		existingPack.Slug = req.Slug
	}
	if req.Title != "" {
		existingPack.Title = req.Title
	}
	if req.Content != "" {
		existingPack.Content = shared.SQL.StringToNullString(&req.Content)
	}
	existingPack.UpdatedAt = now

	if req.Status != "" {
		existingPack.Status = req.Status
		if req.Status == "published" {
			existingPack.PublishedAt = shared.SQL.TimeToNullTime(&now)
		}
	}

	err = s.repo.Update(ctx, existingPack)
	if err != nil {
		return nil, err
	}

	var resp UpdateResponse

	if req.Slug != "" {
		resp.Slug = req.Slug
	}
	if req.Title != "" {
		resp.Title = req.Title
	}
	if req.Content != "" {
		resp.Content = &req.Content
	}
	resp.UpdatedAt = now

	if req.Status != "" {
		resp.Status = req.Status
		if req.Status == "published" {
			resp.PublishedAt = &now
		}
	}

	return &resp, nil
}

func (s *Service) Delete(ctx context.Context, authorID uuid.UUID, packID uuid.UUID) (*DeleteResponse, error) {
	err := s.repo.Delete(ctx, authorID, packID)
	if err != nil {
		return nil, err
	}

	return &DeleteResponse{ID: packID}, nil
}
