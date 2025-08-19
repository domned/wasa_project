package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var name string
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&name); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	user, err := rt.db.SetMyUserName(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var photo string
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&photo); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	user, err := rt.db.SetMyPhoto(photo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}