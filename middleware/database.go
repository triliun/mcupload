package middleware

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/triliun/mcupload/backend/shared"
)

// Database returns a middleware for database instace
func Database(db *sqlx.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = shared.Context.Set(r, "db", db)

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
