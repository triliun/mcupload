package auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/triliun/mcupload/backend/logger"
	"github.com/triliun/mcupload/backend/middleware/auth"
	"github.com/triliun/mcupload/backend/shared"
	"go.uber.org/zap"
)

type Handler struct {
	service    *Service
	middleware *auth.AuthMiddleware
}

func NewHandler(service *Service, middleware *auth.AuthMiddleware) *Handler {
	return &Handler{
		service:    service,
		middleware: middleware,
	}
}

func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Route("/auth", func(authGroup chi.Router) {
		authGroup.Group(func(onlyNonAuth chi.Router) {
			onlyNonAuth.Use(h.middleware.OnlyNonLoggedIn)

			onlyNonAuth.Post("/register", h.Register)
			onlyNonAuth.Post("/login", h.Login)
		})

		authGroup.Post("/refresh", h.RefreshToken)

		authGroup.With(h.middleware.Validate).Post("/logout", h.Logout)
	})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FieldsWhithRequestID(ctx)
	log.Info("Creating user")

	var req RegisterRequest
	err := shared.Binder.BindJSON(r, &req)
	if err != nil {
		log.Error("Failed to bind JSON", zap.Error(err))
		shared.HTTP.WithBindError(w)
		return
	}

	err = req.Validate(r)
	if err != nil {
		log.Error("Validation failed", zap.Error(err))
		shared.HTTP.WithValidationError(w, err)
		return
	}

	errorMessage := "Failed to create user"
	successMessage := "User created successfully"

	var resp *RegisterResponse
	resp, err = h.service.Register(ctx, &req)
	if err != nil {
		log.Error(errorMessage, zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusInternalServerError, errorMessage, "Internal server error")
		return
	}

	log.Info(successMessage)
	shared.HTTP.WithSuccess(w, http.StatusCreated, successMessage, resp)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FieldsWhithRequestID(ctx)
	log.Info("Trying to Loggin")

	var req LoginRequest
	err := shared.Binder.BindJSON(r, &req)
	if err != nil {
		log.Error("Failed to bind JSON", zap.Error(err))
		shared.HTTP.WithBindError(w)
		return
	}

	err = req.Validate()
	if err != nil {
		log.Error("Validation failed", zap.Error(err))
		shared.HTTP.WithValidationError(w, err)
		return
	}

	errorMessage := "Failed to login"
	successMessage := "Login successfully"

	var resp *auth.TokenPair
	resp, err = h.service.Login(ctx, &req, h.middleware)
	if err != nil {
		log.Error(errorMessage, zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusUnauthorized, errorMessage, err.Error())
		return
	}

	log.Info(successMessage)
	shared.HTTP.WithSuccess(w, http.StatusOK, successMessage, resp)
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FieldsWhithRequestID(ctx)
	log.Info("Refreshing token")

	var req RefreshTokenRequest
	err := shared.Binder.BindJSON(r, &req)
	if err != nil {
		log.Warn("Invalid refresh token request", zap.Error(err))
		shared.HTTP.WithBindError(w)
		return
	}

	err = req.Validate()
	if err != nil {
		log.Error("Validation failed", zap.Error(err))
		shared.HTTP.WithValidationError(w, err)
		return
	}

	errorMessage := "Authentication failed"
	successMessage := "Token refreshed successfully"

	var resp *auth.TokenPair
	resp, err = h.service.RefreshToken(ctx, &req, h.middleware)
	if err != nil {
		log.Error("Failed to generate new token pair", zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusUnauthorized, errorMessage, "Invalid refresh token")
		return
	}

	log.Info(successMessage)
	shared.HTTP.WithSuccess(w, http.StatusOK, successMessage, resp)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FieldsWhithRequestID(ctx)
	log.Info("Trying to logout")

	jwtID := shared.Context.GetString(ctx, "jwt_id")
	jwtExp := shared.Context.GetTime(ctx, "jwt_exp")
	userID := shared.Context.GetUserID(ctx)

	errorMessage := "Failed to logout"
	successMessage := "Logout successfully"

	resp, err := h.service.Logout(ctx, jwtID, userID, jwtExp)
	if err != nil {
		log.Warn(errorMessage, zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusUnauthorized, errorMessage, "Invalid token")
		return
	}

	log.Info(successMessage)
	shared.HTTP.WithSuccess(w, http.StatusOK, successMessage, resp)
}
