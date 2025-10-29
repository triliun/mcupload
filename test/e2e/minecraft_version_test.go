package e2e

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/triliun/mcupload/backend/features/minecraft_version"
	"github.com/triliun/mcupload/backend/shared"
)

func TestMinecraftVersionSuccess(t *testing.T) {
	const baseURL = baseURL + "/minecraft-version"

	var minecraftVersionID []string

	tokenPair := createAuthToken(t, UserWithRoleCEO.Username, UserWithRoleCEO.Password)

	tests := []TestCase{
		{
			name:           "Create Minecraft Version - Java",
			method:         "POST",
			url:            baseURL,
			requestBody:    minecraft_version.CreateRequest{Edition: "Java", Version: GenerateTestVersion()},
			expectedStatus: http.StatusCreated,
			setupAuth:      true,
		},
		{
			name:           "Create Minecraft Version - Bedrock",
			method:         "POST",
			url:            baseURL,
			requestBody:    minecraft_version.CreateRequest{Edition: "Bedrock", Version: GenerateTestVersion()},
			expectedStatus: http.StatusCreated,
			setupAuth:      true,
		},
		{
			name:           "Get All Minecraft Version",
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

			if strings.Contains(tc.name, "Get") {
				var resp TestResponse[[]map[string]any]
				err := shared.Binder.BindWJSON(w.Result(), &resp)
				if err != nil {
					t.Fatal(err)
				}

				data := resp.Data

				for _, data := range data {
					minecraftVersionID = append(minecraftVersionID, data["id"].(string))
				}
			}
		})
	}

	updateTest := TestCase{
		name:           "Update Minecraft Version",
		method:         "PUT",
		url:            baseURL,
		expectedStatus: http.StatusOK,
		expectedBody:   successContains,
		setupAuth:      true,
	}

	t.Run(updateTest.name, func(t *testing.T) {
		if len(minecraftVersionID) == 0 {
			t.Fatalf("Minecraft version id is required")
		}
		for _, versionID := range minecraftVersionID {
			updateTest.requestBody = minecraft_version.UpdateRequest{
				ID:      uuid.MustParse(versionID),
				Edition: GenerateTestEdition(),
				Version: GenerateTestVersion(),
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
		name:           "Delete Minecraft Version",
		method:         "DELETE",
		url:            baseURL + "/",
		expectedStatus: http.StatusOK,
		expectedBody:   successContains,
		setupAuth:      true,
	}

	t.Run(deleteTest.name, func(t *testing.T) {
		if len(minecraftVersionID) == 0 {
			t.Fatalf("Minecraft version id is required")
		}
		for _, versionID := range minecraftVersionID {
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
