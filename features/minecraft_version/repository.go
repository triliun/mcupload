package minecraft_version

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, version *Version) error

	// GetAll retuns slice GetResponse
	GetAll(ctx context.Context) ([]GetResponse, error)
	Update(ctx context.Context, version *Version) error
	Delete(ctx context.Context, minecraftVersionID uuid.UUID) error

	// GetOneByID retuns id, edition, version
	GetOneByID(ctx context.Context, minecraftVersionID uuid.UUID) (*Version, error)
}

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, version *Version) error {
	query := `
	INSERT INTO minecraft_versions (
		id, edition, version, created_at
	) VALUES (
		:id, :edition, :version, :created_at
	)
	`

	_, err := r.db.NamedExecContext(ctx, query, version)

	return err
}

func (r PostgresRepository) GetAll(ctx context.Context) ([]GetResponse, error) {
	var resp []GetResponse

	query := "SELECT * FROM minecraft_versions"
	err := r.db.SelectContext(ctx, &resp, query)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *PostgresRepository) Update(ctx context.Context, version *Version) error {
	query := `
	UPDATE minecraft_versions SET
		edition = :edition,
		version = :version,
		updated_at = :updated_at
		WHERE id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, version)
	return err
}

func (r *PostgresRepository) Delete(ctx context.Context, minecraftVersionID uuid.UUID) error {
	query := "DELETE FROM minecraft_versions WHERE id = $1"

	result, err := r.db.ExecContext(ctx, query, minecraftVersionID)
	if err != nil {
		return err
	}

	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("id not found")
	}

	return nil
}

func (r *PostgresRepository) GetOneByID(ctx context.Context, minecraftVersionID uuid.UUID) (*Version, error) {
	var version Version
	query := "SELECT id, edition, version FROM minecraft_versions WHERE id = $1"

	err := r.db.GetContext(ctx, &version, query, minecraftVersionID)
	if err != nil {
		return nil, err
	}

	return &version, nil
}
