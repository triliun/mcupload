package pack_category

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreatePackCategory(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	authorID, err := s.repo.GetAuthorIdFromResourcePack(ctx, req.PackID)
	if err != nil {
		return nil, err
	}

	if authorID != req.authorID {
		return nil, errors.New("not your own resource pack")
	}

	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	var resp *CreateResponse
	resp, err = s.repo.Create(ctx, &PackCategory{
		ID:         uuid,
		AuthorID:   req.authorID,
		PackID:     req.PackID,
		CategoryID: req.CategoryID,
		CreatedAt:  time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Service) GetAllPackCategories(ctx context.Context, req *GetAllRequest) ([]GetResponse, error) {
	resp, err := s.repo.GetAll(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Service) DeletePackCategory(ctx context.Context, packCategoryID uuid.UUID, authorID uuid.UUID) (*DeleteResponse, error) {
	err := s.repo.Delete(ctx, packCategoryID, authorID)
	if err != nil {
		return nil, err
	}

	return &DeleteResponse{ID: packCategoryID}, err
}
