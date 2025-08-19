package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) createConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	   // Get user ID from URL
	   userId := ps.ByName("id")
	   if userId == "" {
		   http.Error(w, "missing user id", http.StatusBadRequest)
		   return
	   }

		// Build a User struct (other fields can be empty)
		user := database.User{UId: userId}

	   // Fetch conversations from DB
	   conversations, err := rt.db.GetMyConversations(user)
	   if err != nil {
		   http.Error(w, err.Error(), http.StatusInternalServerError)
		   return
	   }

	   w.Header().Set("Content-Type", "application/json")
	   w.WriteHeader(http.StatusOK)
	   if err := json.NewEncoder(w).Encode(conversations); err != nil {
		   http.Error(w, err.Error(), http.StatusInternalServerError)
	   }
}

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get user ID from URL
	userId := ps.ByName("id")
	if userId == "" {
		http.Error(w, "missing user id", http.StatusBadRequest)
		return
	}

	// Build a User struct (other fields can be empty)
	user := database.User{UId: userId}

	// Fetch conversations from DB
	conversations, err := rt.db.GetMyConversations(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(conversations); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (rt *_router) addtoGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.WriteHeader(http.StatusNotImplemented)
} 