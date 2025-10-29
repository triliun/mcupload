package category

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/triliun/mcupload/backend/shared"
)

// Representaions table categories
type Category struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CreateRequest struct {
	Name string `json:"name"`
}

type CreateResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetResponse struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type UpdateRequest struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type UpdateResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DeleteResponse struct {
	ID uuid.UUID `json:"id"`
}

func (r CreateRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Name,
			validation.Required.Error("Name is required"),
			validation.Length(3, 0).Error("Name must be at least 3 characters"),
		),
	)
}

func (r UpdateRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.ID, shared.Validator.Rule.UUIDRules(true)...),
		validation.Field(&r.Name,
			validation.Required.Error("Name is required"),
			validation.Length(3, 0).Error("Name must be at least 3 characters"),
		),
	)
}
