package categories

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type RestHandler struct {
	service   *CategoryService
	validator *validator.Validate
}

func NewRestHandler(service *CategoryService, validator *validator.Validate) *RestHandler {
	return &RestHandler{
		service:   service,
		validator: validator,
	}
}

func (h *RestHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var request CategoryRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.service.CreateCategory(r.Context(), request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RestHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "id")

	if categoryID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	var request CategoryRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.service.UpdateCategory(r.Context(), categoryID, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RestHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "id")

	if categoryID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	response, err := h.service.DeleteCategory(r.Context(), categoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RestHandler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	categoryId := chi.URLParam(r, "id")

	if categoryId == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	response, err := h.service.GetCategoryByID(r.Context(), categoryId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusBadRequest)
	}

}

func (h *RestHandler) GetAllCategory(w http.ResponseWriter, r *http.Request) {
	response, err := h.service.GetAllCategory(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusBadRequest)
	}

}
