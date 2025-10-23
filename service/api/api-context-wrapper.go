package api

import (
	"context"
	"encoding/json"
	"strings"

	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// httpRouterHandler is the signature for functions that accepts a reqcontext.RequestContext in addition to those
// required by the httprouter package.
type httpRouterHandler func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext)

// wrap parses the request and adds a reqcontext.RequestContext instance related to the request.
func (rt *_router) wrap(fn httpRouterHandler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		reqUUID, err := uuid.NewV4()
		if err != nil {
			rt.baseLogger.WithError(err).Error("can't generate a request UUID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var ctx = reqcontext.RequestContext{
			ReqUUID: reqUUID,
		}

		// Create a request-specific logger
		ctx.Logger = rt.baseLogger.WithFields(logrus.Fields{
			"reqid":     ctx.ReqUUID.String(),
			"remote-ip": r.RemoteAddr,
		})

		// Call the next handler in chain (usually, the handler function for the path)
		fn(w, r, ps, ctx)
	}
}

// wrapAuth combines authentication middleware with request context wrapping
func (rt *_router) wrapAuth(fn httpRouterHandler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Extract Authorization header and validate
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

		// Create request context
		reqUUID, genErr := uuid.NewV4()
		if genErr != nil {
			rt.baseLogger.WithError(genErr).Error("can't generate a request UUID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var ctx = reqcontext.RequestContext{
			ReqUUID: reqUUID,
		}

		// Create a request-specific logger
		ctx.Logger = rt.baseLogger.WithFields(logrus.Fields{
			"reqid":     ctx.ReqUUID.String(),
			"remote-ip": r.RemoteAddr,
			"user-id":   user.UId,
		})

		// Add user to request context
		r = r.WithContext(context.WithValue(r.Context(), AuthUserKey, user))

		// Call the actual handler
		fn(w, r, ps, ctx)
	}
}

// sendError sends a JSON error response
func (rt *_router) sendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := map[string]string{
		"message": message,
	}

	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		rt.baseLogger.WithError(err).Error("Error encoding error response")
	}
}
