package resource_pack

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, resourcePack *ResourcePack) error
	GetOneBySlug(ctx context.Context, slug string) (*GetResponse, error)
	GetAllByAuthorID(ctx context.Context, authorID uuid.UUID) ([]GetAllMyResponse, error)
	GetOneByID(ctx context.Context, packID uuid.UUID) (*ResourcePack, error)
	List(ctx context.Context, req *ListRequest) (*ListResponse, error)
	Update(ctx context.Context, resourcePack *ResourcePack) error
	Delete(ctx context.Context, authorID uuid.UUID, packID uuid.UUID) error
}

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, resourcePack *ResourcePack) error {
	query := `
	INSERT INTO resource_packs (
		id, slug, title, author_id, created_at, updated_at
	) VALUES (
		:id, :slug, :title, :author_id, :created_at, :updated_at
	)
	`

	_, err := r.db.NamedExecContext(ctx, query, resourcePack)

	return err
}

func (r *PostgresRepository) GetOneBySlug(ctx context.Context, slug string) (*GetResponse, error) {
	var resp GetResponse
	query := `
	SELECT
		r.title,
		u.username,
		r.content,
		r.created_at,
		r.updated_at
		FROM resource_packs r
		JOIN users u ON u.id = r.author_id
		WHERE slug = $1 AND status = 'published'
	`

	err := r.db.GetContext(ctx, &resp, query, slug)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (r PostgresRepository) GetAllByAuthorID(ctx context.Context, authorID uuid.UUID) ([]GetAllMyResponse, error) {
	var resp []GetAllMyResponse

	query := `
	SELECT
	id,
	title,
	status,
	created_at,
	updated_at,
	published_at
	FROM resource_packs
	WHERE author_id = $1`

	err := r.db.SelectContext(ctx, &resp, query, authorID)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *PostgresRepository) GetOneByID(ctx context.Context, packID uuid.UUID) (*ResourcePack, error) {
	var resourcePack ResourcePack
	query := "SELECT * FROM resource_packs WHERE id = $1"

	err := r.db.GetContext(ctx, &resourcePack, query, packID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &resourcePack, nil
}

func (r *PostgresRepository) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	// Base query dengan cursor-based pagination
	query := `
		SELECT r.id, r.title, r.slug, u.username, r.created_at, r.updated_at
		FROM resource_packs r
		JOIN users u ON u.id = r.author_id
		WHERE 1=1
	`

	var args []any
	argCounter := 1

	// Cursor condition
	if req.Cursor != "" {
		query += fmt.Sprintf(" AND r.created_at < (SELECT created_at FROM resource_packs WHERE id = $%d)", argCounter)
		args = append(args, req.Cursor)
		argCounter++
	}

	// Status filter
	if req.Status != "" {
		query += fmt.Sprintf(" AND r.status = $%d", argCounter)
		args = append(args, req.Status)
		argCounter++
	} else {
		// Default hanya tampilkan yang published
		query += fmt.Sprintf(" AND r.status = $%d", argCounter)
		args = append(args, "published")
		argCounter++
	}

	// Author filter
	if req.AuthorID != "" {
		query += fmt.Sprintf(" AND r.author_id = $%d", argCounter)
		args = append(args, req.AuthorID)
		argCounter++
	}

	// Ordering dan limit
	query += fmt.Sprintf(" ORDER BY r.created_at DESC LIMIT $%d", argCounter)
	args = append(args, req.Limit+1) // Fetch one extra to check if there's more

	var list []List
	err := r.db.SelectContext(ctx, &list, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list resource packs: %w", err)
	}

	// Check if there are more items
	hasMore := len(list) > req.Limit
	if hasMore {
		// Remove the extra item
		list = list[:req.Limit]
	}

	// Calculate next cursor
	var nextCursor string
	if hasMore && len(list) > 0 {
		lastItem := list[len(list)-1]
		cursor := lastItem.ID.String()
		nextCursor = cursor
	}

	return &ListResponse{
		List:       list,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (r *PostgresRepository) Update(ctx context.Context, resourcePack *ResourcePack) error {
	query := `
	UPDATE resource_packs SET
		slug = :slug,
		title = :title,
		content = :content,
		status = :status,
		updated_at = :updated_at,
		published_at = :published_at
		WHERE id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, resourcePack)
	return err
}

func (r *PostgresRepository) Delete(ctx context.Context, authorID uuid.UUID, packID uuid.UUID) error {
	query := "DELETE FROM resource_packs WHERE author_id = $1 AND id = $2"

	result, err := r.db.ExecContext(ctx, query, authorID, packID)
	if err != nil {
		return err
	}

	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("resource pack not found")
	}

	return nil
}
