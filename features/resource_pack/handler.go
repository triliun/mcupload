package resource_pack

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
	router.Route("/resource-pack", func(packGroup chi.Router) {
		packGroup.Get("/{slug}", h.GetResourcePack)
		packGroup.Get("/list", h.ListResourcePacks)

		packGroup.Group(func(protected chi.Router) {
			protected.Use(h.middleware.Validate)
			protected.Post("/", h.CreateResourcePack)
			protected.Get("/my", h.GetAllMyResourcePack)
			protected.Get("/my/{id}", h.GetMyResourcePack)
			protected.Put("/{id}", h.UpdateResourcePack)
			protected.Delete("/{id}", h.DeleteResourcePack)
		})
	})
}

func (h *Handler) GetResourcePack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FieldsWhithRequestID(ctx)
	log.Info("Fetching resource pack")

	slug := chi.URLParam(r, "slug")

	errorMessage := "Failed to get resource pack"
	successMessage := "Resource pack retrieved successfully"

	req := GetRequest{Slug: slug}
	err := req.Validate()
	if err != nil {
		log.Error("Validation failed", zap.Error(err))
		shared.HTTP.WithValidationError(w, err)
		return
	}

	var resp *GetResponse
	resp, err = h.service.GetResourcePack(ctx, req)
	if err != nil {
		log.Error(errorMessage, zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusNotFound, errorMessage, "Resource pack not found")
		return
	}

	log.Info(successMessage)
	shared.HTTP.WithSuccess(w, http.StatusOK, successMessage, resp)
}

func (h *Handler) CreateResourcePack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FieldsWhithRequestID(ctx)
	log.Info("Creating resource pack")

	userID := shared.Context.GetUserID(ctx)

	errorMessage := "Failed to create resource pack"
	successMessage := "Resource pack created successfully"

	resp, err := h.service.CreateResourcePack(ctx, userID)
	if err != nil {
		log.Error(errorMessage, zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusInternalServerError, errorMessage, "Internal server error")
		return
	}

	log.Info(successMessage)
	shared.HTTP.WithSuccess(w, http.StatusCreated, successMessage, resp)
}

func (h *Handler) GetAllMyResourcePack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FieldsWhithRequestID(ctx)
	log.Info("Fetching all my resource pack")

	errorMessage := "Failed to get all my resource pack"
	successMessage := "Resource pack retrieved successfully"

	userID := shared.Context.GetUserID(ctx)

	var resp []GetAllMyResponse
	resp, err := h.service.GetAllMyResourcePack(ctx, userID)
	if err != nil {
		log.Error(errorMessage, zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusNotFound, errorMessage, "Resource pack not found")
		return
	}

	log.Info(successMessage)
	shared.HTTP.WithSuccess(w, http.StatusOK, successMessage, resp)
}

func (h *Handler) GetMyResourcePack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FieldsWhithRequestID(ctx)
	log.Info("Fetching my resource pack")

	errorMessage := "Failed to get my resource pack"
	successMessage := "Resource pack retrieved successfully"

	userID := shared.Context.GetUserID(ctx)
	packID, err := shared.Param.GetUUID(r, "id")
	if err != nil {
		log.Error("Invalid ID format", zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusBadRequest, errorMessage, "Invalid ID format")
		return
	}

	var resp *GetMyResponse
	resp, err = h.service.GetMyResourcePack(ctx, userID, packID)
	if err != nil {
		log.Error(errorMessage, zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusNotFound, errorMessage, "Resource pack not found")
		return
	}

	log.Info(successMessage)
	shared.HTTP.WithSuccess(w, http.StatusOK, successMessage, resp)
}

func (h *Handler) ListResourcePacks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FieldsWhithRequestID(ctx)
	log.Info("Fetching list resource packs")

	req := ListRequest{
		InfiniteScrollRequest: shared.InfiniteScrollRequest{
			Cursor: shared.Query.Get(r, "cursor"),
			Limit:  shared.Query.GetInt(r, "limit"),
		},
		Status:   shared.Query.Get(r, "status"),
		AuthorID: shared.Query.Get(r, "author_id"),
	}

	err := req.Validate()
	if err != nil {
		log.Error("Validation failed", zap.Error(err))
		shared.HTTP.WithValidationError(w, err)
		return
	}

	// Parse limit dengan default 20 dan max 100
	// if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
	// 	if limit, err := strconv.Atoi(limitStr); err == nil {
	// 		req.Limit = limit
	// 	}
	// }

	// // Validate cursor format jika ada
	// if req.Cursor != "" {
	// 	if _, err := uuid.Parse(req.Cursor); err != nil {
	// 		h.sendErrorResponse(w, "Invalid cursor format", http.StatusBadRequest)
	// 		return
	// 	}
	// }

	// // Validate author_id format jika ada
	// if req.AuthorID != "" {
	// 	if _, err := uuid.Parse(req.AuthorID); err != nil {
	// 		h.sendErrorResponse(w, "Invalid author_id format", http.StatusBadRequest)
	// 		return
	// 	}
	// }

	errorMessage := "Failed to get list resource packs"
	successMessage := "List Resource pack retrieved successfully"

	var result *ListResponse
	result, err = h.service.ListResourcePacks(ctx, &req)
	if err != nil {
		log.Error(errorMessage, zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusInternalServerError, errorMessage, "Internal server error")
		return
	}

	log.Info(successMessage)
	shared.HTTP.WithMetadataSuccess(w, http.StatusOK, successMessage, result.List, result.NextCursor, result.HasMore)
}

func (h *Handler) UpdateResourcePack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FieldsWhithRequestID(ctx)
	log.Info("Updating resource pack")

	errorMessage := "Failed to update resource pack"
	successMessage := "Resource pack updated successfully"

	userID := shared.Context.GetUserID(ctx)
	packID, err := shared.Param.GetUUID(r, "id")
	if err != nil {
		log.Error("Invalid ID format", zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusBadRequest, errorMessage, "Invalid ID format")
		return
	}

	req := UpdateRequest{AuthorID: userID, PackID: packID}
	err = shared.Binder.BindJSON(r, &req)
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

	var resp *UpdateResponse
	resp, err = h.service.UpdateResourcePack(ctx, &req)
	if err != nil {
		log.Error(errorMessage, zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusInternalServerError, errorMessage, "Internal server error")
		return
	}

	log.Info(successMessage)
	shared.HTTP.WithSuccess(w, http.StatusOK, successMessage, resp)
}

func (h *Handler) DeleteResourcePack(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := logger.FieldsWhithRequestID(ctx)
	log.Info("Deleting resource pack")

	errorMessage := "Failed to delete resource pack"
	successMessage := "Resource pack deleted successfully"

	userID := shared.Context.GetUserID(ctx)
	packID, err := shared.Param.GetUUID(r, "id")
	if err != nil {
		log.Error("Invalid ID format", zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusBadRequest, errorMessage, "Invalid ID format")
		return
	}

	var resp *DeleteResponse
	resp, err = h.service.Delete(ctx, userID, packID)
	if err != nil {
		log.Error(errorMessage, zap.Error(err))
		shared.HTTP.WithGeneralError(w, http.StatusInternalServerError, errorMessage, "Internal server error")
		return
	}

	log.Info(successMessage)
	shared.HTTP.WithSuccess(w, http.StatusOK, successMessage, resp)
}
