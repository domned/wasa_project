package api

import (
	"log"
	"net/http"
	"sync"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin in development
		// In production, you should check the origin properly
		return true
	},
}

// WebSocket message types
type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// Client represents a WebSocket client
type Client struct {
	UserID string
	Conn   *websocket.Conn
	Send   chan WSMessage
	Hub    *Hub
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan WSMessage
	mutex      sync.RWMutex
	router     *_router
}

// Global hub instance
var hub *Hub

// InitializeHub initializes the global hub with a router reference
func InitializeHub(rt *_router) {
	hub = &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan WSMessage),
		router:     rt,
	}
	go hub.run()
	rt.sysLogger.LogInfo("WebSocket hub initialized successfully")
}

// Run the hub
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()
			h.router.sysLogger.LogInfo("WebSocket client connected: " + client.UserID)
			log.Printf("WebSocket client connected: %s", client.UserID)

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mutex.Unlock()
			h.router.sysLogger.LogInfo("WebSocket client disconnected: " + client.UserID)
			log.Printf("WebSocket client disconnected: %s", client.UserID)

		case message := <-h.broadcast:
			h.mutex.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

// serveWs handles WebSocket requests from clients
func (rt *_router) serveWs(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get user ID from query parameter (try both userId and user_id for compatibility)
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		userID = r.URL.Query().Get("user_id")
	}
	if userID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	// Upgrade connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		rt.sysLogger.LogError("WebSocket upgrade failed for user " + userID + ": " + err.Error())
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Create client
	client := &Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan WSMessage, 256),
		Hub:    hub,
	}

	// Register client
	if hub != nil {
		hub.register <- client
	}

	// Update user's last_seen timestamp
	if err := rt.db.UpdateUserLastSeen(userID); err != nil {
		rt.sysLogger.LogWarn("Failed to update last_seen for WebSocket user " + userID + ": " + err.Error())
	}

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}

// readPump handles messages from the WebSocket connection
func (c *Client) readPump() {
	defer func() {
		if c.Hub != nil {
			c.Hub.unregister <- c
		}
		c.Conn.Close()
	}()

	for {
		var msg WSMessage
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle different message types
		switch msg.Type {
		case "typing_start":
			// Broadcast typing indicator to other clients
			if c.Hub != nil {
				c.Hub.broadcast <- WSMessage{
					Type: "user_typing",
					Payload: map[string]interface{}{
						"userId": c.UserID,
						"typing": true,
					},
				}
			}
		case "typing_stop":
			// Broadcast stop typing indicator
			if c.Hub != nil {
				c.Hub.broadcast <- WSMessage{
					Type: "user_typing",
					Payload: map[string]interface{}{
						"userId": c.UserID,
						"typing": false,
					},
				}
			}
		}
	}
}

// writePump handles sending messages to the WebSocket connection
func (c *Client) writePump() {
	defer c.Conn.Close()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				if err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					// Log the error but continue with return since connection is closing
					// We can't use ctx.Logger here as we don't have access to it
				}
				return
			}

			if err := c.Conn.WriteJSON(message); err != nil {
				// Log the error and close the connection
				// Client will be removed from hub by the cleanup routine
				return
			}
		}
	}
}

// BroadcastMessage broadcasts a message to all connected clients
func BroadcastMessage(msgType string, payload interface{}) {
	if hub != nil {
		hub.broadcast <- WSMessage{
			Type:    msgType,
			Payload: payload,
		}
	}
}
