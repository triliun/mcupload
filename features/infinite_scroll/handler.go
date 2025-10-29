// handler/resource_pack_handler.go
package infinite_scroll

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// GET /api/resource-packs?limit=20
// GET /api/resource-packs?cursor=a1b2c3d4-...&limit=20
// GET /api/resource-packs?author_id=author-uuid-here&limit=20

// -- Pastikan ada index untuk performa infinite scroll
// CREATE INDEX idx_resource_packs_status_created_at ON resource_packs(status, created_at DESC);
// CREATE INDEX idx_resource_packs_author_status_created_at ON resource_packs(author_id, status, created_at DESC);

// func SetupResourcePackRoutes(router chi.Router, db *sqlx.DB) {
// 	repo := repository.NewResourcePackRepository(db)
// 	handler := handler.NewResourcePackHandler(repo)

// 	router.Route("/api/resource-packs", func(r chi.Router) {
// 		r.Get("/", handler.ListResourcePacks)      // GET /api/resource-packs?cursor=&limit=20
// 		r.Get("/{id}", handler.GetResourcePack)    // GET /api/resource-packs/{id}
// 	})
// }

type ResourcePackHandler struct {
	repo Repository
}

func NewResourcePackHandler(repo Repository) *ResourcePackHandler {
	return &ResourcePackHandler{repo: repo}
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func (h *ResourcePackHandler) ListResourcePacks(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	params := ListResourcePacksParams{
		Cursor:   r.URL.Query().Get("cursor"),
		Status:   r.URL.Query().Get("status"),
		AuthorID: r.URL.Query().Get("author_id"),
	}

	// Parse limit dengan default 20 dan max 100
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			params.Limit = limit
		}
	}

	// Validate cursor format jika ada
	if params.Cursor != "" {
		if _, err := uuid.Parse(params.Cursor); err != nil {
			h.sendErrorResponse(w, "Invalid cursor format", http.StatusBadRequest)
			return
		}
	}

	// Validate author_id format jika ada
	if params.AuthorID != "" {
		if _, err := uuid.Parse(params.AuthorID); err != nil {
			h.sendErrorResponse(w, "Invalid author_id format", http.StatusBadRequest)
			return
		}
	}

	// Get resource packs from repository
	result, err := h.repo.ListResourcePacks(r.Context(), params)
	if err != nil {
		h.sendErrorResponse(w, "Failed to fetch resource packs", http.StatusInternalServerError)
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, result, nil)
}

func (h *ResourcePackHandler) GetResourcePack(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.sendErrorResponse(w, "Invalid resource pack ID", http.StatusBadRequest)
		return
	}

	resourcePack, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		h.sendErrorResponse(w, "Resource pack not found", http.StatusNotFound)
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, resourcePack, nil)
}

func (h *ResourcePackHandler) sendSuccessResponse(w http.ResponseWriter, status int, data interface{}, meta interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

func (h *ResourcePackHandler) sendErrorResponse(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Error:   message,
	})
}
