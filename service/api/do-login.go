package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

// doLogin handles user login
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	type loginRequest struct {
		Name string `json:"name"`
	}
	var req loginRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	if len(req.Name) < 3 || len(req.Name) > 16 {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	// Try to find user by name
	users, err := rt.db.ListUsers(req.Name)
	var user database.User
	if err != nil {
		rt.sysLogger.LogError("Database error during user lookup: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(users) > 0 {
		user = users[0]
		rt.sysLogger.LogInfo("User " + req.Name + " logged in successfully")
	} else {
		user, err = rt.db.SetMyUserName(req.Name)
		if err != nil {
			rt.sysLogger.LogError("Failed to create new user " + req.Name + ": " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rt.sysLogger.LogInfo("New user " + req.Name + " registered and logged in")
	}

	resp := map[string]string{"identifier": user.UId}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		ctx.Logger.WithError(err).Error("failed to encode login response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// User represents a user in the system
type User struct {
	UId      string `json:"identifier"`
	Username string `json:"name"`
}
