package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/VishalHilal/e-commerce-api/internal/json"
)

type contextKey string

const UserContextKey = contextKey("user")

func (j *JWTService) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			json.Write(w, http.StatusUnauthorized, map[string]string{"error": "Authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			json.Write(w, http.StatusUnauthorized, map[string]string{"error": "Bearer token required"})
			return
		}

		claims, err := j.ValidateToken(tokenString)
		if err != nil {
			json.Write(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRole(allowedRoles ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserContextKey).(*JWTClaims)
			if !ok {
				json.Write(w, http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
				return
			}

			for _, role := range allowedRoles {
				if claims.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			json.Write(w, http.StatusForbidden, map[string]string{"error": "Insufficient permissions"})
		})
	}
}

func GetUserFromContext(ctx context.Context) *JWTClaims {
	if claims, ok := ctx.Value(UserContextKey).(*JWTClaims); ok {
		return claims
	}
	return nil
}
