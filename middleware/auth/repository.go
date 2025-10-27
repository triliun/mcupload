package auth

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	IsTokenRevoked(ctx context.Context, tokenID string) (bool, error)
	CleanupRevokedTokens()
}

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) IsTokenRevoked(ctx context.Context, tokenID string) (bool, error) {
	var count int

	query := "SELECT COUNT(*) FROM revoked_tokens WHERE token_id = $1 AND expires_at > NOW()"

	err := r.db.QueryRowContext(ctx, query, tokenID).Scan(&count)

	return count > 0, err
}

func (r *PostgresRepository) CleanupRevokedTokens() {
	query := "DELETE FROM revoked_tokens WHERE expires_at < NOW() - INTERVAL '7 days'"

	r.db.Exec(query)
}
