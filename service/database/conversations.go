package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/gofrs/uuid"
)

var demoNames = []string{"Alice & Bob", "Charlie Group", "Delta Squad", "Echo Team", "Foxtrot Chat"}
var demoAvatars = []string{
	"https://randomuser.me/api/portraits/men/1.jpg",
	"https://randomuser.me/api/portraits/women/2.jpg",
	"https://randomuser.me/api/portraits/men/3.jpg",
	"https://randomuser.me/api/portraits/women/4.jpg",
	"https://randomuser.me/api/portraits/men/5.jpg",
}

func (db *appdbimpl) CreateConversation(participants []User, name string) (Conversation, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return Conversation{}, err
	}

	// Use provided name or generate a random one if empty
	var conversationName string
	if name != "" {
		conversationName = name
	} else {
		rand.Seed(time.Now().UnixNano())
		conversationName = demoNames[rand.Intn(len(demoNames))]
	}

	rand.Seed(time.Now().UnixNano())
	avatar := demoAvatars[rand.Intn(len(demoAvatars))]

	// Extract only user IDs for storage
	var participantIDs []string
	for _, participant := range participants {
		participantIDs = append(participantIDs, participant.UId)
	}

	participantsJSON, err := json.Marshal(participantIDs)
	if err != nil {
		return Conversation{}, err
	}

	_, err = db.c.Exec("INSERT INTO conversations (id, participants, name, picture) VALUES (?, ?, ?, ?)",
		id.String(), string(participantsJSON), conversationName, avatar)
	if err != nil {
		return Conversation{}, err
	}

	return Conversation{CId: id.String(), Participants: participants, Name: conversationName, Picture: avatar}, nil
}

