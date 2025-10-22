package api

import (
	"encoding/json"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	convId := ps.ByName("conversationId")
	userId := ps.ByName("id")
	db := rt.db
	sqldb := db.GetRawDB()

	// Get messages using database interface
	dbMessages, err := db.GetConversationMessages(convId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var messages []map[string]interface{}
	var messageIds []string
	for _, msg := range dbMessages {
		// Collect message IDs for marking as read
		if msg.SenderId != userId {
			messageIds = append(messageIds, msg.Id)
		}

		// Get reactions for this message
		reactionRows, err := sqldb.Query(`
	       	SELECT emoji, COUNT(*) as count, GROUP_CONCAT(u.username, ',') as usernames
	       	FROM reactions r
	       	JOIN users u ON r.sender_id = u.id
	       	WHERE r.message_id = ?
	       	GROUP BY emoji`, msg.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		reactions := make(map[string]interface{})
		for reactionRows.Next() {
			var emoji string
			var count int
			var usernames string
			if scanErr := reactionRows.Scan(&emoji, &count, &usernames); scanErr != nil {
				reactionRows.Close()
				http.Error(w, scanErr.Error(), http.StatusInternalServerError)
				return
			}
			reactions[emoji] = map[string]interface{}{
				"count": count,
				"users": usernames,
			}
		}
		if err = reactionRows.Err(); err != nil {
			reactionRows.Close()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		reactionRows.Close()

		// Check read status for this message
		var isMessageRead bool

		if msg.SenderId == userId {
			// For messages sent by current user: check if any other participant has read it
			var readCount int
			err = sqldb.QueryRow(`
	       		SELECT COUNT(*) 
	       		FROM read_status 
	       		WHERE message_id = ? AND user_id != ?`,
				msg.Id, userId).Scan(&readCount)
			if err != nil {
				readCount = 0
			}
			isMessageRead = readCount > 0
		} else {
			// For messages received by current user: always show as "delivered"
			// (we could check if current user has read it, but for sender's view we care about recipients)
			isMessageRead = false
		}

		messages = append(messages, map[string]interface{}{
			"id":             msg.Id,
			"senderId":       msg.SenderId,
			"text":           msg.Text,
			"imageUrl":       msg.ImageUrl,
			"senderUsername": msg.SenderUsername,
			"time":           msg.Time,
			"reactions":      reactions,
			"isRead":         isMessageRead,
		})
	}

	// Mark all unread messages as read
	for _, msgId := range messageIds {
		err := db.MarkMessageAsRead(msgId, userId)
		if err != nil {
			ctx.Logger.WithError(err).Error("failed to mark message as read")
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		ctx.Logger.WithError(err).Error("failed to encode messages response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
