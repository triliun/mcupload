package category

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
	router.Route("/category", func(categoryGroup chi.Router) {
		categoryGroup.Use(h.middleware.Validate)

		categoryGroup.Get("/", h.GetAllCategories)

		categoryGroup.Group(func(staff chi.Router) {
			staff.Use(h.middleware.ValidateRoles([]string{"admin", "ceo"}))

			staff.Post("/", h.CreateCategory)
			staff.Put("/", h.UpdateCategory)
			staff.Delete("/{id}", h.DeleteCategory)
		})
	})
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FieldsWhithRequestID(ctx)
	log.Info("Creating category")

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

	errorMessage := "Failed to create category"
	successMessage := "Category created successfully"

	var resp *CreateResponse
	resp, err = h.service.CreateCategory(ctx, &req)
	if err != nil {
		log.Error(errorMessage, zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusInternalServerError, errorMessage, "Internal server error")
		return
	}

	log.Info(successMessage)
	shared.HTTP.WithSuccess(w, http.StatusCreated, successMessage, resp)
}

func (h *Handler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FieldsWhithRequestID(ctx)
	log.Info("Fetching category")

	errorMessage := "Failed to get category"
	successMessage := "Category retrieved successfully"

	resp, err := h.service.GetAllCategories(ctx)
	if err != nil {
		log.Error(errorMessage, zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusNotFound, errorMessage, "Category not found")
		return
	}

	log.Info(successMessage)
	shared.HTTP.WithSuccess(w, http.StatusOK, successMessage, resp)
}

func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FieldsWhithRequestID(ctx)
	log.Info("Updating category")

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

	errorMessage := "Failed to update category"
	successMessage := "Category updated successfully"

	var resp *UpdateResponse
	resp, err = h.service.UpdateCategory(ctx, &req)
	if err != nil {
		log.Error(errorMessage, zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusInternalServerError, errorMessage, "Internal server error")
		return
	}

	log.Info(successMessage)
	shared.HTTP.WithSuccess(w, http.StatusOK, successMessage, resp)
}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FieldsWhithRequestID(ctx)
	log.Info("Deleting category")

	categoryID, err := shared.Param.GetUUID(r, "id")

	errorMessage := "Failed to delete category"
	successMessage := "Category deleted successfully"

	if err != nil {
		log.Error("Invalid ID format", zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusBadRequest, errorMessage, "Invalid ID format")
		return
	}

	var resp *DeleteResponse
	resp, err = h.service.DeleteCategory(ctx, categoryID)
	if err != nil {
		log.Error(errorMessage, zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusInternalServerError, errorMessage, "Internal server error")
		return
	}

	log.Info(successMessage)
	shared.HTTP.WithSuccess(w, http.StatusOK, successMessage, resp)
}
