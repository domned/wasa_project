/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)
type User struct {
	UId string `json:"id"`
	Username string `json:"username"`
	Picture string `json:"picture,omitempty"`
}

type Conversation struct {
	CId string `json:"id"`
	Participants []User `json:"participants"`
}
// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	Ping() error
	doLogin (User) ()
	listUsers (Username string) ([]User, error)
	setMyUserName (Username string) (User, error)
	setMyPhoto (Picture string) (User, error)
	createConversation ([]User) (Conversation, error)
	getMyConversations (User) ([]Conversation, error)
	getConversation (CId string) (Conversation, error)
	addtoGroup (CId string, User User) (Conversation, error)
	leaveGroup (CId string, User User) (Conversation, error)
	setGroupName (CId string, Name string) (Conversation, error)
	setGroupPhoto (CId string, Picture string) (Conversation, error)
	sendMessage (CId string, User User, Message string) (Conversation, error)
	deleteMessage (CId string, User User, MId string) (Conversation, error)
	forwardMessage (CId string, User User, MId string) (Conversation, error)
	reactToMessage (CId string, User User, MId string, Emoji string) (Conversation, error)
	removeReaction (CId string, User User, MId string, Emoji string) (Conversation, error)
	commentMessage (CId string, User User, MId string, Comment string) (Conversation, error)
	uncommentMessage (CId string, User User, MId string, CommentId string) (Conversation, error)
	getContextReply () (string, error)
	addContact (User User, Contact User) (User, error)
	listContacts (User User) ([]User, error)
	removeContact (User User, Contact User) (User, error)
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}
	_, err := db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, err
	}

	//create all tables
	usersTable := `CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT NOT NULL,
		picture TEXT
	);`
	if _, err := db.Exec(usersTable); err != nil {
		return nil, fmt.Errorf("error creating users table: %w", err)
	}
	conversationsTable := `CREATE TABLE IF NOT EXISTS conversations (
		id TEXT PRIMARY KEY,
		participants TEXT NOT NULL
		picture TEXT
	);`
	if _, err := db.Exec(conversationsTable); err != nil {
		return nil, fmt.Errorf("error creating conversations table: %w", err)
	}

	messagesTable := `CREATE TABLE IF NOT EXISTS messages (
		id TEXT PRIMARY KEY,
		conversation_id TEXT NOT NULL,
		sender_id TEXT NOT NULL,
		message TEXT NOT NULL,
		FOREIGN KEY(conversation_id) REFERENCES conversations(id) ON DELETE CASCADE
	);`
	if _, err := db.Exec(messagesTable); err != nil {
		return nil, fmt.Errorf("error creating messages table: %w",err)
	}

	reactionsTable := `CREATE TABLE IF NOT EXISTS reactions (
		id TEXT PRIMARY KEY,
		message_id TEXT NOT NULL,
		sender_id TEXT NOT NULL,
		emoji TEXT NOT NULL,
		FOREIGN KEY(message_id) REFERENCES messages(id) ON DELETE CASCADE
	);`
	if _, err := db.Exec(reactionsTable); err != nil {
		return nil, fmt.Errorf("error creating reactions table: %w",err)
	}

	commentsTable := `CREATE TABLE IF NOT EXISTS comments (
		id TEXT PRIMARY KEY,
		message_id TEXT NOT NULL,
		sender_id TEXT NOT NULL,
		comment TEXT NOT NULL,
		FOREIGN KEY(message_id) REFERENCES messages(id) ON DELETE CASCADE
	);`
	if _, err := db.Exec(commentsTable); err != nil {
		return nil, fmt.Errorf("error creating comments table: %w",err)
	}
	contactsTable := `CREATE TABLE IF NOT EXISTS contacts (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		contact_id TEXT NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY(contact_id) REFERENCES users(id) ON DELETE CASCADE
	);`
	if _, err := db.Exec(contactsTable); err != nil {
		return nil, fmt.Errorf("error creating contacts table: %w",err)
	}

	_, err = db.Exec(usersTable)


}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
