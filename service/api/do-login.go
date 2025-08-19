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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(users) > 0 {
		user = users[0]
	} else {
		user, err = rt.db.SetMyUserName(req.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	resp := map[string]string{"identifier": user.UId}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// User represents a user in the system
type User struct {
	UId      string `json:"identifier"`
	Username string `json:"name"`
}
