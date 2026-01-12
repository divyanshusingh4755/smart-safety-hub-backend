package user

import (
	"net/http"
	"strings"
)

func JWTMiddleware(jm *JwtManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, "Missing Token", 401)
				return
			}
			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "Inavlid auth token", 401)
				return
			}

			if _, err := jm.Verify(parts[1]); err != nil {
				http.Error(w, "Invalid Token", 401)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
