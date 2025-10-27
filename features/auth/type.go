package auth

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

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	// Age      int    `json:"age"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Representaions table revoked_tokens
type Logout struct {
	ID        uuid.UUID `db:"id"`
	TokenID   uuid.UUID `db:"token_id"`
	UserID    uuid.UUID `db:"user_id"`
	Reason    string    `db:"reason"`
	ExpiresAt time.Time `db:"expires_at"`
	RevokedAt time.Time `db:"revoked_at"`
}

type LogoutRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LogoutResponse struct {
	Reason string `json:"reason"`
}

func (r RegisterRequest) Validate(request *http.Request) error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Username, shared.Validator.Rule.UsernameRules(true, true, request)...),
		validation.Field(&r.Email, shared.Validator.Rule.EmailRules(true, true, request)...),
		// validation.Field(&r.Age, AgeRules()...),
		validation.Field(&r.Password, shared.Validator.Rule.PasswordRules(true)...),
	)
}

func (r LoginRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.Username, shared.Validator.Rule.UsernameRules(true, false, nil)...),
		validation.Field(&r.Password, shared.Validator.Rule.PasswordRules(true)...),
	)
}

func (r RefreshTokenRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.RefreshToken,
			validation.Required.Error("Refresh token is required"),
		),
	)
}

func (r LogoutRequest) Validate() error {
	return validation.ValidateStruct(
		&r,
		validation.Field(&r.AccessToken,
			validation.Required.Error("Access token token is required"),
		),
		validation.Field(&r.RefreshToken,
			validation.Required.Error("Refresh token is required"),
		),
	)
}
