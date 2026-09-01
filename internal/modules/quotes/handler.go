package quotes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type RestHandler struct {
	service   *QuoteService
	validator *validator.Validate
}

func NewRestHandler(service *QuoteService, validator *validator.Validate) *RestHandler {
	return &RestHandler{
		service:   service,
		validator: validator,
	}
}

func (h *RestHandler) CreateQuote(w http.ResponseWriter, r *http.Request) {
	var request CreateQuoteRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.service.CreateQuote(r.Context(), request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RestHandler) GetQuoteByID(w http.ResponseWriter, r *http.Request) {
	quoteID := chi.URLParam(r, "id")
	if quoteID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	response, err := h.service.GetQuoteByID(r.Context(), quoteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RestHandler) UpdateQuoteStatus(w http.ResponseWriter, r *http.Request) {
	quoteID := chi.URLParam(r, "id")

	if quoteID == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	var request UpdateQuoteStatusDTO
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.service.UpdateQuoteStatus(r.Context(), quoteID, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode request", http.StatusInternalServerError)
	}
}

func (h *RestHandler) GetAllQuotes(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 20

	if pageParam := r.URL.Query().Get("page"); pageParam != "" {
		parsedPage, err := strconv.Atoi(pageParam)

		if err != nil || parsedPage < 1 {
			http.Error(w, "invalid page", http.StatusBadRequest)
			return
		}
		page = parsedPage
	}

	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		parsedLimit, err := strconv.Atoi(limitParam)

		if err != nil || parsedLimit < 1 || parsedLimit > 100 {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}

		limit = parsedLimit
	}

	request := GetQuotesQueryDTO{
		Page:   page,
		Limit:  limit,
		Search: strings.TrimSpace(r.URL.Query().Get("search")),
		Status: strings.TrimSpace(r.URL.Query().Get("status")),
	}

	if err := h.validator.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	repsone, err := h.service.GetAllQuotes(r.Context(), request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(repsone); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
