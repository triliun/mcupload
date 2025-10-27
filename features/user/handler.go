package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/triliun/mcupload/backend/logger"
	authM "github.com/triliun/mcupload/backend/middleware/auth"
	"github.com/triliun/mcupload/backend/shared"
	"go.uber.org/zap"
)

type Handler struct {
	service    *Service
	middleware *authM.AuthMiddleware
}

func NewHandler(service *Service, middleware *authM.AuthMiddleware) *Handler {
	return &Handler{
		service:    service,
		middleware: middleware,
	}
}

func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Route("/user", func(userGroup chi.Router) {
		userGroup.Group(func(public chi.Router) {
			public.Get("/{username}", h.GetUser)
		})

		userGroup.Route("/profile", func(profileGroup chi.Router) {
			profileGroup.Use(h.middleware.Validate)

			profileGroup.Get("/", h.GetProfile)
			profileGroup.Put("/", h.UpdateProfile)
			// profileGroup.Get("/settings", h.GetSettings)
			// profileGroup.Put("/settings", h.UpdateSettings)
		})

		// userGroup.Route("/staff", func(staff chi.Router) {
		// 	staff.Use(h.middleware.ValidateRoles([]string{"admin", "ceo"}))

		// 	staff.Get("/", h.ListUsers)
		// 	staff.Put("/{id}/ban", h.BanUser)
		// })
	})
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FieldsWhithRequestID(ctx)
	log.Info("Fetching user by username")

	req := GetRequest{Username: chi.URLParam(r, "username")}

	err := req.Validate()
	if err != nil {
		log.Error("Validation failed", zap.Error(err))
		shared.HTTP.WithValidationError(w, err)
		return
	}

	errorMessage := "Failed to get user"
	successMessage := "User retrieved successfully"

	var resp *GetResponse
	resp, err = h.service.GetUser(ctx, &req)
	if err != nil {
		log.Error(errorMessage, zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusNotFound, errorMessage, "User not found")
		return
	}

	log.Info(successMessage)
	shared.HTTP.WithSuccess(w, http.StatusOK, successMessage, resp)
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FieldsWhithRequestID(ctx)
	log.Info("Fetching user profile")

	userID := shared.Context.GetUserID(ctx)

	errorMessage := "Failed to get profile"
	successMessage := "Profile retrieved successfully"

	resp, err := h.service.GetProfile(ctx, userID)
	if err != nil {
		log.Error(errorMessage, zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusNotFound, errorMessage, "User not found")
		return
	}

	log.Info(successMessage)
	shared.HTTP.WithSuccess(w, http.StatusOK, successMessage, resp)
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FieldsWhithRequestID(ctx)
	log.Info("Updating user")

	userID := shared.Context.GetUserID(ctx)

	var req UpdateRequest
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

	errorMessage := "Failed to update user"
	successMessage := "User updated successfully"

	var resp *UpdateResponse
	resp, err = h.service.UpdateProfile(ctx, userID, &req)
	if err != nil {
		log.Error(errorMessage, zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusInternalServerError, errorMessage, "Internal server error")
		return
	}

	log.Info(successMessage)
	shared.HTTP.WithSuccess(w, http.StatusOK, successMessage, resp)
}
