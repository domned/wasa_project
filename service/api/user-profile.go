package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get user ID from URL
	userId := ps.ByName("id")
	if userId == "" {
		http.Error(w, "missing user id", http.StatusBadRequest)
		return
	}

	var name string
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&name); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if len(name) < 3 || len(name) > 16 {
		http.Error(w, "username must be 3-16 characters", http.StatusBadRequest)
		return
	}

	// Check if username already exists (except for current user)
	var existingUserId string
	err := rt.db.GetRawDB().QueryRow("SELECT id FROM users WHERE username = ? AND id != ?", name, userId).Scan(&existingUserId)
	if err == nil {
		http.Error(w, "username already in use", http.StatusBadRequest)
		return
	}

	// Update username
	_, err = rt.db.GetRawDB().Exec("UPDATE users SET username = ? WHERE id = ?", name, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get updated user details
	var picture string
	err = rt.db.GetRawDB().QueryRow("SELECT COALESCE(picture, '') FROM users WHERE id = ?", userId).Scan(&picture)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := database.User{UId: userId, Username: name, Picture: picture}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		ctx.Logger.WithError(err).Error("failed to encode user response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get user ID from URL
	userId := ps.ByName("id")
	if userId == "" {
		http.Error(w, "missing user id", http.StatusBadRequest)
		return
	}

	var photo string
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&photo); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	// Update picture
	_, err := rt.db.GetRawDB().Exec("UPDATE users SET picture = ? WHERE id = ?", photo, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get updated user details
	var username string
	err = rt.db.GetRawDB().QueryRow("SELECT username FROM users WHERE id = ?", userId).Scan(&username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := database.User{UId: userId, Username: username, Picture: photo}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		ctx.Logger.WithError(err).Error("failed to encode user response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
