package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	GetOneByID(ctx context.Context, userID uuid.UUID) (*User, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*User, error)
	GetOneByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, req *User) error
}

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetOneByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	var user User
	query := "SELECT * FROM users WHERE id = $1"

	err := r.db.GetContext(ctx, &user, query, userID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresRepository) GetProfile(ctx context.Context, userID uuid.UUID) (*User, error) {
	var user User
	query := "SELECT id, username, email, created_at, updated_at FROM users WHERE id = $1"

	err := r.db.GetContext(ctx, &user, query, userID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresRepository) GetOneByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	query := "SELECT username, created_at FROM users WHERE username = $1"

	err := r.db.GetContext(ctx, &user, query, username)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresRepository) Update(ctx context.Context, req *User) error {
	query := `
	UPDATE users SET
		username = :username,
		email = :email,
		password = :password,
		updated_at = :updated_at
		WHERE id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, req)
	return err
}
