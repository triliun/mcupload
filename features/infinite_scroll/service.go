package infinite_scroll

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListResourcePacks(ctx context.Context, params ListResourcePacksParams) (*ResourcePackListResponse, error) {
	// Set default values
	if params.Limit == 0 {
		params.Limit = 20
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	resp, err := s.repo.ListResourcePacks(ctx, params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// func (s *Service) GetUsersInfiniteScroll(ctx context.Context, lastID, perPage int) (*InfiniteScrollResponse, error) {
// 	users, err := s.repo.GetUsersWithInfiniteScroll(ctx, lastID, perPage)
// 	if err != nil {
// 		return nil, err
// 	}

// 	uuid, err := uuid.NewV7()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Cek apakah masih ada data lagi
// 	hasMore := len(users) == perPage

// 	// Ambil last ID untuk cursor berikutnya
// 	var nextLastID uuid.UUID
// 	if len(users) > 0 {
// 		nextLastID = users[len(users)-1].ID
// 	}

// 	response := &InfiniteScrollResponse{
// 		Data:    users,
// 		HasMore: hasMore,
// 		LastID:  nextLastID,
// 	}

// 	return response, nil
// }
