package pack_category

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/triliun/mcupload/backend/shared"
)

// Representaions table pack_categories
type PackCategory struct {
	ID         uuid.UUID `db:"id"`
	AuthorID   uuid.UUID `db:"author_id"`
	PackID     uuid.UUID `db:"pack_id"`
	CategoryID uuid.UUID `db:"category_id"`
	CreatedAt  time.Time `db:"created_at"`
}

type CreateRequest struct {
	PackID     uuid.UUID `json:"pack_id"`
	CategoryID uuid.UUID `json:"category_id"`
	authorID   uuid.UUID
}

type CreateResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type GetAllRequest struct {
	PackID   uuid.UUID `json:"pack_id"`
	authorID uuid.UUID
}

type GetResponse struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type DeleteResponse struct {
	ID uuid.UUID `json:"id"`
}

func (r CreateRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.PackID,
			shared.Validator.Rule.UUIDRules(true)...,
		),
		validation.Field(&r.CategoryID,
			shared.Validator.Rule.UUIDRules(true)...,
		),
		validation.Field(&r.authorID,
			shared.Validator.Rule.UUIDRules(true)...,
		),
	)
}

func (r GetAllRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.PackID,
			shared.Validator.Rule.UUIDRules(true)...,
		),
		validation.Field(&r.authorID,
			shared.Validator.Rule.UUIDRules(true)...,
		),
	)
}
