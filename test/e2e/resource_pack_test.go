package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/triliun/mcupload/backend/features/resource_pack"
	"github.com/triliun/mcupload/backend/shared"
)

func TestResourcePackSuccess(t *testing.T) {
	const baseURL = baseURL + "/resource-pack"

	user := NewTestUserBuilder().Build()
	createTestUser(t, user.Username, user.Email, user.Password)
	tokenPair := createAuthToken(t, user.Username, user.Password)

	var packID string

	tests := []TestCase{
		{
			name:           "Create Resource Pack",
			method:         "POST",
			url:            baseURL,
			expectedStatus: http.StatusCreated,
			setupAuth:      true,
		},
		{
			name:           "Get All My Resource Pack",
			method:         "GET",
			url:            baseURL + "/my",
			expectedStatus: http.StatusOK,
			setupAuth:      true,
		},
		{
			name:           "Get My Resource Pack",
			method:         "GET",
			url:            baseURL + "/my/",
			expectedStatus: http.StatusOK,
			setupAuth:      true,
		},
		{
			name:           "Update Resource Pack",
			method:         "PUT",
			url:            baseURL + "/",
			requestBody:    resource_pack.UpdateRequest{Slug: "amphora-32", Title: "Amphora 32x32", Content: "Texture Pack PvP Amphora 32x32 by Seth", Status: "draft"},
			expectedStatus: http.StatusOK,
			setupAuth:      true,
		},
		{
			name:           "Delete Resource Pack",
			method:         "DELETE",
			url:            baseURL + "/",
			expectedStatus: http.StatusOK,
			setupAuth:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "Get My Resource Pack" || tc.name == "Update Resource Pack" || tc.name == "Delete Resource Pack" {
				tc.url += packID
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

			if tc.name == "Get All My Resource Pack" {
				var resp TestResponse[[]map[string]any]
				err := shared.Binder.BindWJSON(w.Result(), &resp)
				if err != nil {
					t.Fatal(err)
				}

				data := resp.Data[0]
				packID = data["id"].(string)
			}
		})
	}
}

func TestGetListResourcePacksSuccess(t *testing.T) {
	const baseURL = baseURL + "/resource-pack"

	user := NewTestUserBuilder().Build()
	createTestUser(t, user.Username, user.Email, user.Password)
	tokenPair := createAuthToken(t, user.Username, user.Password)

	var packID []string

	t.Run("Create 50 Resource Packs and Update Status To Published", func(t *testing.T) {
		for i := 0; i < 50; i++ {
			req, _, err := createRequest("POST", baseURL, nil)
			if err != nil {
				t.Fatal(err)
			}

			if tokenPair.AccessToken == "" {
				t.Fatalf("Access token is required")
			}

			req.Header.Set(authHeader, bearerPrefix+tokenPair.AccessToken)

			w := httptest.NewRecorder()
			GetRouter().ServeHTTP(w, req)

			assert.Equal(t, http.StatusCreated, w.Code)
			assert.Contains(t, w.Body.String(), successContains)

			var resp TestResponse[map[string]any]
			err = shared.Binder.BindWJSON(w.Result(), &resp)
			if err != nil {
				t.Fatal(err)
			}

			data := resp.Data["id"].(string)
			packID = append(packID, data)

			// Update status to "published"
			req, _, err = createRequest("PUT", baseURL+"/"+data, resource_pack.UpdateRequest{Status: "published"})
			if err != nil {
				t.Fatal(err)
			}

			if tokenPair.AccessToken == "" {
				t.Fatalf("Access token is required")
			}

			req.Header.Set(authHeader, bearerPrefix+tokenPair.AccessToken)

			w = httptest.NewRecorder()
			GetRouter().ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Contains(t, w.Body.String(), successContains)
		}
	})

	t.Run("Get List Resource Packs", func(t *testing.T) {
		const url = baseURL + "/list"
		cursor := ""
		hasMore := true

		for hasMore {
			params := "?limit=10"
			if cursor != "" {
				params += "&cursor=" + cursor
			}

			req, _, err := createRequest("GET", url+params, nil)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			GetRouter().ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Contains(t, w.Body.String(), successContains)

			var resp TestResponse[[]map[string]any]
			err = shared.Binder.BindWJSON(w.Result(), &resp)
			if err != nil {
				t.Fatal(err)
			}

			cursor = resp.NextCursor
			hasMore = resp.HasMore

		}
	})
}
