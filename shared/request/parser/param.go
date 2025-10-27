package parser

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Param struct{}

// GetUUID returns error and value as uuid
func (p *Param) GetUUID(r *http.Request, key string) (uuid.UUID, error) {
	value, err := uuid.Parse(chi.URLParam(r, key))
	if err != nil {
		return value, err
	}

	return value, nil
}

// GetUint64 returns error and value as uint64
func (p *Param) GetUint64(r *http.Request, key string) (uint64, error) {
	value, err := strconv.ParseUint(chi.URLParam(r, key), 10, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}

// GetUint32 returns error and value as uint32
func (p *Param) GetUint32(r *http.Request, key string) (uint32, error) {
	value, err := strconv.ParseUint(chi.URLParam(r, key), 10, 32)
	if err != nil {
		return 0, err
	}

	return uint32(value), nil
}
