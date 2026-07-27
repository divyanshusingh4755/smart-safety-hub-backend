package contacts

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type RestHandler struct {
	service   *ContactService
	validator *validator.Validate
}

func NewRestHandler(
	servcie *ContactService,
	validator *validator.Validate,
) *RestHandler {
	return &RestHandler{
		service:   servcie,
		validator: validator,
	}
}

func (h *RestHandler) CreateContact(
	w http.ResponseWriter,
	r *http.Request,
) {

	var request CreateContactRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	response, err := h.service.CreateContact(
		r.Context(),
		request,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)
	}
}

func (h *RestHandler) GetContactByID(
	w http.ResponseWriter,
	r *http.Request,
) {

	contactID := chi.URLParam(r, "id")

	if contactID == "" {
		http.Error(
			w,
			"ID is required",
			http.StatusBadRequest,
		)
		return
	}

	response, err := h.service.GetContactByID(
		r.Context(),
		contactID,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)
	}
}

func (h *RestHandler) GetAllContacts(
	w http.ResponseWriter,
	r *http.Request,
) {

	page := 1
	limit := 20

	if pageParam := r.URL.Query().Get("page"); pageParam != "" {
		parsedPage, err := strconv.Atoi(pageParam)

		if err != nil || parsedPage < 1 {
			http.Error(
				w,
				"invalid page",
				http.StatusBadRequest,
			)
			return
		}

		page = parsedPage
	}

	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		parsedLimit, err :=
			strconv.Atoi(limitParam)

		if err != nil ||
			parsedLimit < 1 ||
			parsedLimit > 100 {

			http.Error(
				w,
				"invalid limit",
				http.StatusBadRequest,
			)
			return
		}

		limit = parsedLimit
	}

	request := GetContactsQueryDTO{
		Page: page,

		Limit: limit,

		Search: strings.TrimSpace(
			r.URL.Query().Get("search"),
		),

		Status: strings.TrimSpace(
			r.URL.Query().Get("status"),
		),

		Source: strings.TrimSpace(
			r.URL.Query().Get("source"),
		),
	}

	if err := h.validator.Struct(request); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	response, err :=
		h.service.GetAllContacts(
			r.Context(),
			request,
		)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	if err :=
		json.NewEncoder(w).
			Encode(response); err != nil {

		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)
	}
}

func (h *RestHandler) UpdateContactStatus(
	w http.ResponseWriter,
	r *http.Request,
) {

	contactID := chi.URLParam(r, "id")

	if contactID == "" {
		http.Error(
			w,
			"ID is required",
			http.StatusBadRequest,
		)
		return
	}

	var request UpdateContactStatusDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	response, err := h.service.UpdateContactStatus(
		r.Context(),
		contactID,
		request,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)
	}
}
