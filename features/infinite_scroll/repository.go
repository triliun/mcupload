// repository/resource_pack_repo.go
package infinite_scroll

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	ListResourcePacks(ctx context.Context, params ListResourcePacksParams) (*ResourcePackListResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*ResourcePack, error)
}

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) ListResourcePacks(ctx context.Context, params ListResourcePacksParams) (*ResourcePackListResponse, error) {
	// Base query dengan cursor-based pagination
	query := `
		SELECT id, slug, title, content, status, author_id, created_at, updated_at
		FROM resource_packs
		WHERE 1=1
	`

	var args []any
	argCounter := 1

	// Cursor condition
	if params.Cursor != "" {
		query += fmt.Sprintf(" AND created_at < (SELECT created_at FROM resource_packs WHERE id = $%d)", argCounter)
		args = append(args, params.Cursor)
		argCounter++
	}

	// Status filter
	if params.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCounter)
		args = append(args, params.Status)
		argCounter++
	} else {
		// Default hanya tampilkan yang published
		query += fmt.Sprintf(" AND status = $%d", argCounter)
		args = append(args, "published")
		argCounter++
	}

	// Author filter
	if params.AuthorID != "" {
		query += fmt.Sprintf(" AND author_id = $%d", argCounter)
		args = append(args, params.AuthorID)
		argCounter++
	}

	// Ordering dan limit
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d", argCounter)
	args = append(args, params.Limit+1) // Fetch one extra to check if there's more

	var resourcePacks []ResourcePack
	err := r.db.SelectContext(ctx, &resourcePacks, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list resource packs: %w", err)
	}

	// Check if there are more items
	hasMore := len(resourcePacks) > params.Limit
	if hasMore {
		// Remove the extra item
		resourcePacks = resourcePacks[:params.Limit]
	}

	// Calculate next cursor
	var nextCursor *string
	if hasMore && len(resourcePacks) > 0 {
		lastItem := resourcePacks[len(resourcePacks)-1]
		cursor := lastItem.ID.String()
		nextCursor = &cursor
	}

	return &ResourcePackListResponse{
		ResourcePacks: resourcePacks,
		NextCursor:    nextCursor,
		HasMore:       hasMore,
	}, nil
}

// GetByID untuk mendapatkan single resource pack
func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*ResourcePack, error) {
	const query = `
		SELECT id, slug, title, content, status, author_id, created_at, updated_at
		FROM resource_packs
		WHERE id = $1 AND status = 'published'
	`

	var resourcePack ResourcePack
	err := r.db.GetContext(ctx, &resourcePack, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource pack: %w", err)
	}

	return &resourcePack, nil
}
