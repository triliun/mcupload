package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, req *User) error
	// FindOneByUsername returns id uuid, username string, email string, password string, role string
	FindOneByUsername(ctx context.Context, username string) (*User, error)
	// GetOneByID returns id uuid, username string, role string
	GetOneByID(ctx context.Context, userID uuid.UUID) (*User, error)
	RevokeToken(ctx context.Context, req *Logout) error
}

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, req *User) error {
	query := `
	INSERT INTO users (
		id, username, email, password, created_at, updated_at
	) VALUES (
		:id, :username, :email, :password, :created_at, :updated_at
	)
	`

	_, err := r.db.NamedExecContext(ctx, query, req)

	return err
}

func (r *PostgresRepository) FindOneByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	query := "SELECT id, username, email, password, role FROM users WHERE username = $1"

	err := r.db.GetContext(ctx, &user, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *PostgresRepository) GetOneByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	var user User
	query := "SELECT id, username, role FROM users WHERE id = $1"

	err := r.db.GetContext(ctx, &user, query, userID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *PostgresRepository) RevokeToken(ctx context.Context, req *Logout) error {
	query := `
	INSERT INTO revoked_tokens (
		id, token_id, user_id, reason, expires_at, revoked_at
	) VALUES (
		:id, :token_id, :user_id, :reason, :expires_at, :revoked_at
	)
	`

	_, err := r.db.NamedExecContext(ctx, query, req)

	return err
}
