package social

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/markbates/goth/gothic"
	"github.com/smart-safety-hub/backend/shared"
)

type RestHandler struct {
	service   *SocialService
	validator *validator.Validate
}

func NewRestHandler(service *SocialService, validator *validator.Validate) *RestHandler {
	return &RestHandler{
		service:   service,
		validator: validator,
	}
}

func getBackendUrl(r *http.Request) string {
	scheme := "https"
	if r.TLS == nil {
		scheme = "http"
	}
	return fmt.Sprintf("%s://%s", scheme, r.Host)
}

func (h *RestHandler) PrepareAuth(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(shared.UserClaimsKey).(*shared.UserClaims)

	provider := chi.URLParam(r, "provider")

	// Save UserID in the session now while we have the JWT header
	session, _ := gothic.Store.Get(r, "gothic_session")
	session.Values["user_id"] = claims.UserID

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	// Returning the starting url
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"redirectUrl": fmt.Sprintf("%s/v1/auth/social/%s", getBackendUrl(r), provider),
	})
}

func (h *RestHandler) BeginAuth(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	if provider == "" {
		http.Error(w, "provider is required", http.StatusBadRequest)
		return
	}

	session, _ := gothic.Store.Get(r, "gothic_session")
	if _, ok := session.Values["user_id"].(string); !ok {
		http.Error(w, "Unauthorized: Please state from the app", 401)
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))
	gothic.BeginAuthHandler(w, r)
}

func (h *RestHandler) Callback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	if provider == "" {
		http.Error(w, "provider is required", http.StatusBadRequest)
		return
	}

	// Add provider to context so gothic.CompleteUserAuth knows which provider to use
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	// Retrieve the UserID from the session
	session, _ := gothic.Store.Get(r, "gothic_session")
	userID, ok := session.Values["user_id"].(string)
	if !ok || userID == "" {
		http.Error(w, "User session lost. Please try again", http.StatusUnauthorized)
		return
	}

	gothUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.service.LinkExternalAccount(r.Context(), userID, gothUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	frontendURL := "http://localhost:3000/dashboard/social?success=true"
	http.Redirect(w, r, frontendURL, http.StatusTemporaryRedirect)
}

func (h *RestHandler) HandleListConnections(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(shared.UserClaimsKey).(*shared.UserClaims)
	userID := claims.UserID

	response, err := h.service.HandleListConnections(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *RestHandler) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(shared.UserClaimsKey).(*shared.UserClaims)
	if !ok || claims == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID := claims.UserID
	err := h.service.CreatePost(r.Context(), userID, req)
	if err != nil {
		fmt.Printf("Post failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post published successfully!"})
}
