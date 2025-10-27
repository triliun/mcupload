package e2e

import (
	"net/http/httptest"
	"testing"

	"github.com/triliun/mcupload/backend/features/auth"
)

type DataConstrain interface {
	map[string]any | []map[string]any
}

type TestResponse[T DataConstrain] struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	Data       T      `json:"data"`
	Errors     any    `json:"errors"`
	NextCursor string `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
}

type TestCase struct {
	name           string
	method         string
	url            string
	requestBody    any
	setupAuth      bool
	expectedStatus int
	expectedBody   string
	checkFunc      func(*testing.T, *httptest.ResponseRecorder)
}

// Test constants
const (
	baseURL         = "/api/v1"
	contentTypeJSON = "application/json"
	authHeader      = "Authorization"
	bearerPrefix    = "Bearer "
	successContains = `"success":true`
	errorContains   = `"success":false`
)

var UserWithRoleCEO = struct{ auth.User }{
	User: auth.User{
		Username: "seth",
		Email:    "seth@example.com",
		Password: "SethPro69",
		Role:     "ceo",
	},
}
