package pack_category

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, packCategory *PackCategory) (*CreateResponse, error)
	GetAll(ctx context.Context, req *GetAllRequest) ([]GetResponse, error)
	Delete(ctx context.Context, packCategoryID uuid.UUID, authorID uuid.UUID) error

	// Helper
	GetAuthorIdFromResourcePack(ctx context.Context, packID uuid.UUID) (uuid.UUID, error)
}

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r PostgresRepository) Create(ctx context.Context, packCategory *PackCategory) (*CreateResponse, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
	INSERT INTO pack_categories (
		id, author_id, pack_id, category_id, created_at
	) VALUES (
		:id, :author_id, :pack_id, :category_id, :created_at
	)
	`

	_, err = tx.NamedExecContext(ctx, query, packCategory)
	if err != nil {
		return nil, err
	}

	var categoryName string

	query = "SELECT name FROM categories WHERE id = $1"

	err = tx.GetContext(ctx, &categoryName, query, packCategory.CategoryID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &CreateResponse{
		ID:        packCategory.ID,
		Name:      categoryName,
		CreatedAt: packCategory.CreatedAt,
	}, nil
}

func (r PostgresRepository) GetAll(ctx context.Context, req *GetAllRequest) ([]GetResponse, error) {
	var resp []GetResponse

	query := `
	SELECT
		p.id,
		c.name,
		p.created_at
		FROM pack_categories p
  	JOIN categories c ON  c.id = p.category_id
  	WHERE pack_id = $1 AND author_id = $2
	`

	err := r.db.SelectContext(ctx, &resp, query, req.PackID, req.authorID)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r PostgresRepository) Delete(ctx context.Context, packCategoryID uuid.UUID, authorID uuid.UUID) error {
	query := "DELETE FROM pack_categories WHERE id = $1 AND author_id = $2"

	result, err := r.db.ExecContext(ctx, query, packCategoryID, authorID)
	if err != nil {
		return err
	}

	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("pack category not found")
	}

	return nil
}

func (r *PostgresRepository) GetAuthorIdFromResourcePack(ctx context.Context, packID uuid.UUID) (uuid.UUID, error) {
	var authorID uuid.UUID
	query := "SELECT author_id FROM resource_packs WHERE id = $1"

	err := r.db.GetContext(ctx, &authorID, query, packID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return authorID, nil
		}
		return authorID, err
	}

	return authorID, nil
}
