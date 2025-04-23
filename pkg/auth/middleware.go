package auth

import (
	"context"
	"net/http"
	"strings"
)

// Clé pour stocker l'ID utilisateur dans le contexte
type contextKey string

const UserIDKey contextKey = "userID"

/**
 * Vérifie la présence et la validité du token dans l'en-tête Authorization
 */
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"code":401,"error_code":"no_authorization","msg":"Missing token"}`, http.StatusUnauthorized)
			return
		}

		// On retire "Bearer" situé devan le token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // Si "Bearer " n'était pas présent
			http.Error(w, `{"code":401,"error_code":"no_authorization","msg":"Invalid token format"}`, http.StatusUnauthorized)
			return
		}

		// Vérification du token
		claims, err := VerifyToken(tokenString)
		if err != nil {
			http.Error(w, `{"code":401,"error_code":"no_authorization","msg":"Invalid token"}`, http.StatusUnauthorized)
			return
		}

		// Stocker l'ID utilisateur dans le contexte
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
