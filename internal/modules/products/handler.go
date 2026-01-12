package products

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type RestHandler struct {
	service   *ProductService
	validator *validator.Validate
}

func NewRestHandler(service *ProductService, validator *validator.Validate) *RestHandler {
	return &RestHandler{
		service:   service,
		validator: validator,
	}
}

func (h *RestHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var request ProductRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.service.CreateProduct(r.Context(), request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RestHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")

	if productID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	var request ProductRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.service.UpdateProduct(r.Context(), productID, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RestHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")

	if productID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	query := r.URL.Query()
	status := query.Get("status")

	if status == "" {
		http.Error(w, "Status is required", http.StatusBadRequest)
		return
	}

	response, err := h.service.DeleteProduct(r.Context(), productID, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RestHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")

	if productID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	response, err := h.service.GetProductByID(r.Context(), productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusBadRequest)
	}
}

func (h *RestHandler) GetProductBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	if slug == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	response, err := h.service.GetProductBySlug(r.Context(), slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusBadRequest)
	}
}

func (h *RestHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	request := ProductFilters{
		Category: query["category"],
		Brand:    query["brand"],
		Search:   query.Get("search"),
		Status:   query.Get("status"),
		Page:     1,
		Limit:    40,
	}

	if p, err := strconv.Atoi(query.Get("page")); err == nil && p > 0 {
		request.Page = p
	}

	if l, err := strconv.Atoi(query.Get("limit")); err == nil && l > 0 {
		request.Limit = l
	}

	response, err := h.service.GetAllProducts(r.Context(), request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusBadRequest)
	}
}

func (h *RestHandler) AddProductAttribute(w http.ResponseWriter, r *http.Request) {
	var request ProductAttributeDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.service.AddProductAttribute(r.Context(), request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RestHandler) GetProductAttributeByID(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")

	if productID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	response, err := h.service.GetProductAttributeByID(r.Context(), productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusBadRequest)
	}
}

func (h *RestHandler) SyncProductVariants(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	var request VariantRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if productID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	response, err := h.service.SyncProductVariants(r.Context(), productID, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusBadRequest)
	}
}

func (h *RestHandler) AddProductMedia(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")

	if productID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	var request []ProductMediaDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Var(request, "required,dive"); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.service.AddProductMedia(r.Context(), productID, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusBadRequest)
	}
}

func (h *RestHandler) GetProductMedia(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")

	if productID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	response, err := h.service.GetProductMedia(r.Context(), productID)
	if err != nil {
		fmt.Println("err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusBadRequest)
	}
}

func (h *RestHandler) GetProductVariants(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")

	if productID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	response, err := h.service.GetProductVariants(r.Context(), productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusBadRequest)
	}
}

func (h *RestHandler) SaveProductSEO(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	if productID == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	var request ProductSEODTO
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the DTO
	if err := h.validator.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.SaveProductSEO(r.Context(), productID, request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *RestHandler) GetProductSEO(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	if productID == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	response, err := h.service.GetProductSEO(r.Context(), productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
