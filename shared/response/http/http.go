package http

import (
	"net/http"

	"github.com/triliun/mcupload/backend/shared/request/binder"
	types "github.com/triliun/mcupload/backend/shared/types/http"
	"github.com/triliun/mcupload/backend/shared/validator"
)

// H is an shortcut for map[string]any
type H map[string]any

type HTTP struct{}

var bind = binder.NewBinder()

// SetStatus sets the HTTP response status code.
func (h *HTTP) SetStatus(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

// SetHeader is an intelligent shortcut for set a response header,
// if value == "" this method removes the header
func (h *HTTP) SetHeader(w http.ResponseWriter, key string, value string) {
	if value == "" {
		w.Header().Del(key)
		return
	}

	w.Header().Set(key, value)
}

// WithJSON sends a JSON response using sonic with the given status code and payload
func (h *HTTP) WithJSON(w http.ResponseWriter, statusCode int, payload any) error {
	h.SetHeader(w, "Content-Type", "application/json")
	h.SetStatus(w, statusCode)

	jsonBytes, err := bind.API.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = w.Write(jsonBytes)

	return err
}

// WithSuccess sends a success JSON response with standardized structure
func (h *HTTP) WithSuccess(w http.ResponseWriter, statusCode int, message string, payload any) {
	h.WithJSON(w, statusCode, types.APIResponse{Success: true, Message: message, Data: payload})
}

// WithMetadataSuccess sends a success JSON response with metadata field
func (h *HTTP) WithMetadataSuccess(w http.ResponseWriter, statusCode int, message string, payload any, nextCursor string, hasMore bool) {
	h.WithJSON(w, statusCode, types.APIResponse{Success: true, Message: message, Data: payload, NextCursor: nextCursor, HasMore: hasMore})
}

// WithError sends an error JSON response with standardized structure
func (h *HTTP) WithError(w http.ResponseWriter, statusCode int, message string, payload any) {
	h.WithJSON(w, statusCode, types.APIResponse{Success: false, Message: message, Errors: payload})
}

// WithGeneralError sends an error response with a general error field
func (h *HTTP) WithGeneralError(w http.ResponseWriter, statusCode int, message string, error string) {
	h.WithError(w, statusCode, message, H{"general": error})
}

// WithBindError sends a bad request error for JSON or Any binding failures
func (h *HTTP) WithBindError(w http.ResponseWriter) {
	h.WithGeneralError(w, http.StatusBadRequest, "Failed to bind request", "Invalid request payload")
}

// WithValidationError sends a validation error response with field-specific errors
func (h *HTTP) WithValidationError(w http.ResponseWriter, error error) {
	validate := &validator.Validator{}
	errors := validate.ValidateResponse(error)

	h.WithError(w, http.StatusBadRequest, "Validation failed", errors)
}
