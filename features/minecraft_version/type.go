package minecraft_version

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/triliun/mcupload/backend/shared"
)

type Version struct {
	ID        uuid.UUID `db:"id"`
	Edition   string    `db:"edition"`
	Version   string    `db:"version"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CreateRequest struct {
	Edition string `json:"edition"`
	Version string `json:"version"`
}

type CreateResponse struct {
	ID        uuid.UUID `json:"id"`
	Edition   string    `json:"edition"`
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}

type GetResponse struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Edition   string    `db:"edition" json:"edition"`
	Version   string    `db:"version" json:"version"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type UpdateRequest struct {
	ID      uuid.UUID `json:"id"`
	Edition string    `json:"edition,omitempty"`
	Version string    `json:"version,omitempty"`
}

type UpdateResponse struct {
	ID        uuid.UUID `json:"id"`
	Edition   string    `json:"edition,omitempty"`
	Version   string    `json:"version,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DeleteRequest struct {
	ID uuid.UUID `json:"id"`
}

type DeleteResponse struct {
	ID uuid.UUID `json:"id"`
}

func (r CreateRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Edition,
			shared.Validator.Rule.EditionRules(true)...,
		),
		validation.Field(&r.Version,
			shared.Validator.Rule.VersionRules(true)...,
		),
	)
}

func (r UpdateRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.ID,
			shared.Validator.Rule.UUIDRules(true)...,
		),
		validation.Field(&r.Edition,
			shared.Validator.Rule.EditionRules(r.Edition != "")...,
		),
		validation.Field(&r.Version,
			shared.Validator.Rule.VersionRules(r.Version != "")...,
		),
	)
}

func (r DeleteRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.ID,
			shared.Validator.Rule.UUIDRules(true)...,
		),
	)
}
