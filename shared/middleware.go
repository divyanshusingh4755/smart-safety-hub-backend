package shared

import (
	"context"
	"log"
	"net/http"
	"strings"
)

type contextKey string

const UserClaimsKey contextKey = "userClaims"

func JWTMiddleware(jm *JwtManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			log.Println("auth", auth)
			if auth == "" {
				http.Error(w, "Missing Token", 401)
				return
			}
			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "Inavlid auth token", 401)
				return
			}
			claims, err := jm.Verify(parts[1])
			log.Println("cleaimss", claims)
			if err != nil {
				http.Error(w, "Invalid Token", 401)
				return
			}
			ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
