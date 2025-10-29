package pack_category

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
	router.Route("/pack-category", func(packCategoryGroup chi.Router) {
		// packCategoryGroup.GET("/{packID}", h.GetAllPackCategoriesFromPackID)

		packCategoryGroup.Group(func(protected chi.Router) {
			protected.Use(h.middleware.Validate)

			protected.Post("/", h.CreatePackCategory)
			protected.Get("/", h.GetAllPackCategories)
			protected.Delete("/{id}", h.DeletePackCategory)
		})
	})
}

func (h *Handler) CreatePackCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FieldsWhithRequestID(ctx)
	log.Info("Creating pack category")

	errorMessage := "Failed to create pack category"
	successMessage := "Pack category created successfully"

	userID := shared.Context.GetUserID(ctx)

	req := CreateRequest{authorID: userID}
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

	var resp *CreateResponse
	resp, err = h.service.CreatePackCategory(ctx, &req)
	if err != nil {
		log.Error(errorMessage, zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusInternalServerError, errorMessage, "Internal server error")
		return
	}

	log.Info(successMessage)
	shared.HTTP.WithSuccess(w, http.StatusCreated, successMessage, resp)
}

func (h *Handler) GetAllPackCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FieldsWhithRequestID(ctx)
	log.Info("Fetching pack category")

	userID := shared.Context.GetUserID(ctx)

	req := GetAllRequest{authorID: userID}
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
	errorMessage := "Failed to get pack category"
	successMessage := "Pack category retrieved successfully"

	var resp []GetResponse
	resp, err = h.service.GetAllPackCategories(ctx, &req)
	if err != nil {
		log.Error(errorMessage, zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusNotFound, errorMessage, "Pack category not found")
		return
	}

	log.Info(successMessage)
	shared.HTTP.WithSuccess(w, http.StatusOK, successMessage, resp)
}

func (h *Handler) DeletePackCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FieldsWhithRequestID(ctx)
	log.Info("Deleting pack category")

	errorMessage := "Failed to delete pack category"
	successMessage := "Pack Category deleted successfully"

	userID := shared.Context.GetUserID(ctx)

	packCategoryID, err := shared.Param.GetUUID(r, "id")
	if err != nil {
		log.Error("Invalid ID format", zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusBadRequest, errorMessage, "Invalid ID format")
		return
	}

	var resp *DeleteResponse
	resp, err = h.service.DeletePackCategory(ctx, packCategoryID, userID)
	if err != nil {
		log.Error(errorMessage, zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusInternalServerError, errorMessage, "Internal server error")
		return
	}

	log.Info(successMessage)
	shared.HTTP.WithSuccess(w, http.StatusOK, successMessage, resp)
}
