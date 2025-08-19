
package database

import (
	"encoding/json"
	"math/rand"
	"time"
	"database/sql"
	"github.com/gofrs/uuid"
)

// GetAllConversations returns all conversations in the database (for user picker)
func (db *appdbimpl) GetAllConversations() ([]Conversation, error) {
	rows, err := db.c.Query("SELECT id, participants, name, picture FROM conversations")
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
		   err := rows.Scan(&conv.CId, &participantsJSON, &name, &picture)
		   if err != nil {
			   return nil, err
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
			  err = json.Unmarshal([]byte(participantsJSON), &participantIDs)
			  if err != nil {
				  return nil, err
			  }
			  var participants []User
			  for _, uid := range participantIDs {
				  var user User
				  var upic sql.NullString
				  err := db.c.QueryRow("SELECT id, username, picture FROM users WHERE id = ?", uid).Scan(&user.UId, &user.Username, &upic)
				  if err != nil {
					  return nil, err
				  }
				  if upic.Valid {
					  user.Picture = upic.String
				  } else {
					  user.Picture = ""
				  }
				  participants = append(participants, user)
			  }
			  conv.Participants = participants
		   conversations = append(conversations, conv)
	   }
	return conversations, nil
}

var demoNames = []string{"Alice & Bob", "Charlie Group", "Delta Squad", "Echo Team", "Foxtrot Chat"}
var demoAvatars = []string{
	"https://randomuser.me/api/portraits/men/1.jpg",
	"https://randomuser.me/api/portraits/women/2.jpg",
	"https://randomuser.me/api/portraits/men/3.jpg",
	"https://randomuser.me/api/portraits/women/4.jpg",
	"https://randomuser.me/api/portraits/men/5.jpg",
}

func (db *appdbimpl) CreateConversation(participants []User) (Conversation, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return Conversation{}, err
	}

	rand.Seed(time.Now().UnixNano())
	name := demoNames[rand.Intn(len(demoNames))]
	avatar := demoAvatars[rand.Intn(len(demoAvatars))]

	participantsJSON, err := json.Marshal(participants)
	if err != nil {
		return Conversation{}, err
	}

	_, err = db.c.Exec("INSERT INTO conversations (id, participants, name, picture) VALUES (?, ?, ?, ?)", 
		id.String(), string(participantsJSON), name, avatar)
	if err != nil {
		return Conversation{}, err
	}

	return Conversation{CId: id.String(), Participants: participants, Name: name, Picture: avatar}, nil
}

func (db *appdbimpl) GetMyConversations(user User) ([]Conversation, error) {
		rows, err := db.c.Query("SELECT id, participants, name, picture FROM conversations")
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
		   err := rows.Scan(&conv.CId, &participantsJSON, &name, &picture)
		   if err != nil {
			   return nil, err
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
		   err = json.Unmarshal([]byte(participantsJSON), &participantIDs)
		   if err != nil {
			   return nil, err
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
				   err := db.c.QueryRow("SELECT id, username, picture FROM users WHERE id = ?", uid).Scan(&u.UId, &u.Username, &upic)
				   if err != nil {
					   return nil, err
				   }
				   if upic.Valid {
					   u.Picture = upic.String
				   } else {
					   u.Picture = ""
				   }
				   participants = append(participants, u)
			   }
			   conv.Participants = participants
			   conversations = append(conversations, conv)
		   }
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
	   err = json.Unmarshal([]byte(participantsJSON), &participantIDs)
	   if err != nil {
		   return Conversation{}, err
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
	   err = json.Unmarshal([]byte(participantsJSON), &participantIDs)
	   if err != nil {
		   return Conversation{}, err
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
	   err = json.Unmarshal([]byte(participantsJSON), &participantIDs)
	   if err != nil {
		   return Conversation{}, err
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