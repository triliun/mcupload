package resource_pack

import (
	"database/sql"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/triliun/mcupload/backend/shared"
)

type List struct {
	ID        uuid.UUID `db:"id" json:"-"`
	Title     string    `db:"title" json:"title"`
	Slug      string    `db:"slug" json:"slug"`
	Username  string    `db:"username" json:"username"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type ListRequest struct {
	shared.InfiniteScrollRequest
	Status   string `json:"status,omitempty"`
	AuthorID string `json:"author_id,omitempty"`
}

type ListResponse struct {
	List       []List
	NextCursor string
	HasMore    bool
}

// ResourcePack represents the resource_packs table
type ResourcePack struct {
	ID          uuid.UUID      `db:"id"`
	Slug        string         `db:"slug"`
	Title       string         `db:"title"`
	Content     sql.NullString `db:"content"`
	Status      string         `db:"status"`
	AuthorID    uuid.UUID      `db:"author_id"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
	PublishedAt sql.NullTime   `db:"published_at"`
}

type CreateRequest struct {
	Slug     string    `json:"-"`
	Title    string    `json:"-"`
	AuthorID uuid.UUID `json:"-"`
}

type CreateResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type GetRequest struct {
	Slug string
}

type GetResponse struct {
	Title     string    `db:"title" json:"title"`
	Username  string    `db:"username" json:"username"`
	Content   string    `db:"content" json:"content"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type GetAllMyResponse struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	Title       string     `db:"title" json:"title"`
	Status      string     `db:"status" json:"status"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	PublishedAt *time.Time `db:"published_at" json:"published_at"`
}

type GetMyResponse struct {
	Slug        string     `json:"slug"`
	Title       string     `json:"title"`
	Content     *string    `json:"content,omitempty"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}

type UpdateRequest struct {
	Slug     string    `json:"slug,omitempty"`
	Title    string    `json:"title,omitempty"`
	Content  string    `json:"content,omitempty"`
	Status   string    `json:"status,omitempty"`
	AuthorID uuid.UUID `json:"-"`
	PackID   uuid.UUID `json:"-"`
}

type UpdateResponse struct {
	Slug        string     `json:"slug,omitempty"`
	Title       string     `json:"title,omitempty"`
	Content     *string    `json:"content,omitempty"`
	Status      string     `json:"status,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}

type DeleteResponse struct {
	ID uuid.UUID `json:"id"`
}

func (r GetRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Slug,
			shared.Validator.Rule.SlugRules(true, false, nil)...,
		),
	)
}

func (r UpdateRequest) Validate(request *http.Request) error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Slug,
			shared.Validator.Rule.SlugRules(r.Slug != "", true, request)...,
		),
		validation.Field(&r.Title,
			shared.Validator.Rule.TitleRules(r.Title != "")...,
		),
		validation.Field(&r.Content,
			validation.When(r.Content != "",
				validation.Required.Error("Content is required"),
				validation.Length(10, 1500),
			),
		),
		validation.Field(&r.Status,
			validation.When(r.Status != "",
				validation.Required.Error("Status is required"),
				validation.In("draft", "published", "archived"),
			),
		),
	)
}

func (r ListRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Limit,
			validation.Required.Error("Limit is required"),
			validation.Min(5).Error("Limit must be at least 5"),
			validation.Max(20).Error("Max limit is 20"),
		),
	)
}