func (db *appdbimpl) GetMyConversations(user User) ([]Conversation, error) {
	// Modified query to include last message information and sort by last message time
	rows, err := db.c.Query(`
		SELECT 
			c.id, 
			c.participants, 
			c.name, 
			c.picture,
			m.id as last_msg_id,
			m.sender_id as last_msg_sender_id,
			m.message as last_msg_text,
			COALESCE(m.image_url, '') as last_msg_image_url,
			u.username as last_msg_sender_username,
			CAST((julianday(m.timestamp) - 2440587.5) * 86400000 AS INTEGER) as last_msg_time
		FROM conversations c
		LEFT JOIN (
			SELECT conversation_id, MAX(timestamp) as max_timestamp
			FROM messages 
			GROUP BY conversation_id
		) latest ON c.id = latest.conversation_id
		LEFT JOIN messages m ON latest.conversation_id = m.conversation_id AND latest.max_timestamp = m.timestamp
		LEFT JOIN users u ON m.sender_id = u.id
		ORDER BY m.timestamp DESC NULLS LAST`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		var conv Conversation
		var participantsJSON string
		var name sql.NullString
		var picture sql.NullString
		var lastMsgId sql.NullString
		var lastMsgSenderId sql.NullString
		var lastMsgText sql.NullString
		var lastMsgImageUrl sql.NullString
		var lastMsgSenderUsername sql.NullString
		var lastMsgTime sql.NullInt64

		if scanErr := rows.Scan(&conv.CId, &participantsJSON, &name, &picture,
			&lastMsgId, &lastMsgSenderId, &lastMsgText, &lastMsgImageUrl, &lastMsgSenderUsername, &lastMsgTime); scanErr != nil {
			return nil, scanErr
		}

		if name.Valid {
			conv.Name = name.String
		} else {
			conv.Name = ""
		}
		if picture.Valid {
			conv.Picture = picture.String
		} else {
			conv.Picture = ""
		}

		var participantIDs []string

		// Try to unmarshal as string array first (new format)
		err = json.Unmarshal([]byte(participantsJSON), &participantIDs)
		if err != nil {
			// If that fails, try to unmarshal as User object array (old format)
			var participantObjects []User
			err = json.Unmarshal([]byte(participantsJSON), &participantObjects)
			if err != nil {
				return nil, err
			}
			// Extract IDs from User objects
			for _, participant := range participantObjects {
				participantIDs = append(participantIDs, participant.UId)
			}
		}

		// Only include conversations where the user is a participant
		found := false
		for _, pid := range participantIDs {
			if pid == user.UId {
				found = true
				break
			}
		}

		if found {
			var participants []User
			for _, uid := range participantIDs {
				var u User
				var upic sql.NullString
				if qErr := db.c.QueryRow("SELECT id, username, picture FROM users WHERE id = ?", uid).Scan(&u.UId, &u.Username, &upic); qErr != nil {
					return nil, qErr
				}
				if upic.Valid {
					u.Picture = upic.String
				}
				participants = append(participants, u)
			}
			conv.Participants = participants

			// Add last message information if available
			if lastMsgId.Valid {
				conv.LastMessage = &Message{
					Id:             lastMsgId.String,
					SenderId:       lastMsgSenderId.String,
					Text:           lastMsgText.String,
					ImageUrl:       lastMsgImageUrl.String,
					SenderUsername: lastMsgSenderUsername.String,
				}
				if lastMsgTime.Valid {
					conv.LastMessageTime = fmt.Sprintf("%d", lastMsgTime.Int64)
				}
			}

			// Get unread count for this conversation
			unreadCount, ucErr := db.GetUnreadCount(conv.CId, user.UId)
			if ucErr == nil {
				conv.UnreadCount = unreadCount
			}

			conversations = append(conversations, conv)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return conversations, nil
}

func (db *appdbimpl) GetConversation(cid string) (Conversation, error) {
	var conv Conversation
	var participantsJSON string
	err := db.c.QueryRow("SELECT id, participants FROM conversations WHERE id = ?", cid).
		Scan(&conv.CId, &participantsJSON)
	if err != nil {
		return Conversation{}, err
	}

	var participantIDs []string

	// Try to unmarshal as string array first (new format)
	err = json.Unmarshal([]byte(participantsJSON), &participantIDs)
	if err != nil {
		// If that fails, try to unmarshal as User object array (old format)
		var participantObjects []User
		err = json.Unmarshal([]byte(participantsJSON), &participantObjects)
		if err != nil {
			return Conversation{}, err
		}
		// Extract IDs from User objects
		for _, participant := range participantObjects {
			participantIDs = append(participantIDs, participant.UId)
		}
	}
	var participants []User
	for _, uid := range participantIDs {
		var user User
		var picture sql.NullString
		err := db.c.QueryRow("SELECT id, username, picture FROM users WHERE id = ?", uid).Scan(&user.UId, &user.Username, &picture)
		if err != nil {
			return Conversation{}, err
		}
		if picture.Valid {
			user.Picture = picture.String
		} else {
			user.Picture = ""
		}
		participants = append(participants, user)
	}
	conv.Participants = participants
	return conv, nil
}

func (db *appdbimpl) AddToGroup(cid string, user User) (Conversation, error) {
	// Fetch current participant IDs
	var participantsJSON string
	err := db.c.QueryRow("SELECT participants FROM conversations WHERE id = ?", cid).Scan(&participantsJSON)
	if err != nil {
		return Conversation{}, err
	}
	var participantIDs []string

	// Try to unmarshal as string array first (new format)
	err = json.Unmarshal([]byte(participantsJSON), &participantIDs)
	if err != nil {
		// If that fails, try to unmarshal as User object array (old format)
		var participantObjects []User
		err = json.Unmarshal([]byte(participantsJSON), &participantObjects)
		if err != nil {
			return Conversation{}, err
		}
		// Extract IDs from User objects
		for _, participant := range participantObjects {
			participantIDs = append(participantIDs, participant.UId)
		}
	}
	// Add new user if not already present
	for _, pid := range participantIDs {
		if pid == user.UId {
			return db.GetConversation(cid) // already a member
		}
	}
	participantIDs = append(participantIDs, user.UId)
	newJSON, _ := json.Marshal(participantIDs)
	_, err = db.c.Exec("UPDATE conversations SET participants = ? WHERE id = ?", string(newJSON), cid)
	if err != nil {
		return Conversation{}, err
	}
	return db.GetConversation(cid)
}

func (db *appdbimpl) LeaveGroup(cid string, user User) (Conversation, error) {
	// Fetch current participant IDs
	var participantsJSON string
	err := db.c.QueryRow("SELECT participants FROM conversations WHERE id = ?", cid).Scan(&participantsJSON)
	if err != nil {
		return Conversation{}, err
	}
	var participantIDs []string

	// Try to unmarshal as string array first (new format)
	err = json.Unmarshal([]byte(participantsJSON), &participantIDs)
	if err != nil {
		// If that fails, try to unmarshal as User object array (old format)
		var participantObjects []User
		err = json.Unmarshal([]byte(participantsJSON), &participantObjects)
		if err != nil {
			return Conversation{}, err
		}
		// Extract IDs from User objects
		for _, participant := range participantObjects {
			participantIDs = append(participantIDs, participant.UId)
		}
	}
	// Remove user
	var newIDs []string
	for _, pid := range participantIDs {
		if pid != user.UId {
			newIDs = append(newIDs, pid)
		}
	}
	newJSON, _ := json.Marshal(newIDs)
	_, err = db.c.Exec("UPDATE conversations SET participants = ? WHERE id = ?", string(newJSON), cid)
	if err != nil {
		return Conversation{}, err
	}
	return db.GetConversation(cid)
}

func (db *appdbimpl) SetGroupName(cid string, name string) (Conversation, error) {
	// Note: Group name is not in the current schema, would need to add a name column
	return db.GetConversation(cid)
}

func (db *appdbimpl) SetGroupPhoto(cid string, picture string) (Conversation, error) {
	_, err := db.c.Exec("UPDATE conversations SET picture = ? WHERE id = ?", picture, cid)
	if err != nil {
		return Conversation{}, err
	}
	return db.GetConversation(cid)
}

func (db *appdbimpl) GetUnreadCount(conversationId string, userId string) (int, error) {
	var count int
	err := db.c.QueryRow(`
		SELECT COUNT(*) 
		FROM messages 
		WHERE conversation_id = ? 
		AND sender_id != ? 
		AND timestamp > COALESCE(
			(SELECT last_read_timestamp FROM conversation_participants 
			 WHERE conversation_id = ? AND user_id = ?), 
			'1970-01-01 00:00:00'
		)`, conversationId, userId, conversationId, userId).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}
