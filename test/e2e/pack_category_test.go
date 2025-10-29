package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/triliun/mcupload/backend/features/category"
	"github.com/triliun/mcupload/backend/features/pack_category"
	"github.com/triliun/mcupload/backend/shared"
)

func TestPackCategorySuccess(t *testing.T) {
	const baseURLCategory = baseURL + "/category"
	const baseURLResourcePack = baseURL + "/resource-pack"
	const baseURL = baseURL + "/pack-category"

	user := NewTestUserBuilder().Build()
	createTestUser(t, user.Username, user.Email, user.Password)
	tokenPair := createAuthToken(t, user.Username, user.Password)

	// ceoTokenPair for create category
	ceoTokenPair := createAuthToken(t, UserWithRoleCEO.Username, UserWithRoleCEO.Password)

	var categoryID uuid.UUID
	var packID uuid.UUID
	var packCategoryID uuid.UUID

	tests := []TestCase{
		{
			// Create category
			name:           "Create Category",
			method:         "POST",
			url:            baseURLCategory,
			requestBody:    category.CreateRequest{Name: "Horror"},
			expectedStatus: http.StatusCreated,
			setupAuth:      true,
		},
		{
			name:           "Get All Categories",
			method:         "GET",
			url:            baseURLCategory,
			expectedStatus: http.StatusOK,
			setupAuth:      true,
		},
		{
			// Create resource pack
			name:           "Create Resource Pack",
			method:         "POST",
			url:            baseURLResourcePack,
			expectedStatus: http.StatusCreated,
			setupAuth:      true,
		},
		{
			// TestPackCategory
			name:           "Create Pack Category",
			method:         "POST",
			url:            baseURL,
			expectedStatus: http.StatusCreated,
			setupAuth:      true,
		},
		{
			name:           "Get All Pack Category",
			method:         "GET",
			url:            baseURL,
			expectedStatus: http.StatusOK,
			setupAuth:      true,
		},
		{
			name:           "Delete Pack Category",
			method:         "DELETE",
			url:            baseURL + "/",
			expectedStatus: http.StatusOK,
			setupAuth:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "Create Pack Category" {
				tc.requestBody = pack_category.CreateRequest{PackID: packID, CategoryID: categoryID}
			}

			if tc.name == "Get All Pack Category" {
				tc.requestBody = pack_category.GetAllRequest{PackID: packID}
			}

			if tc.name == "Delete Pack Category" {
				tc.url += packCategoryID.String()
			}

			req, _, err := createRequest(tc.method, tc.url, tc.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			if tc.setupAuth {
				if tokenPair.AccessToken == "" {
					t.Fatalf("Access token is required")
				}

				if tc.name == "Create Category" {
					req.Header.Set(authHeader, bearerPrefix+ceoTokenPair.AccessToken)
				} else {
					req.Header.Set(authHeader, bearerPrefix+tokenPair.AccessToken)
				}
			}

			w := httptest.NewRecorder()
			GetRouter().ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), successContains)

			if tc.name == "Get All Categories" {
				var resp TestResponse[[]map[string]any]
				err := shared.Binder.BindWJSON(w.Result(), &resp)
				if err != nil {
					t.Fatal(err)
				}

				data := resp.Data[0]
				categoryID = uuid.MustParse(data["id"].(string))
			}

			if tc.name == "Create Resource Pack" {
				var resp TestResponse[map[string]any]
				err := shared.Binder.BindWJSON(w.Result(), &resp)
				if err != nil {
					t.Fatal(err)
				}

				data := resp.Data
				packID = uuid.MustParse(data["id"].(string))
			}

			if tc.name == "Get All Pack Category" {
				var resp TestResponse[[]map[string]any]
				err := shared.Binder.BindWJSON(w.Result(), &resp)
				if err != nil {
					t.Fatal(err)
				}

				data := resp.Data[0]
				packCategoryID = uuid.MustParse(data["id"].(string))
			}
		})
	}
}
