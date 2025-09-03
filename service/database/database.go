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
	"strings"
)
type User struct {
	UId string `json:"id"`
	Username string `json:"username"`
	Picture string `json:"picture,omitempty"`
}

type Message struct {
	Id string `json:"id"`
	SenderId string `json:"senderId"`
	Text string `json:"text"`
	ImageUrl string `json:"imageUrl,omitempty"`
	SenderUsername string `json:"senderUsername"`
	Time string `json:"time,omitempty"`
	Reactions map[string]interface{} `json:"reactions,omitempty"`
	IsRead bool `json:"isRead,omitempty"`
	ReadBy []string `json:"readBy,omitempty"`
}

type Conversation struct {
	CId string `json:"id"`
	Name string `json:"name"`
	Picture string `json:"picture"`
	Participants []User `json:"participants"`
	LastMessage *Message `json:"lastMessage,omitempty"`
	LastMessageTime string `json:"lastMessageTime,omitempty"`
	UnreadCount int `json:"unreadCount,omitempty"`
}
// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	Ping() error
	DoLogin(user User)
	ListUsers(username string) ([]User, error)
	SetMyUserName(username string) (User, error)
	SetMyPhoto(picture string) (User, error)
	CreateConversation(participants []User, name string) (Conversation, error)
	GetMyConversations(user User) ([]Conversation, error)
	GetConversation(cid string) (Conversation, error)
	AddToGroup(cid string, user User) (Conversation, error)
	LeaveGroup(cid string, user User) (Conversation, error)
	SetGroupName(cid string, name string) (Conversation, error)
	SetGroupPhoto(cid string, picture string) (Conversation, error)
	SendMessage(cid string, user User, message string) (Conversation, error)
	SendMessageWithImage(cid string, user User, message string, imageUrl string) (Conversation, error)
	DeleteMessage(cid string, user User, mid string) (Conversation, error)
	ForwardMessage(cid string, user User, mid string) (Conversation, error)
	ReactToMessage(cid string, user User, mid string, emoji string) (Conversation, error)
	RemoveReaction(cid string, user User, mid string, emoji string) (Conversation, error)
	CommentMessage(cid string, user User, mid string, comment string) (Conversation, error)
	UncommentMessage(cid string, user User, mid string, commentId string) (Conversation, error)
	MarkMessageAsRead(messageId string, userId string) error
	GetUnreadCount(conversationId string, userId string) (int, error)
	GetContextReply() (string, error)
	AddContact(user User, contact User) (User, error)
	ListContacts(user User) ([]User, error)
	RemoveContact(user User, contact User) (User, error)
	GetAllConversations() ([]Conversation, error)
	GetRawDB() *sql.DB
}

type appdbimpl struct {
	c *sql.DB
}

func (db *appdbimpl) GetRawDB() *sql.DB {
	return db.c
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}
	// Enable WAL mode for better concurrency
	_, err := db.Exec("PRAGMA journal_mode = WAL")
	if err != nil {
		return nil, fmt.Errorf("error enabling WAL mode: %w", err)
	}
	
	// Set proper SQLite settings for concurrency
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, fmt.Errorf("error enabling foreign keys: %w", err)
	}
	
	_, err = db.Exec("PRAGMA busy_timeout = 30000") // 30 second timeout
	if err != nil {
		return nil, fmt.Errorf("error setting busy timeout: %w", err)
	}
	
	_, err = db.Exec("PRAGMA synchronous = NORMAL")
	if err != nil {
		return nil, fmt.Errorf("error setting synchronous mode: %w", err)
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
		participants TEXT NOT NULL,
		name TEXT,
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
		image_url TEXT,
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
		FOREIGN KEY(message_id) REFERENCES messages(id) ON DELETE CASCADE,
		UNIQUE(message_id, sender_id, emoji)
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

	readStatusTable := `CREATE TABLE IF NOT EXISTS read_status (
		id TEXT PRIMARY KEY,
		message_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		read_at INTEGER NOT NULL,
		FOREIGN KEY(message_id) REFERENCES messages(id) ON DELETE CASCADE,
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE(message_id, user_id)
	);`
	if _, err := db.Exec(readStatusTable); err != nil {
		return nil, fmt.Errorf("error creating read_status table: %w",err)
	}

	// Add image_url column to messages table if it doesn't exist
	// This is a migration for existing databases
	_, err = db.Exec("ALTER TABLE messages ADD COLUMN image_url TEXT")
	if err != nil {
		// Column might already exist, check if it's a "duplicate column" error
		// If it's not a duplicate column error, return the error
		if !strings.Contains(err.Error(), "duplicate column") {
			return nil, fmt.Errorf("error adding image_url column: %w", err)
		}
		// Otherwise, column already exists, continue
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
