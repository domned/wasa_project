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

	// Parse request body for participants
	var request struct {
		Participants []string `json:"participants"`
		Name         string   `json:"name,omitempty"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate participants
	if len(request.Participants) == 0 {
		http.Error(w, "participants cannot be empty", http.StatusBadRequest)
		return
	}

	// Add current user to participants if not already included
	found := false
	for _, pid := range request.Participants {
		if pid == userId {
			found = true
			break
		}
	}
	if !found {
		request.Participants = append(request.Participants, userId)
	}

	// Convert participant IDs to User structs and validate they exist
	var participants []database.User
	for _, pid := range request.Participants {
		// Check if user exists in database
		var username, picture string
		err := rt.db.GetRawDB().QueryRow("SELECT username, COALESCE(picture, '') FROM users WHERE id = ?", pid).
			Scan(&username, &picture)
		if err != nil {
			http.Error(w, "participant not found: "+pid, http.StatusBadRequest)
			return
		}
		participants = append(participants, database.User{
			UId:      pid,
			Username: username,
			Picture:  picture,
		})
	}

	// Create conversation
	conversation, err := rt.db.CreateConversation(participants, request.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(conversation); err != nil {
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
	// Get user ID and conversation ID from URL
	userId := ps.ByName("id")
	conversationId := ps.ByName("conversationId")
	
	if userId == "" || conversationId == "" {
		http.Error(w, "missing user id or conversation id", http.StatusBadRequest)
		return
	}

	// Get conversation details
	conversation, err := rt.db.GetConversation(conversationId)
	if err != nil {
		http.Error(w, "conversation not found", http.StatusNotFound)
		return
	}

	// Check if user is participant in conversation
	isParticipant := false
	for _, participant := range conversation.Participants {
		if participant.UId == userId {
			isParticipant = true
			break
		}
	}
	
	if !isParticipant {
		http.Error(w, "unauthorized", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(conversation); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (rt *_router) addtoGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get user ID and conversation ID from URL
	userId := ps.ByName("id")
	conversationId := ps.ByName("conversationId")
	
	if userId == "" || conversationId == "" {
		http.Error(w, "missing user id or conversation id", http.StatusBadRequest)
		return
	}

	// Parse request body for new member
	var request struct {
		Name string `json:"name"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if request.Name == "" {
		http.Error(w, "member name cannot be empty", http.StatusBadRequest)
		return
	}

	// Find user by username
	var memberUserId, memberUsername string
	err := rt.db.GetRawDB().QueryRow("SELECT id, username FROM users WHERE username = ?", request.Name).
		Scan(&memberUserId, &memberUsername)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// Add user to group
	memberUser := database.User{UId: memberUserId, Username: memberUsername}
	_, err = rt.db.AddToGroup(conversationId, memberUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(memberUserId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get user ID and conversation ID from URL
	userId := ps.ByName("id")
	conversationId := ps.ByName("conversationId")
	
	if userId == "" || conversationId == "" {
		http.Error(w, "missing user id or conversation id", http.StatusBadRequest)
		return
	}

	// Get user details
	var username string
	err := rt.db.GetRawDB().QueryRow("SELECT username FROM users WHERE id = ?", userId).Scan(&username)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	// Remove user from group
	user := database.User{UId: userId, Username: username}
	_, err = rt.db.LeaveGroup(conversationId, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get user ID and conversation ID from URL
	userId := ps.ByName("id")
	conversationId := ps.ByName("conversationId")
	
	if userId == "" || conversationId == "" {
		http.Error(w, "missing user id or conversation id", http.StatusBadRequest)
		return
	}

	// Parse request body for new name
	var newName string
	if err := json.NewDecoder(r.Body).Decode(&newName); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if newName == "" {
		http.Error(w, "name cannot be empty", http.StatusBadRequest)
		return
	}

	// Update conversation name
	_, err := rt.db.GetRawDB().Exec("UPDATE conversations SET name = ? WHERE id = ?", newName, conversationId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get user ID and conversation ID from URL
	userId := ps.ByName("id")
	conversationId := ps.ByName("conversationId")
	
	if userId == "" || conversationId == "" {
		http.Error(w, "missing user id or conversation id", http.StatusBadRequest)
		return
	}

	// Parse request body for new picture URL
	var newPicture string
	if err := json.NewDecoder(r.Body).Decode(&newPicture); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if newPicture == "" {
		http.Error(w, "picture URL cannot be empty", http.StatusBadRequest)
		return
	}

	// Update conversation picture
	_, err := rt.db.SetGroupPhoto(conversationId, newPicture)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
} 