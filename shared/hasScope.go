package shared

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UserID      string   `json:"sub"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	Token       *jwt.Token
}

func HasScope(requriedScope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Pull claims out of context
			claims, ok := r.Context().Value(UserClaimsKey).(*UserClaims)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if the required scope exists in the slice
			hasPermission := false
			for _, s := range claims.Permissions {
				if s == requriedScope {
					hasPermission = true
					break
				}
			}
			if !hasPermission {
				http.Error(w, "Forbidden: Missing scope"+requriedScope, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
