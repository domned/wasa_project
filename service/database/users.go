package database

import (
	"database/sql"

	"github.com/gofrs/uuid"
)

func (db *appdbimpl) DoLogin(user User) {
	// TODO: Implement login logic
}

func (db *appdbimpl) GetUserByID(userID string) (User, error) {
	var user User
	var picture sql.NullString

	err := db.c.QueryRow("SELECT id, username, picture FROM users WHERE id = ?", userID).Scan(
		&user.UId, &user.Username, &picture)
	if err != nil {
		return User{}, err
	}

	// Handle NULL picture values
	if picture.Valid {
		user.Picture = picture.String
	} else {
		user.Picture = ""
	}

	return user, nil
}

func (db *appdbimpl) ListUsers(username string) ([]User, error) {
	rows, err := db.c.Query("SELECT id, username, picture FROM users WHERE username LIKE ?",
		"%"+username+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0) // Initialize as empty slice instead of nil
	for rows.Next() {
		var user User
		var picture sql.NullString
		if scanErr := rows.Scan(&user.UId, &user.Username, &picture); scanErr != nil {
			return nil, scanErr
		}
		// Handle NULL picture values
		if picture.Valid {
			user.Picture = picture.String
		} else {
			user.Picture = ""
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (db *appdbimpl) SetMyUserName(username string) (User, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return User{}, err
	}

	_, err = db.c.Exec("INSERT INTO users (id, username) VALUES (?, ?)",
		id.String(), username)
	if err != nil {
		return User{}, err
	}

	return User{UId: id.String(), Username: username}, nil
}

func (db *appdbimpl) SetMyPhoto(picture string) (User, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return User{}, err
	}

	_, err = db.c.Exec("UPDATE users SET picture = ? WHERE id = ?",
		picture, id.String())
	if err != nil {
		return User{}, err
	}

	return User{UId: id.String(), Picture: picture}, nil
}
