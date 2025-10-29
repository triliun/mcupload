package e2e

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/triliun/mcupload/backend/features/category"
	"github.com/triliun/mcupload/backend/shared"
)

func TestCategorySuccess(t *testing.T) {
	const baseURL = baseURL + "/category"

	var categoryID []string

	tokenPair := createAuthToken(t, UserWithRoleCEO.Username, UserWithRoleCEO.Password)

	tests := []TestCase{
		{
			name:           "Create Category: Anarcy",
			method:         "POST",
			url:            baseURL,
			requestBody:    category.CreateRequest{Name: "Anarcy"},
			expectedStatus: http.StatusCreated,
			setupAuth:      true,
		},
		{
			name:           "Create Category: Horror",
			method:         "POST",
			url:            baseURL,
			requestBody:    category.CreateRequest{Name: "Horror"},
			expectedStatus: http.StatusCreated,
			setupAuth:      true,
		},
		{
			name:           "Get All Categories",
			method:         "GET",
			url:            baseURL,
			expectedStatus: http.StatusOK,
			setupAuth:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _, err := createRequest(tc.method, tc.url, tc.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			if tc.setupAuth && tokenPair.AccessToken != "" {
				req.Header.Set(authHeader, bearerPrefix+tokenPair.AccessToken)
			}

			w := httptest.NewRecorder()
			GetRouter().ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), successContains)

			if strings.Contains(tc.name, "Get") {
				var resp TestResponse[[]map[string]any]
				err := shared.Binder.BindWJSON(w.Result(), &resp)
				if err != nil {
					t.Fatal(err)
				}

				data := resp.Data

				for _, data := range data {
					categoryID = append(categoryID, data["id"].(string))
				}
			}
		})
	}

	updateTest := TestCase{
		name:           "Update Category",
		method:         "PUT",
		url:            baseURL,
		expectedStatus: http.StatusOK,
		expectedBody:   successContains,
		setupAuth:      true,
	}

	t.Run(updateTest.name, func(t *testing.T) {
		if len(categoryID) == 0 {
			t.Fatalf("Category id is required")
		}
		for _, versionID := range categoryID {
			updateTest.requestBody = category.UpdateRequest{
				ID:   uuid.MustParse(versionID),
				Name: GenerateTestCategory(),
			}

			req, _, err := createRequest(updateTest.method, updateTest.url, updateTest.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			if !updateTest.setupAuth && tokenPair.AccessToken == "" {
				t.Fatalf("Access token is required")
			}

			req.Header.Set(authHeader, bearerPrefix+tokenPair.AccessToken)

			w := httptest.NewRecorder()
			GetRouter().ServeHTTP(w, req)

			assert.Equal(t, updateTest.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), updateTest.expectedBody)
		}
	})

	deleteTest := TestCase{
		name:           "Delete Category",
		method:         "DELETE",
		url:            baseURL + "/",
		expectedStatus: http.StatusOK,
		expectedBody:   successContains,
		setupAuth:      true,
	}

	t.Run(deleteTest.name, func(t *testing.T) {
		if len(categoryID) == 0 {
			t.Fatalf("Category id is required")
		}
		for _, versionID := range categoryID {
			url := deleteTest.url + versionID

			req, _, err := createRequest(deleteTest.method, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			if !deleteTest.setupAuth && tokenPair.AccessToken == "" {
				t.Fatalf("Access token is required")
			}
			req.Header.Set(authHeader, bearerPrefix+tokenPair.AccessToken)

			w := httptest.NewRecorder()
			GetRouter().ServeHTTP(w, req)

			assert.Equal(t, deleteTest.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), deleteTest.expectedBody)
		}
	})
}
