package database

import (
	"database/sql"
	"time"
	"github.com/gofrs/uuid"
)

func (db *appdbimpl) SendMessage(cid string, user User, message string) (Conversation, error) {
	return db.SendMessageWithImage(cid, user, message, "")
}

func (db *appdbimpl) SendMessageWithImage(cid string, user User, message string, imageUrl string) (Conversation, error) {
	id, err := uuid.NewV4()
	if err != nil {
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
	var message, imageUrl string
	err := db.c.QueryRow("SELECT message, COALESCE(image_url, '') FROM messages WHERE id = ?", mid).Scan(&message, &imageUrl)
	if err != nil {
		return Conversation{}, err
	}

	return db.SendMessageWithImage(cid, user, message, imageUrl)
}

func (db *appdbimpl) ReactToMessage(cid string, user User, mid string, emoji string) (Conversation, error) {
	// Check if reaction already exists
	var existingId string
	err := db.c.QueryRow("SELECT id FROM reactions WHERE message_id = ? AND sender_id = ? AND emoji = ?",
		mid, user.UId, emoji).Scan(&existingId)
	
	if err == nil {
		// Reaction exists, remove it (toggle off)
		return db.RemoveReaction(cid, user, mid, emoji)
	} else if err != sql.ErrNoRows {
		// Some other error occurred
		return Conversation{}, err
	}

	// Reaction doesn't exist (sql.ErrNoRows), add it
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

func (db *appdbimpl) RemoveReaction(cid string, user User, mid string, emoji string) (Conversation, error) {
	_, err := db.c.Exec("DELETE FROM reactions WHERE message_id = ? AND sender_id = ? AND emoji = ?",
		mid, user.UId, emoji)
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
	_, err := db.c.Exec("DELETE FROM comments WHERE id = ? AND sender_id = ?", commentId, user.UId)
	if err != nil {
		return Conversation{}, err
	}

	return db.GetConversation(cid)
}

func (db *appdbimpl) MarkMessageAsRead(messageId string, userId string) error {
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	// Use INSERT OR REPLACE to handle the case where the message is already marked as read
	_, err = db.c.Exec(`INSERT OR REPLACE INTO read_status (id, message_id, user_id, read_at) 
		VALUES (?, ?, ?, ?)`, id.String(), messageId, userId, time.Now().Unix())
	return err
}

func (db *appdbimpl) GetUnreadCount(conversationId string, userId string) (int, error) {
	var count int
	err := db.c.QueryRow(`
		SELECT COUNT(*) 
		FROM messages m 
		LEFT JOIN read_status r ON m.id = r.message_id AND r.user_id = ?
		WHERE m.conversation_id = ? AND m.sender_id != ? AND r.message_id IS NULL`,
		userId, conversationId, userId).Scan(&count)
	return count, err
}