package database

import (
	"github.com/gofrs/uuid"
)

func (db *appdbimpl) sendMessage(cid string, user User, message string) (Conversation, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return Conversation{}, err
	}

	_, err = db.c.Exec("INSERT INTO messages (id, conversation_id, sender_id, message) VALUES (?, ?, ?, ?)",
		id.String(), cid, user.UId, message)
	if err != nil {
		return Conversation{}, err
	}

	return db.getConversation(cid)
}

func (db *appdbimpl) deleteMessage(cid string, user User, mid string) (Conversation, error) {
	_, err := db.c.Exec("DELETE FROM messages WHERE id = ? AND sender_id = ?", mid, user.UId)
	if err != nil {
		return Conversation{}, err
	}
	return db.getConversation(cid)
}

func (db *appdbimpl) forwardMessage(cid string, user User, mid string) (Conversation, error) {
	var message string
	err := db.c.QueryRow("SELECT message FROM messages WHERE id = ?", mid).Scan(&message)
	if err != nil {
		return Conversation{}, err
	}

	return db.sendMessage(cid, user, message)
}

func (db *appdbimpl) reactToMessage(cid string, user User, mid string, emoji string) (Conversation, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return Conversation{}, err
	}

	_, err = db.c.Exec("INSERT INTO reactions (id, message_id, sender_id, emoji) VALUES (?, ?, ?, ?)",
		id.String(), mid, user.UId, emoji)
	if err != nil {
		return Conversation{}, err
	}

	return db.getConversation(cid)
}

func (db *appdbimpl) removeReaction(cid string, user User, mid string, emoji string) (Conversation, error) {
	_, err := db.c.Exec("DELETE FROM reactions WHERE message_id = ? AND sender_id = ? AND emoji = ?",
		mid, user.UId, emoji)
	if err != nil {
		return Conversation{}, err
	}

	return db.getConversation(cid)
}

func (db *appdbimpl) commentMessage(cid string, user User, mid string, comment string) (Conversation, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return Conversation{}, err
	}

	_, err = db.c.Exec("INSERT INTO comments (id, message_id, sender_id, comment) VALUES (?, ?, ?, ?)",
		id.String(), mid, user.UId, comment)
	if err != nil {
		return Conversation{}, err
	}

	return db.getConversation(cid)
}

func (db *appdbimpl) uncommentMessage(cid string, user User, mid string, commentId string) (Conversation, error) {
	_, err := db.c.Exec("DELETE FROM comments WHERE id = ? AND sender_id = ?", commentId, user.UId)
	if err != nil {
		return Conversation{}, err
	}

	return db.getConversation(cid)
} 