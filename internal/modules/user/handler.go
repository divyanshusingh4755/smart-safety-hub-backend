package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type RestHandler struct {
	service   *UserService
	validator *validator.Validate
}

func NewRestHandler(service *UserService, validator *validator.Validate) *RestHandler {
	return &RestHandler{service: service, validator: validator}
}

func (h *RestHandler) Register(w http.ResponseWriter, r *http.Request) {
	var request RegisterDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.service.Register(r.Context(), &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// http.SetCookie(w, &http.Cookie{
	// 	Name:     "refresh_token",
	// 	Value:    refreshToken,
	// 	HttpOnly: true,
	// 	Secure:   false,
	// 	Path:     "/",
	// 	SameSite: http.SameSiteLaxMode,
	// 	Expires:  time.Now().Add(30 * 24 * time.Hour),
	// })

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RestHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request LoginDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Println("error here", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		fmt.Println("error there", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.service.Login(r.Context(), &request)
	if err != nil {
		fmt.Println("error where", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// if refreshToken == "" {
	// 	http.Error(w, "Failed to generate session", http.StatusInternalServerError)
	// 	return
	// }

	// http.SetCookie(w, &http.Cookie{
	// 	Name:     "refresh_token",
	// 	Value:    refreshToken,
	// 	HttpOnly: true,
	// 	Secure:   false,
	// 	Path:     "/",
	// 	SameSite: http.SameSiteLaxMode,
	// 	Expires:  time.Now().Add(30 * 24 * time.Hour),
	// })

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RestHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var request ForgotPasswordDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.service.ForgotPassword(r.Context(), &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RestHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var request ResetPasswordDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.service.ResetPassword(r.Context(), &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RestHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var request LogoutDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Println("rese", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.service.Logout(r.Context(), request)
	fmt.Println("response", response, err)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// http.SetCookie(w, &http.Cookie{
	// 	Name:     "refresh_token",
	// 	Value:    "",
	// 	Path:     "/",
	// 	Expires:  time.Unix(0, 0),
	// 	MaxAge:   -1,
	// 	HttpOnly: true,
	// 	Secure:   false,
	// })

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RestHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var request RefreshTokenDTO

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.service.RefreshToken(r.Context(), request)
	fmt.Println("err", err)
	if err != nil {
		http.Error(w, "Session expired", http.StatusUnauthorized)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
