package api

import (
	"encoding/json"
	"net/http"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
       convId := ps.ByName("convId")
	db := rt.db
	sqldb := db.GetRawDB()
	rows, err := sqldb.Query(`
		SELECT m.id, m.sender_id, m.message, u.username 
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
       for rows.Next() {
	       var id, senderId, message, senderUsername string
	       if err := rows.Scan(&id, &senderId, &message, &senderUsername); err != nil {
		       http.Error(w, err.Error(), http.StatusInternalServerError)
		       return
	       }
	       messages = append(messages, map[string]interface{}{
		       "id": id,
		       "senderId": senderId,
		       "text": message,
		       "senderUsername": senderUsername,
	       })
       }
       w.Header().Set("Content-Type", "application/json")
       json.NewEncoder(w).Encode(messages)
}
