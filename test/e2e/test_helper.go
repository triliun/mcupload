package e2e

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/triliun/mcupload/backend/features/auth"
	authM "github.com/triliun/mcupload/backend/middleware/auth"
	"github.com/triliun/mcupload/backend/shared"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// TestUserBuilder implements the Builder pattern for creating test user data.
// It provides a fluent interface for constructing user registration requests
// with customizable fields while maintaining sensible defaults.
type TestUserBuilder struct {
	username string
	email    string
	password string
}

// NewTestUserBuilder creates a new TestUserBuilder with randomized default values.
// This ensures unique test data across test executions and prevents conflicts.
func NewTestUserBuilder() *TestUserBuilder {
	randomNumber := rand.Intn(10000)
	username := fmt.Sprintf("test_user_%d", randomNumber)
	email := fmt.Sprintf("test%d@example.com", randomNumber)

	return &TestUserBuilder{
		username: username,
		email:    email,
		password: "securePassword123",
	}
}

// WithUsername sets a custom username for the test user.
func (b *TestUserBuilder) WithUsername(username string) *TestUserBuilder {
	b.username = username
	return b
}

// WithEmail sets a custom email for the test user.
func (b *TestUserBuilder) WithEmail(email string) *TestUserBuilder {
	b.email = email
	return b
}

// WithPassword sets a custom password for the test user.
func (b *TestUserBuilder) WithPassword(password string) *TestUserBuilder {
	b.password = password
	return b
}

// Build constructs and returns a RegisterRequest with the configured values.
func (b *TestUserBuilder) Build() *auth.RegisterRequest {
	return &auth.RegisterRequest{
		Username: b.username,
		Email:    b.email,
		Password: b.password,
	}
}

// createTestUser helper function creates a new user via the registration endpoint.
// It validates the response to ensure successful user creation.
func createTestUser(t *testing.T, username string, email string, password string) {
	url := baseURL + "/auth/register"

	registerData := auth.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	}

	req, _, err := createRequest("POST", url, &registerData)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	GetRouter().ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), successContains)
}

// createAuthToken authenticates a user and returns the token pair from the response.
// It validates the login was successful before extracting tokens.
func createAuthToken(t *testing.T, username string, password string) *authM.TokenPair {
	url := baseURL + "/auth/login"

	loginData := auth.LoginRequest{
		Username: username,
		Password: password,
	}

	req, _, err := createRequest("POST", url, &loginData)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	GetRouter().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), successContains)

	var tokenPair *authM.TokenPair
	tokenPair, err = extractTokenFromLoginResponse(w)
	if err != nil {
		t.Fatal(err)
	}

	return tokenPair
}

// createRequest creates a new HTTP request with JSON content type.
// It marshals the request body to JSON and returns both the request and buffer.
func createRequest(method, url string, body any) (*http.Request, *bytes.Buffer, error) {
	var reqBody bytes.Buffer
	if body != nil {

		jsonData, err := shared.Binder.API.Marshal(body)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = *bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, &reqBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", contentTypeJSON)
	return req, &reqBody, nil
}

// extractTokenFromLoginResponse extracts the TokenPair from the HTTP response.
// It decodes the JSON response and maps it to the TokenPair structure.
func extractTokenFromLoginResponse(w *httptest.ResponseRecorder) (*authM.TokenPair, error) {
	var resp TestResponse[map[string]any]

	err := shared.Binder.BindWJSON(w.Result(), &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	data := resp.Data
	accessToken, ok1 := data["access_token"].(string)
	refreshToken, ok2 := data["refresh_token"].(string)
	expiresAt, ok3 := data["expires_at"].(float64)

	if !ok1 || !ok2 || !ok3 {
		return nil, fmt.Errorf("invalid token data in response")
	}

	return &authM.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    int64(expiresAt),
	}, nil
}

var (
	editions   = []string{"Java", "Bedrock"}
	categories = []string{"Bedwars", "Survival", "PvP", "Skywars"}
	letters    = []rune("abcdefghijklmnopqrstuvwxyz")
)

// GenerateTestEdition generates random edition for testing
func GenerateTestEdition() string {
	return editions[rand.Intn(len(editions))]
}

// GenerateTestCategory generates random category for testing
func GenerateTestCategory() string {
	return categories[rand.Intn(len(categories))]
}

// GenerateTestVersion generates a random version for testing
func GenerateTestVersion() string {
	versionTypes := []func() string{
		generateReleaseVersion,
		generateSnapshotVersion,
		generatePreReleaseVersion,
	}
	return versionTypes[rand.Intn(len(versionTypes))]()
}

func generateReleaseVersion() string {
	major := rand.Intn(5) + 1
	minor := rand.Intn(20) + 1

	if rand.Float32() < 0.3 {
		return strconv.Itoa(major) + "." + strconv.Itoa(minor)
	}

	patch := rand.Intn(10)
	return strconv.Itoa(major) + "." + strconv.Itoa(minor) + "." + strconv.Itoa(patch)
}

func generateSnapshotVersion() string {
	year := rand.Intn(5) + 20
	week := rand.Intn(52) + 1
	letter := string(letters[rand.Intn(len(letters))])
	return strconv.Itoa(year) + "w" + strconv.Itoa(week) + letter
}

func generatePreReleaseVersion() string {
	major := rand.Intn(5) + 1
	minor := rand.Intn(20) + 1
	preNum := rand.Intn(5) + 1
	return strconv.Itoa(major) + "." + strconv.Itoa(minor) + "-pre" + strconv.Itoa(preNum)
}
