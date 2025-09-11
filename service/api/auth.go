package api

import (
	"context"
	"net/http"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/gofrs/uuid"
)

// Context key for authenticated user
type contextKey string

const AuthUserKey contextKey = "auth_user"

// AuthMiddleware validates Bearer token and adds user to context
func (rt *_router) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			rt.sendError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		// Check Bearer format
		if !strings.HasPrefix(authHeader, "Bearer ") {
			rt.sendError(w, http.StatusUnauthorized, "Invalid authorization format. Use 'Bearer <token>'")
			return
		}

		// Extract token (user ID)
		token := strings.TrimPrefix(authHeader, "Bearer ")
		token = strings.TrimSpace(token)

		// Validate UUID format
		userID, err := uuid.FromString(token)
		if err != nil {
			rt.sendError(w, http.StatusUnauthorized, "Invalid token format")
			return
		}

		// Verify user exists in database
		user, err := rt.db.GetUserByID(userID.String())
		if err != nil {
			rt.sendError(w, http.StatusUnauthorized, "Invalid token - user not found")
			return
		}

		// Add user to request context
		ctx := context.WithValue(r.Context(), AuthUserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// GetAuthenticatedUser extracts the authenticated user from request context
func GetAuthenticatedUser(r *http.Request) (database.User, bool) {
	user, ok := r.Context().Value(AuthUserKey).(database.User)
	return user, ok
}

// RequireAuth is a helper to check if user is authenticated and matches path parameter
func (rt *_router) RequireAuth(w http.ResponseWriter, r *http.Request, pathUserID string) (database.User, bool) {
	user, ok := GetAuthenticatedUser(r)
	if !ok {
		rt.sendError(w, http.StatusUnauthorized, "Authentication required")
		return database.User{}, false
	}

	// Validate that path user ID matches authenticated user
	if pathUserID != "" && user.UId != pathUserID {
		rt.sendError(w, http.StatusForbidden, "Access denied - can only access your own resources")
		return database.User{}, false
	}

	return user, true
}
