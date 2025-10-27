package infinite_scroll

import (
	"time"

	"github.com/google/uuid"
)

// InfiniteScrollRequest untuk menerima parameter infinite scroll
type InfiniteScrollRequest struct {
	LastID  int `form:"last_id" json:"last_id"`
	PerPage int `form:"per_page" json:"per_page" binding:"min=1,max=50"`
}

// InfiniteScrollResponse untuk response infinite scroll
type InfiniteScrollResponse struct {
	Data    any       `json:"data"`
	HasMore bool      `json:"has_more"`
	LastID  uuid.UUID `json:"last_id,omitempty"`
	Total   int64     `json:"total,omitempty"`
}

// Representaions table resource_packs
type ResourcePack struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	Slug        string     `db:"slug" json:"slug"`
	Title       string     `db:"title" json:"title"`
	Content     *string    `db:"content" json:"content,omitempty"`
	Status      string     `db:"status" json:"status"`
	AuthorID    uuid.UUID  `db:"author_id" json:"author_id"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	PublishedAt *time.Time `db:"published_at" json:"published_at,omitempty"`
}

type ResourcePackListResponse struct {
	ResourcePacks []ResourcePack `json:"resource_packs"`
	NextCursor    *string        `json:"next_cursor,omitempty"`
	HasMore       bool           `json:"has_more"`
}

type ListResourcePacksParams struct {
	Cursor   string `json:"cursor" query:"cursor"`
	Limit    int    `json:"limit" query:"limit" validate:"max=100"`
	Status   string `json:"status" query:"status"`
	AuthorID string `json:"author_id" query:"author_id"`
}
