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
	
	// Get messages
	rows, err := sqldb.Query(`
		SELECT m.id, m.sender_id, m.message, COALESCE(m.image_url, '') as image_url, u.username 
		FROM messages m 
		JOIN users u ON m.sender_id = u.id 
		WHERE m.conversation_id = ? 
		ORDER BY m.rowid ASC`, convId)
       if err != nil {
	       http.Error(w, err.Error(), http.StatusInternalServerError)
	       return
       }
       defer rows.Close()
       
       var messages []map[string]interface{}
       var messageIds []string
       for rows.Next() {
	       var id, senderId, message, imageUrl, senderUsername string
	       if err := rows.Scan(&id, &senderId, &message, &imageUrl, &senderUsername); err != nil {
		       http.Error(w, err.Error(), http.StatusInternalServerError)
		       return
	       }
	       
	       // Collect message IDs for marking as read
	       if senderId != userId {
		       messageIds = append(messageIds, id)
	       }
	       
	       // Get reactions for this message
	       reactionRows, err := sqldb.Query(`
	       	SELECT emoji, COUNT(*) as count, GROUP_CONCAT(u.username, ',') as usernames
	       	FROM reactions r
	       	JOIN users u ON r.sender_id = u.id
	       	WHERE r.message_id = ?
	       	GROUP BY emoji`, id)
	       if err != nil {
		       http.Error(w, err.Error(), http.StatusInternalServerError)
		       return
	       }
	       
	       reactions := make(map[string]interface{})
	       for reactionRows.Next() {
		       var emoji string
		       var count int
		       var usernames string
		       if err := reactionRows.Scan(&emoji, &count, &usernames); err != nil {
			       reactionRows.Close()
			       http.Error(w, err.Error(), http.StatusInternalServerError)
			       return
		       }
		       reactions[emoji] = map[string]interface{}{
			       "count": count,
			       "users": usernames,
		       }
	       }
	       reactionRows.Close()
	       
	       // Check read status for this message
	       var isMessageRead bool
	       
	       if senderId == userId {
	       	// For messages sent by current user: check if any other participant has read it
	       	var readCount int
	       	err = sqldb.QueryRow(`
	       		SELECT COUNT(*) 
	       		FROM read_status 
	       		WHERE message_id = ? AND user_id != ?`, 
	       		id, userId).Scan(&readCount)
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
		       "id": id,
		       "senderId": senderId,
		       "text": message,
		       "imageUrl": imageUrl,
		       "senderUsername": senderUsername,
		       "reactions": reactions,
		       "isRead": isMessageRead,
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
       json.NewEncoder(w).Encode(messages)
}
