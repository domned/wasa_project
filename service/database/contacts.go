package database

import (
	"github.com/gofrs/uuid"
)

func (db *appdbimpl) addContact(user User, contact User) (User, error) {
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

func (db *appdbimpl) listContacts(user User) ([]User, error) {
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
		err := rows.Scan(&contact.UId, &contact.Username, &contact.Picture)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}
	return contacts, nil
}

func (db *appdbimpl) removeContact(user User, contact User) (User, error) {
	_, err := db.c.Exec("DELETE FROM contacts WHERE user_id = ? AND contact_id = ?", 
		user.UId, contact.UId)
	if err != nil {
		return User{}, err
	}
	return contact, nil
} 