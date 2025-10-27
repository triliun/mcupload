package e2e

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserSuccess(t *testing.T) {
	baseURL := baseURL + "/user"

	user := NewTestUserBuilder().Build()
	createTestUser(t, user.Username, user.Email, user.Password)
	tokenPair := createAuthToken(t, user.Username, user.Password)

	// userUpdateData := NewTestUserBuilder()

	tests := []TestCase{
		{
			name:           "Get User",
			method:         "GET",
			url:            baseURL + "/" + user.Username,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Profile",
			method:         "GET",
			url:            baseURL + "/profile",
			setupAuth:      true,
			expectedStatus: http.StatusOK,
		},
		// {
		// 	name:           "Update Profile",
		// 	method:         "PUT",
		// 	url:            baseURL + "/profile",
		// 	requestBody:    userUpdateData,
		// 	setupAuth:      true,
		// 	expectedStatus: http.StatusOK,
		// },
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
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
		})
	}

	// Concurrent test
	t.Run("Concurrent get user", func(t *testing.T) {
		var wg sync.WaitGroup
		requests := 100

		for i := 0; i < requests; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				req, _, err := createRequest("GET", baseURL+"/seth", nil)
				if err != nil {
					t.Fatal(err)
				}

				w := httptest.NewRecorder()
				GetRouter().ServeHTTP(w, req)

				assert.Equal(t, http.StatusOK, w.Code)
				assert.Contains(t, w.Body.String(), successContains)
			}()
		}
		wg.Wait()
	})
}
