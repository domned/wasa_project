package database

import (
	"time"

	"github.com/gofrs/uuid"
)

func (db *appdbimpl) SendMessage(cid string, user User, content string) (Conversation, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return Conversation{}, err
	}

	conversation, err := db.GetConversation(cid)
	if err != nil {
		return Conversation{}, err
	}

	// Check if user is participant
	found := false
	for _, participant := range conversation.Participants {
		if participant.UId == user.UId {
			found = true
			break
		}
	}
	if !found {
		return Conversation{}, err
	}

	_, err = db.c.Exec("INSERT INTO messages (id, conversation_id, sender_id, message) VALUES (?, ?, ?, ?)",
		id.String(), cid, user.UId, content)
	if err != nil {
		return Conversation{}, err
	}

	return db.GetConversation(cid)
}

func (db *appdbimpl) SendMessageWithImage(cid string, user User, message string, imageUrl string) (Conversation, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return Conversation{}, err
	}

	conversation, err := db.GetConversation(cid)
	if err != nil {
		return Conversation{}, err
	}

	// Check if user is participant
	found := false
	for _, participant := range conversation.Participants {
		if participant.UId == user.UId {
			found = true
			break
		}
	}
	if !found {
		return Conversation{}, err
	}

	_, err = db.c.Exec("INSERT INTO messages (id, conversation_id, sender_id, message, image_url) VALUES (?, ?, ?, ?, ?)",
		id.String(), cid, user.UId, message, imageUrl)
	if err != nil {
		return Conversation{}, err
	}

	return db.GetConversation(cid)
}

func (db *appdbimpl) DeleteMessage(cid string, user User, mid string) (Conversation, error) {
	_, err := db.c.Exec("DELETE FROM messages WHERE id = ? AND sender_id = ?", mid, user.UId)
	if err != nil {
		return Conversation{}, err
	}
	return db.GetConversation(cid)
}

func (db *appdbimpl) ForwardMessage(cid string, user User, mid string) (Conversation, error) {
	var message string
	err := db.c.QueryRow("SELECT message FROM messages WHERE id = ?", mid).Scan(&message)
	if err != nil {
		return Conversation{}, err
	}

	_, err = db.SendMessage(cid, user, message)
	if err != nil {
		return Conversation{}, err
	}

	return db.GetConversation(cid)
}

func (db *appdbimpl) GetConversationMessages(cid string) ([]Message, error) {
	rows, err := db.c.Query(`
		SELECT m.id, m.conversation_id, m.message, m.sender_id, m.timestamp, u.username
		FROM messages m 
		JOIN users u ON m.sender_id = u.id 
		WHERE m.conversation_id = ? 
		ORDER BY m.timestamp ASC`, cid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		var sender User
		var timestamp time.Time
		var conversationId string
		if scanErr := rows.Scan(&m.Id, &conversationId, &m.Text, &sender.UId, &timestamp, &sender.Username); scanErr != nil {
			return nil, scanErr
		}
		m.SenderId = sender.UId
		m.SenderUsername = sender.Username
		m.Time = timestamp.Format(time.RFC3339)
		messages = append(messages, m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (db *appdbimpl) ReactToMessage(cid string, user User, mid string, emoji string) (Conversation, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return Conversation{}, err
	}

	_, err = db.c.Exec("INSERT INTO reactions (id, message_id, sender_id, emoji) VALUES (?, ?, ?, ?)",
		id.String(), mid, user.UId, emoji)
	if err != nil {
		return Conversation{}, err
	}

	return db.GetConversation(cid)
}

func (db *appdbimpl) RemoveReaction(cid string, user User, mid string, reactionId string) (Conversation, error) {
	_, err := db.c.Exec("DELETE FROM reactions WHERE id = ? AND sender_id = ?",
		reactionId, user.UId)
	if err != nil {
		return Conversation{}, err
	}

	return db.GetConversation(cid)
}

func (db *appdbimpl) CommentMessage(cid string, user User, mid string, comment string) (Conversation, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return Conversation{}, err
	}

	_, err = db.c.Exec("INSERT INTO comments (id, message_id, sender_id, comment) VALUES (?, ?, ?, ?)",
		id.String(), mid, user.UId, comment)
	if err != nil {
		return Conversation{}, err
	}

	return db.GetConversation(cid)
}

func (db *appdbimpl) UncommentMessage(cid string, user User, mid string, commentId string) (Conversation, error) {
	_, err := db.c.Exec("DELETE FROM comments WHERE id = ? AND sender_id = ?",
		commentId, user.UId)
	if err != nil {
		return Conversation{}, err
	}

	return db.GetConversation(cid)
}

func (db *appdbimpl) MarkMessageAsRead(messageId string, userId string) error {
	_, err := db.c.Exec(`
		UPDATE conversation_participants 
		SET last_read_timestamp = (
			SELECT timestamp FROM messages WHERE id = ?
		) 
		WHERE user_id = ? AND conversation_id = (
			SELECT conversation_id FROM messages WHERE id = ?
		)`, messageId, userId, messageId)

	return err
}
