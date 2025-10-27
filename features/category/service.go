package category

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

func (s *Service) CreateCategory(ctx context.Context, req *CreateRequest) (*CreateResponse, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	err = s.repo.Create(ctx, &Category{
		ID:        uuid,
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return nil, err
	}

	return &CreateResponse{
		ID:        uuid,
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (s *Service) GetAllCategories(ctx context.Context) ([]GetResponse, error) {
	resp, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Service) UpdateCategory(ctx context.Context, req *UpdateRequest) (*UpdateResponse, error) {
	now := time.Now()

	err := s.repo.Update(ctx, &Category{
		ID:        req.ID,
		Name:      req.Name,
		UpdatedAt: now,
	})
	if err != nil {
		return nil, err
	}

	return &UpdateResponse{
		ID:        req.ID,
		Name:      req.Name,
		UpdatedAt: now,
	}, nil
}

func (s *Service) DeleteCategory(ctx context.Context, categoryID uuid.UUID) (*DeleteResponse, error) {
	err := s.repo.Delete(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	return &DeleteResponse{ID: categoryID}, nil
}

// // InfiniteScrollRequest untuk menerima parameter infinite scroll
// type InfiniteScrollRequest struct {
//     LastID    int `form:"last_id" json:"last_id"`
//     PerPage   int `form:"per_page" json:"per_page" binding:"min=1,max=50"`
// }

// // InfiniteScrollResponse untuk response infinite scroll
// type InfiniteScrollResponse struct {
//     Data      any `json:"data"`
//     HasMore   bool        `json:"has_more"`
//     LastID    int         `json:"last_id,omitempty"`
//     Total     int64       `json:"total,omitempty"`
// }

// func (r *PostgresRepository) GetUsersWithInfiniteScroll(ctx context.Context, lastID, perPage int) ([]PackCategory, error) {
// 	var users []PackCategory

// 	query := `
//         SELECT id, name, email, created_at
//         FROM users
//         WHERE deleted_at IS NULL
//         AND id > $1
//         ORDER BY id ASC
//         LIMIT $2
//     `

// 	err := r.db.SelectContext(ctx, &users, query, lastID, perPage)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return users, nil
// }

// func (s *Service) GetUsersInfiniteScroll(ctx context.Context, lastID, perPage int) (*InfiniteScrollResponse, error) {
//     users, err := s.repo.GetUsersWithInfiniteScroll(ctx, lastID, perPage)
//     if err != nil {
//         return nil, err
//     }

//     // Cek apakah masih ada data lagi
//     hasMore := len(users) == perPage

//     // Ambil last ID untuk cursor berikutnya
//     var nextLastID int
//     if len(users) > 0 {
//         nextLastID = users[len(users)-1].ID
//     }

//     response := &InfiniteScrollResponse{
//         Data:    users,
//         HasMore: hasMore,
//         LastID:  nextLastID,
//     }

//     return response, nil
// }

// // GetUsersInfiniteScroll - Handler untuk infinite scroll
// func (h *UserHandler) GetUsersInfiniteScroll(c *gin.Context) {
//     var request models.InfiniteScrollRequest

//     // Bind query parameters
//     if err := c.ShouldBindQuery(&request); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
//         return
//     }

//     // Set default values
//     if request.PerPage == 0 {
//         request.PerPage = 20 // Default untuk infinite scroll
//     }
//     if request.PerPage > 50 {
//         request.PerPage = 50 // Limit max per page
//     }

//     var result *models.InfiniteScrollResponse
//     var err error

//     if request.LastID == 0 {
//         // Initial load
//         result, err = h.userService.GetInitialUsers(c.Request.Context(), request.PerPage)
//     } else {
//         // Subsequent loads
//         result, err = h.userService.GetUsersInfiniteScroll(c.Request.Context(), request.LastID, request.PerPage)
//     }

//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
//         return
//     }

//     c.JSON(http.StatusOK, result)
// }
