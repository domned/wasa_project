package database

import (
	"encoding/json"

	"github.com/gofrs/uuid"
)

func (db *appdbimpl) createConversation(participants []User) (Conversation, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return Conversation{}, err
	}

	participantsJSON, err := json.Marshal(participants)
	if err != nil {
		return Conversation{}, err
	}

	_, err = db.c.Exec("INSERT INTO conversations (id, participants) VALUES (?, ?)", 
		id.String(), string(participantsJSON))
	if err != nil {
		return Conversation{}, err
	}

	return Conversation{CId: id.String(), Participants: participants}, nil
}

func (db *appdbimpl) getMyConversations(user User) ([]Conversation, error) {
	rows, err := db.c.Query("SELECT id, participants FROM conversations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		var conv Conversation
		var participantsJSON string
		err := rows.Scan(&conv.CId, &participantsJSON)
		if err != nil {
			return nil, err
		}

		var participants []User
		err = json.Unmarshal([]byte(participantsJSON), &participants)
		if err != nil {
			return nil, err
		}

		conv.Participants = participants
		conversations = append(conversations, conv)
	}
	return conversations, nil
}

func (db *appdbimpl) getConversation(cid string) (Conversation, error) {
	var conv Conversation
	var participantsJSON string
	err := db.c.QueryRow("SELECT id, participants FROM conversations WHERE id = ?", cid).
		Scan(&conv.CId, &participantsJSON)
	if err != nil {
		return Conversation{}, err
	}

	var participants []User
	err = json.Unmarshal([]byte(participantsJSON), &participants)
	if err != nil {
		return Conversation{}, err
	}

	conv.Participants = participants
	return conv, nil
}

func (db *appdbimpl) addtoGroup(cid string, user User) (Conversation, error) {
	conv, err := db.getConversation(cid)
	if err != nil {
		return Conversation{}, err
	}

	conv.Participants = append(conv.Participants, user)
	participantsJSON, err := json.Marshal(conv.Participants)
	if err != nil {
		return Conversation{}, err
	}

	_, err = db.c.Exec("UPDATE conversations SET participants = ? WHERE id = ?", 
		string(participantsJSON), cid)
	if err != nil {
		return Conversation{}, err
	}

	return conv, nil
}

func (db *appdbimpl) leaveGroup(cid string, user User) (Conversation, error) {
	conv, err := db.getConversation(cid)
	if err != nil {
		return Conversation{}, err
	}

	// Remove user from participants
	var newParticipants []User
	for _, p := range conv.Participants {
		if p.UId != user.UId {
			newParticipants = append(newParticipants, p)
		}
	}
	conv.Participants = newParticipants

	participantsJSON, err := json.Marshal(conv.Participants)
	if err != nil {
		return Conversation{}, err
	}

	_, err = db.c.Exec("UPDATE conversations SET participants = ? WHERE id = ?", 
		string(participantsJSON), cid)
	if err != nil {
		return Conversation{}, err
	}

	return conv, nil
}

func (db *appdbimpl) setGroupName(cid string, name string) (Conversation, error) {
	// Note: Group name is not in the current schema, would need to add a name column
	return db.getConversation(cid)
}

func (db *appdbimpl) setGroupPhoto(cid string, picture string) (Conversation, error) {
	_, err := db.c.Exec("UPDATE conversations SET picture = ? WHERE id = ?", picture, cid)
	if err != nil {
		return Conversation{}, err
	}
	return db.getConversation(cid)
} 