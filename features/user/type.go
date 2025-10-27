package user

import (
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/triliun/mcupload/backend/shared"
)

// Representaions table users
type User struct {
	ID        uuid.UUID `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Role      string    `db:"role" `
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type GetProfileResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetRequest struct {
	Username string `json:"-"`
}

type GetResponse struct {
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateRequest struct {
	Username string `db:"username" json:"username,omitempty"`
	Email    string `db:"email" json:"email,omitempty"`
	Password string `db:"password" json:"password,omitempty"`
}

type UpdateResponse struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

func (r GetRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Username, shared.Validator.Rule.UsernameRules(true, false, nil)...),
	)
}

func (r UpdateRequest) Validate(request *http.Request) error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Username, shared.Validator.Rule.UsernameRules(r.Username != "", true, request)...),
		validation.Field(&r.Email, shared.Validator.Rule.EmailRules(r.Email != "", true, request)...),
		validation.Field(&r.Password, shared.Validator.Rule.PasswordRules(r.Password != "")...),
	)
}
