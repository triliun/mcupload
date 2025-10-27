package parser

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type Query struct{}

// Get returns value as string
// if error returns ""
func (q *Query) Get(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// GetUUID returns error and value as uuid
func (q *Query) GetUUID(r *http.Request, key string) (uuid.UUID, error) {
	value, err := uuid.Parse(q.Get(r, key))
	if err != nil {
		return value, err
	}

	return value, nil
}

// GetInt returns value as int
// if error returns 0
func (q *Query) GetInt(r *http.Request, key string) int {
	value, err := strconv.Atoi(q.Get(r, key))
	if err != nil {
		return 0
	}

	return value
}

// GetUint64 returns error and value as uint64
func (q *Query) GetUint64(r *http.Request, key string) (uint64, error) {
	value, err := strconv.ParseUint(q.Get(r, key), 10, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}

// GetUint32 returns error and value as uint32
func (q *Query) GetUint32(r *http.Request, key string) (uint32, error) {
	value, err := strconv.ParseUint(q.Get(r, key), 10, 32)
	if err != nil {
		return 0, err
	}

	return uint32(value), nil
}
