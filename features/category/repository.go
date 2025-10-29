package category

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, category *Category) error
	GetAll(ctx context.Context) ([]GetResponse, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Helper
	GetOneByID(ctx context.Context, id uuid.UUID) (*Category, error)
	GetAuthorIdFromResourcePack(ctx context.Context, packID uuid.UUID) (*uuid.UUID, error)
}

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r PostgresRepository) Create(ctx context.Context, category *Category) error {
	query := `
	INSERT INTO categories (
		id, name, created_at, updated_at
	) VALUES (
		:id, :name, :created_at, :updated_at
	)
	`

	_, err := r.db.NamedExecContext(ctx, query, category)

	return err
}

func (r PostgresRepository) GetAll(ctx context.Context) ([]GetResponse, error) {
	var resp []GetResponse

	query := "SELECT * FROM categories"
	err := r.db.SelectContext(ctx, &resp, query)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM categories WHERE id = $1"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("categories not found")
	}

	return nil
}

func (r *PostgresRepository) Update(ctx context.Context, category *Category) error {
	query := `
	UPDATE categories SET
		name = :name,
		updated_at = :updated_at
		WHERE id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, category)
	return err
}

// T
func (r *PostgresRepository) GetOneByID(ctx context.Context, id uuid.UUID) (*Category, error) {
	var category Category

	query := "SELECT name, updated_at FROM categories WHERE id = $1"

	err := r.db.GetContext(ctx, &category, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}

func (r *PostgresRepository) GetAuthorIdFromResourcePack(ctx context.Context, packID uuid.UUID) (*uuid.UUID, error) {
	var authorID uuid.UUID
	query := "SELECT author_id FROM resource_packs WHERE id = $1"

	err := r.db.GetContext(ctx, &authorID, query, packID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return &authorID, nil
}
