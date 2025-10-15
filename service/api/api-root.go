package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getApiRoot returns basic API information and status
func (rt *_router) getApiRoot(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Create API information response
	apiInfo := map[string]interface{}{
		"name":      "WASAText API",
		"version":   "1.2.0",
		"status":    "running",
		"endpoints": 27, // Total number of endpoints including this one
	}

	// Set content type header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode and send response
	if err := json.NewEncoder(w).Encode(apiInfo); err != nil {
		ctx.Logger.WithError(err).Error("failed to encode API root response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
