package blogs

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type RestHandler struct {
	service   *BlogService
	validator *validator.Validate
}

func NewRestHandler(service *BlogService, validator *validator.Validate) *RestHandler {
	return &RestHandler{
		service:   service,
		validator: validator,
	}
}

func (h *RestHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	var request BlogRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.service.CreateBlog(r.Context(), request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RestHandler) UpdateBlog(w http.ResponseWriter, r *http.Request) {
	blogID := chi.URLParam(r, "id")
	if blogID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	var request UpdateBlogDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.service.UpdateBlog(r.Context(), blogID, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RestHandler) DeleteBlog(w http.ResponseWriter, r *http.Request) {
	blogID := chi.URLParam(r, "id")
	if blogID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	response, err := h.service.DeleteBlog(r.Context(), blogID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RestHandler) GetBlogByID(w http.ResponseWriter, r *http.Request) {
	blogID := chi.URLParam(r, "id")
	if blogID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	response, err := h.service.GetBlogByID(r.Context(), blogID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RestHandler) GetBlogBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		http.Error(w, "slug is required", http.StatusBadRequest)
		return
	}

	response, err := h.service.GetBlogBySlug(r.Context(), slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RestHandler) GetAllBlogs(w http.ResponseWriter, r *http.Request) {
	response, err := h.service.GetAllBlogs(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RestHandler) GetPublishedBlogs(w http.ResponseWriter, r *http.Request) {
	response, err := h.service.GetPublishedBlogs(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RestHandler) GetPublishedPages(w http.ResponseWriter, r *http.Request) {
	response, err := h.service.GetPublishedPages(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
