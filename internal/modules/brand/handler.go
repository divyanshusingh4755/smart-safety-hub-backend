package brand

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type RestHandler struct {
	service   *BrandService
	validator *validator.Validate
}

func NewRestHandler(service *BrandService, validator *validator.Validate) *RestHandler {
	return &RestHandler{
		service:   service,
		validator: validator,
	}
}

func (h *RestHandler) CreateBrand(w http.ResponseWriter, r *http.Request) {
	var request BrandsRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.service.CreateBrand(r.Context(), request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RestHandler) UpdateBrand(w http.ResponseWriter, r *http.Request) {
	brandID := chi.URLParam(r, "id")

	if brandID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	var request BrandsRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.service.UpdateBrand(r.Context(), brandID, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RestHandler) DeleteBrand(w http.ResponseWriter, r *http.Request) {
	brandID := chi.URLParam(r, "id")

	if brandID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	response, err := h.service.DeleteBrand(r.Context(), brandID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RestHandler) GetBrandByID(w http.ResponseWriter, r *http.Request) {
	brandId := chi.URLParam(r, "id")

	if brandId == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	response, err := h.service.GetBrandByID(r.Context(), brandId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusBadRequest)
	}

}

func (h *RestHandler) GetAllBrand(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page := query.Get("page")
	limit := query.Get("limit")

	if page == "" || limit == "" {
		http.Error(w, "Page and Limit is missing in params", http.StatusBadRequest)
		return
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return

	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	response, err := h.service.GetAllBrand(r.Context(), pageInt, limitInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusBadRequest)
	}

}
