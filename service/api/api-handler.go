package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
    r := gin.Default()

    r.POST("/session", doLogin)

    // Users
    r.GET("/users", listUsers)
    r.PUT("/users/:id", setMyUserName)
    r.PUT("/users/:id/photo", setMyPhoto)

    // Conversations
    r.POST("/users/:id/conversations", createGroup)
    r.GET("/users/:id/conversations", getMyConversations)
    r.GET("/users/:id/conversations/:convId", getConversation)
    r.POST("/users/:id/conversations/:convId/members", addtoGroup)
    r.DELETE("/users/:id/conversations/:convId/members", leaveGroup)
    r.PUT("/users/:id/conversations/:convId/name", setGroupName)
    r.PUT("/users/:id/conversations/:convId/photo", setGroupPhoto)

    // Messages
    r.POST("/users/:id/conversations/:convId/messages", sendMessage)
    r.DELETE("/users/:id/conversations/:convId/messages/:msgId", deleteMessage)
    r.POST("/users/:id/conversations/:convId/messages/:msgId/forward", forwardMessage)

    // Reactions
    r.POST("/users/:id/conversations/:convId/messages/:msgId/reaction", reactToMessage)
    r.DELETE("/users/:id/conversations/:convId/messages/:msgId/reaction/:emoji", removeReaction)

    // Comments
    r.POST("/users/:id/conversations/:convId/messages/:msgId/comment", commentMessage)
    r.DELETE("/users/:id/conversations/:convId/messages/:msgId/comment/:commentId", uncommentMessage)
    r.GET("/context", getContextReply)

    // Contacts
    r.POST("/users/:id/contacts", addContact)
    r.GET("/users/:id/contacts", listContacts)
    r.DELETE("/users/:id/contacts/:contactId", removeContact)

    // WebSocket endpoint for realtime updates
    r.GET("/ws", serveWs)

    r.Run() // Start the server
	return rt.router
}

