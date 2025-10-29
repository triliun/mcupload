package auth

import (
	"errors"
	"net/http"

	"github.com/triliun/mcupload/backend/logger"
	"github.com/triliun/mcupload/backend/shared"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	Service       *Service
	JWTSecret     string
	RefreshSecret string
}

func NewAuthMiddleware(service *Service, jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		Service:       service,
		JWTSecret:     jwtSecret,
		RefreshSecret: jwtSecret + jwtSecret,
	}
}

func (m *AuthMiddleware) Validate(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := ExtractToken(r)
		if err != nil {
			logger.Warn("Token extraction failed", zap.Error(err))
			shared.HTTP.WithError(w, http.StatusUnauthorized, "Authentication failed", err.Error())
			return
		}

		var accessClaims *AccessClaims
		accessClaims, err = m.Service.ValidateAccessToken(r.Context(), tokenString, m.JWTSecret)
		if err != nil {
			logger.Warn("Token validation failed", zap.Error(err))
			shared.HTTP.WithError(w, http.StatusUnauthorized, "Authentication failed", "Invalid token")
			return
		}

		// Save accessClaims to context
		r = shared.Context.Set(r, "jwt_id", accessClaims.ID)
		r = shared.Context.Set(r, "jwt_exp", accessClaims.ExpiresAt.Time)
		r = shared.Context.Set(r, "user_id", accessClaims.UserID)
		r = shared.Context.Set(r, "user_username", accessClaims.Username)
		r = shared.Context.Set(r, "user_role", accessClaims.Role)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (m *AuthMiddleware) OnlyNonLoggedIn(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := ExtractToken(r)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		_, err = m.Service.ValidateAccessToken(r.Context(), tokenString, m.JWTSecret)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		logger.Warn("User already logged in")
		shared.HTTP.WithError(w, http.StatusForbidden, "Access denied", "You are already logged in")
	}

	return http.HandlerFunc(fn)
}

func (m *AuthMiddleware) ValidateRole(requiredRole string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			userRole, err := m.getUserRole(r)
			if err != nil {
				logger.Warn("Failed to get user role", zap.Error(err))
				shared.HTTP.WithError(w, http.StatusForbidden, "Access denied", err.Error())
				return
			}

			if userRole != requiredRole {
				logger.Warn("Insufficient permissions", zap.String("required", requiredRole), zap.String("actual", userRole))
				shared.HTTP.WithError(w, http.StatusForbidden, "Access denied", "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func (m *AuthMiddleware) ValidateRoles(requiredRoles []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			userRole, err := m.getUserRole(r)
			if err != nil {
				logger.Warn("Failed to get user role", zap.Error(err))
				shared.HTTP.WithError(w, http.StatusForbidden, "Access denied", err.Error())
				return
			}

			if !m.hasRequiredRole(userRole, requiredRoles) {
				logger.Warn("Insufficient permissions",
					zap.String("user_role", userRole),
					zap.Strings("required_roles", requiredRoles))
				shared.HTTP.WithError(w, http.StatusForbidden, "Access denied", "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

// Jika ada hierarki role (e.g., Admin > User > Guest)
var roleHierarchy = map[string]int{
	"superadmin": 100,
	"admin":      90,
	"ceo":        85,
	"manager":    70,
	"user":       50,
	"guest":      10,
}

func (m *AuthMiddleware) ValidateMinRole(minRole string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			userRole, err := m.getUserRole(r)
			if err != nil {
				shared.HTTP.WithError(w, http.StatusForbidden, "Access denied", err.Error())
				return
			}

			minLevel, minExists := roleHierarchy[minRole]
			userLevel, userExists := roleHierarchy[userRole]

			if !minExists || !userExists || userLevel < minLevel {
				logger.Warn("Insufficient permissions",
					zap.String("user_role", userRole),
					zap.Int("user_level", userLevel),
					zap.String("required_role", minRole),
					zap.Int("required_level", minLevel))
				shared.HTTP.WithError(w, http.StatusForbidden, "Access denied", "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func (m *AuthMiddleware) getUserRole(r *http.Request) (string, error) {
	userRole := shared.Context.GetUserRole(r.Context())
	if userRole == "" {
		return "", errors.New("role information missing")
	}

	return userRole, nil
}

func (m *AuthMiddleware) hasRequiredRole(userRole string, requiredRoles []string) bool {
	for _, requiredRole := range requiredRoles {
		if userRole == requiredRole {
			return true
		}
	}
	return false
}
