package database

import (
	"database/sql"

	"github.com/gofrs/uuid"
)

func (db *appdbimpl) AddContact(user User, contact User) (User, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return User{}, err
	}

	_, err = db.c.Exec("INSERT INTO contacts (id, user_id, contact_id) VALUES (?, ?, ?)",
		id.String(), user.UId, contact.UId)
	if err != nil {
		return User{}, err
	}

	return contact, nil
}

func (db *appdbimpl) ListContacts(user User) ([]User, error) {
	rows, err := db.c.Query(`
		SELECT u.id, u.username, u.picture 
		FROM users u 
		JOIN contacts c ON u.id = c.contact_id 
		WHERE c.user_id = ?`, user.UId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []User
	for rows.Next() {
		var contact User
		var picture sql.NullString
		if scanErr := rows.Scan(&contact.UId, &contact.Username, &picture); scanErr != nil {
			return nil, scanErr
		}
		// Handle NULL picture values
		if picture.Valid {
			contact.Picture = picture.String
		} else {
			contact.Picture = ""
		}
		contacts = append(contacts, contact)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}

func (db *appdbimpl) RemoveContact(user User, contact User) (User, error) {
	_, err := db.c.Exec("DELETE FROM contacts WHERE user_id = ? AND contact_id = ?",
		user.UId, contact.UId)
	if err != nil {
		return User{}, err
	}
	return contact, nil
}
