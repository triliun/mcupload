package e2e

// go test -v ./test/e2e/main_test.go ./test/e2e/config.go ./test/e2e/auth_test.go

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/triliun/mcupload/backend/middleware/auth"
)

func TestAuthSuccess(t *testing.T) {
	var tokenPair *auth.TokenPair
	baseURL := baseURL + "/auth"

	userData := NewTestUserBuilder().Build()

	tests := []TestCase{
		{
			name:           "Register",
			method:         "POST",
			url:            baseURL + "/register",
			requestBody:    userData,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Login",
			method:         "POST",
			url:            baseURL + "/login",
			requestBody:    userData,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Refresh Access Token",
			method:         "POST",
			url:            baseURL + "/refresh",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Logout",
			method:         "POST",
			url:            baseURL + "/logout",
			setupAuth:      true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Login",
			method:         "POST",
			url:            baseURL + "/login",
			requestBody:    userData,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "Refresh Access Token" {
				tc.requestBody = tokenPair
			}

			req, _, err := createRequest(tc.method, tc.url, tc.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			if tc.setupAuth {
				if tokenPair.AccessToken == "" {
					t.Fatalf("Access token is required")
				}
				req.Header.Set(authHeader, bearerPrefix+tokenPair.AccessToken)
			}

			w := httptest.NewRecorder()
			GetRouter().ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), successContains)

			// Extract token pair from login response
			if tc.name == "Login" || tc.name == "Refresh Access Token" {
				tokenPair, err = extractTokenFromLoginResponse(w)
				if err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
