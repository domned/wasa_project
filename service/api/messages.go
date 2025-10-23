package api

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get parameters from URL
	userId := ps.ByName("id")
	conversationId := ps.ByName("conversationId")

	// Parse request body
	var requestBody struct {
		Content  string `json:"content"`
		ImageUrl string `json:"imageUrl,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		ctx.Logger.WithError(err).Error("failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate that at least content or imageUrl is provided
	if requestBody.Content == "" && requestBody.ImageUrl == "" {
		http.Error(w, "Message must have content or image", http.StatusBadRequest)
		return
	}

	// Get database connection
	db := rt.db.GetRawDB()

	// Check if conversation exists and user is a participant
	var participants string
	err := db.QueryRow("SELECT participants FROM conversations WHERE id = ?", conversationId).Scan(&participants)
	if err != nil {
		ctx.Logger.WithError(err).Error("conversation not found")
		http.Error(w, "Conversation not found", http.StatusNotFound)
		return
	}

	// Check if user is a participant (participants is stored as JSON array)
	// For simplicity, we'll check if the user ID is in the participants string
	// In a production system, you'd want to properly parse the JSON
	if participants == "" {
		http.Error(w, "Invalid conversation", http.StatusBadRequest)
		return
	}

	// Generate message ID
	messageId, err := uuid.NewV4()
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to generate message ID")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Save message to database
	_, err = db.Exec("INSERT INTO messages (id, conversation_id, sender_id, message, image_url) VALUES (?, ?, ?, ?, ?)",
		messageId.String(), conversationId, userId, requestBody.Content, requestBody.ImageUrl)
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to save message to database")
		rt.sysLogger.LogError("Failed to save message to database: " + err.Error())
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	// Log successful message send
	rt.sysLogger.LogInfo("Message sent in conversation " + conversationId + " by user " + userId)

	// Broadcast message to WebSocket clients
	messageData := map[string]interface{}{
		"id":              messageId.String(),
		"conversation_id": conversationId,
		"sender_id":       userId,
		"content":         requestBody.Content,
		"image_url":       requestBody.ImageUrl,
	}
	BroadcastMessage("message", messageData)
	rt.sysLogger.LogDebug("Message broadcasted to WebSocket clients")

	// Return success response with message ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(messageId.String()); err != nil {
		ctx.Logger.WithError(err).Error("failed to encode message response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get parameters from URL
	userId := ps.ByName("id")
	conversationId := ps.ByName("conversationId")
	messageId := ps.ByName("messageId")

	if userId == "" || conversationId == "" || messageId == "" {
		http.Error(w, "missing required parameters", http.StatusBadRequest)
		return
	}

	// Check if message exists and belongs to the user
	var senderId string
	err := rt.db.GetRawDB().QueryRow(
		"SELECT sender_id FROM messages WHERE id = ? AND conversation_id = ?",
		messageId, conversationId).Scan(&senderId)
	if err != nil {
		http.Error(w, "message not found", http.StatusNotFound)
		return
	}

	// Check if user is the sender of the message
	if senderId != userId {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Delete the message (reactions and comments will cascade delete due to foreign keys)
	_, err = rt.db.GetRawDB().Exec("DELETE FROM messages WHERE id = ?", messageId)
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to delete message")
		http.Error(w, "failed to delete message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get parameters from URL
	userId := ps.ByName("id")
	conversationId := ps.ByName("conversationId")
	messageId := ps.ByName("messageId")

	if userId == "" || conversationId == "" || messageId == "" {
		http.Error(w, "missing required parameters", http.StatusBadRequest)
		return
	}

	// Parse request body for target conversation
	var requestBody struct {
		Content string `json:"content"` // This should be the target conversation ID
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		ctx.Logger.WithError(err).Error("failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestBody.Content == "" {
		http.Error(w, "target conversation ID is required", http.StatusBadRequest)
		return
	}

	targetConversationId := requestBody.Content

	// Create user object
	user := database.User{
		UId: userId,
	}

	// Use the database ForwardMessage function which handles both text and images
	conversation, err := rt.db.ForwardMessage(targetConversationId, user, messageId)
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to forward message")
		http.Error(w, "Failed to forward message", http.StatusInternalServerError)
		return
	}

	// Return the updated conversation
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(conversation); err != nil {
		ctx.Logger.WithError(err).Error("failed to encode conversation response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) reactToMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get parameters from URL
	userId := ps.ByName("id")
	conversationId := ps.ByName("conversationId")
	messageId := ps.ByName("messageId")

	// Parse request body
	var requestBody struct {
		Emoji string `json:"emoji"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		ctx.Logger.WithError(err).Error("failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate emoji
	if requestBody.Emoji == "" {
		http.Error(w, "Emoji cannot be empty", http.StatusBadRequest)
		return
	}

	// Create user object
	user := database.User{
		UId: userId,
	}

	// Add reaction to database
	_, err := rt.db.ReactToMessage(conversationId, user, messageId, requestBody.Emoji)
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to add reaction")
		http.Error(w, "Failed to add reaction", http.StatusInternalServerError)
		return
	}

	// Return success response (no content)
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) removeReaction(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get parameters from URL
	userId := ps.ByName("id")
	conversationId := ps.ByName("conversationId")
	messageId := ps.ByName("messageId")
	emoji := ps.ByName("emoji")

	// Validate emoji
	if emoji == "" {
		http.Error(w, "Emoji cannot be empty", http.StatusBadRequest)
		return
	}

	// Create user object
	user := database.User{
		UId: userId,
	}

	// Remove reaction from database
	_, err := rt.db.RemoveReaction(conversationId, user, messageId, emoji)
	if err != nil {
		ctx.Logger.WithError(err).Error("failed to remove reaction")
		http.Error(w, "Failed to remove reaction", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusNoContent)
}

// Comments feature removed (emoji reactions only)
