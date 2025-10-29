package minecraft_version

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
	router.Route("/minecraft-version", func(versionGroup chi.Router) {
		versionGroup.Use(h.middleware.Validate)

		versionGroup.Get("/", h.GetAllVersions)

		versionGroup.Group(func(staff chi.Router) {
			staff.Use(h.middleware.ValidateRoles([]string{"ceo", "admin"}))

			staff.Post("/", h.CreateVersion)
			staff.Put("/", h.UpdateVersion)
			staff.Delete("/{id}", h.DeleteVersion)
		})
	})
}

func (h *Handler) CreateVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FieldsWhithRequestID(ctx)
	log.Info("Creating minecraft version")

	var req CreateRequest
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

	errorMessage := "Failed to create minecraft version"
	successMessage := "Minecraft version created successfully"

	var resp *CreateResponse
	resp, err = h.service.CreateVersion(ctx, &req)
	if err != nil {
		log.Error(errorMessage, zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusInternalServerError, errorMessage, "Internal server error")
		return
	}

	log.Info(successMessage)
	shared.HTTP.WithSuccess(w, http.StatusCreated, successMessage, resp)
}

func (h *Handler) GetAllVersions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FieldsWhithRequestID(ctx)
	log.Info("Fetching minecraft version")

	errorMessage := "Failed to get minecraft version"
	successMessage := "Minecraft version retrieved successfully"

	resp, err := h.service.GetAllVersions(ctx)
	if err != nil {
		log.Error(errorMessage, zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusNotFound, errorMessage, "Minecraft version not found")
		return
	}

	log.Info(successMessage)
	shared.HTTP.WithSuccess(w, http.StatusOK, successMessage, resp)
}

func (h *Handler) UpdateVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FieldsWhithRequestID(ctx)
	log.Info("Updating minecraft version")

	var req UpdateRequest
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

	errorMessage := "Failed to minecraft version"
	successMessage := "Minecraft version updated successfully"

	var resp *UpdateResponse
	resp, err = h.service.UpdateVersion(ctx, &req)
	if err != nil {
		log.Error(errorMessage, zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusInternalServerError, errorMessage, "Internal server error")
		return
	}

	log.Info(successMessage)
	shared.HTTP.WithSuccess(w, http.StatusOK, successMessage, resp)
}

func (h *Handler) DeleteVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FieldsWhithRequestID(ctx)
	log.Info("Deleting minecraft version")

	errorMessage := "Failed to delete minecraft version"
	successMessage := "Minecraft version deleted successfully"

	minecraftVersionID, err := shared.Param.GetUUID(r, "id")
	if err != nil {
		log.Error("Invalid ID format", zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusBadRequest, errorMessage, "Invalid ID format")
		return
	}

	var resp *DeleteResponse
	resp, err = h.service.DeleteVersion(ctx, minecraftVersionID)
	if err != nil {
		log.Error(errorMessage, zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusInternalServerError, errorMessage, "Internal server error")
		return
	}

	log.Info(successMessage)
	shared.HTTP.WithSuccess(w, http.StatusOK, successMessage, resp)
}
