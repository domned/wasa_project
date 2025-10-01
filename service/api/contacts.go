package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) addContact(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get user ID from URL
	userId := ps.ByName("id")
	if userId == "" {
		http.Error(w, "missing user id", http.StatusBadRequest)
		return
	}

	// Parse request body for contact user ID
	var request struct {
		ContactUserId string `json:"contactUserId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if request.ContactUserId == "" {
		http.Error(w, "contactUserId cannot be empty", http.StatusBadRequest)
		return
	}

	// Get user details
	var userUsername string
	err := rt.db.GetRawDB().QueryRow("SELECT username FROM users WHERE id = ?", userId).Scan(&userUsername)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// Find contact by user ID
	var contactId, contactUsername, contactPicture string
	err = rt.db.GetRawDB().QueryRow("SELECT id, username, COALESCE(picture, '') FROM users WHERE id = ?", request.ContactUserId).
		Scan(&contactId, &contactUsername, &contactPicture)
	if err != nil {
		http.Error(w, "contact not found", http.StatusNotFound)
		return
	}

	// Don't allow adding self as contact
	if contactId == userId {
		http.Error(w, "cannot add yourself as contact", http.StatusBadRequest)
		return
	}

	// Check if contact already exists
	var existingContactId string
	err = rt.db.GetRawDB().QueryRow("SELECT id FROM contacts WHERE user_id = ? AND contact_id = ?", userId, contactId).
		Scan(&existingContactId)
	if err == nil {
		http.Error(w, "contact already exists", http.StatusConflict)
		return
	}

	// Add contact
	user := database.User{UId: userId, Username: userUsername}
	contact := database.User{UId: contactId, Username: contactUsername, Picture: contactPicture}

	_, err = rt.db.AddContact(user, contact)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(contact); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (rt *_router) listContacts(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get user ID from URL
	userId := ps.ByName("id")
	if userId == "" {
		http.Error(w, "missing user id", http.StatusBadRequest)
		return
	}

	// Get user details
	var username string
	err := rt.db.GetRawDB().QueryRow("SELECT username FROM users WHERE id = ?", userId).Scan(&username)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// Get contacts list
	user := database.User{UId: userId, Username: username}
	contacts, err := rt.db.ListContacts(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(contacts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (rt *_router) removeContact(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get user ID and contact ID from URL
	userId := ps.ByName("id")
	contactId := ps.ByName("contactId")

	if userId == "" || contactId == "" {
		http.Error(w, "missing user id or contact id", http.StatusBadRequest)
		return
	}

	// Get user details
	var userUsername string
	err := rt.db.GetRawDB().QueryRow("SELECT username FROM users WHERE id = ?", userId).Scan(&userUsername)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// Get contact details
	var contactUsername string
	err = rt.db.GetRawDB().QueryRow("SELECT username FROM users WHERE id = ?", contactId).Scan(&contactUsername)
	if err != nil {
		http.Error(w, "contact not found", http.StatusNotFound)
		return
	}

	// Remove contact
	user := database.User{UId: userId, Username: userUsername}
	contact := database.User{UId: contactId, Username: contactUsername}

	_, err = rt.db.RemoveContact(user, contact)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
