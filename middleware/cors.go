package middleware

import (
	"net/http"

	"github.com/triliun/mcupload/backend/shared"
)

func CORS(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		shared.HTTP.SetHeader(w, "Access-Control-Allow-Origin", "https://mcupload.com")
		shared.HTTP.SetHeader(w, "Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		shared.HTTP.SetHeader(w, "Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			shared.HTTP.SetStatus(w, http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
