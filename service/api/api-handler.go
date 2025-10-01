package api

import (
	"net/http"
)

func (rt *_router) Handler() http.Handler {
	r := rt.router

	// Public routes (no authentication required)
	r.GET("/", rt.wrap(rt.getApiRoot))
	r.POST("/session", rt.wrap(rt.doLogin))
	r.GET("/liveness", rt.wrap(rt.liveness))

	// Authenticated routes
	r.GET("/users", rt.wrapAuth(rt.listUsers))

	// User specific routes
	r.PUT("/users/:id", rt.wrapAuth(rt.setMyUserName))
	r.PUT("/users/:id/photo", rt.wrapAuth(rt.setMyPhoto))
	r.GET("/users/:id/context", rt.wrapAuth(rt.getContextReply))

	// Conversations
	r.POST("/users/:id/conversations", rt.wrapAuth(rt.createConversation))
	r.GET("/users/:id/conversations", rt.wrapAuth(rt.getMyConversations))
	r.GET("/users/:id/conversations/:conversationId", rt.wrapAuth(rt.getConversation))
	r.GET("/conversations/all", rt.wrapAuth(rt.getAllConversations))
	r.POST("/users/:id/conversations/:conversationId/members", rt.wrapAuth(rt.addtoGroup))
	r.DELETE("/users/:id/conversations/:conversationId/members", rt.wrapAuth(rt.leaveGroup))
	r.PUT("/users/:id/conversations/:conversationId/name", rt.wrapAuth(rt.setGroupName))
	r.PUT("/users/:id/conversations/:conversationId/photo", rt.wrapAuth(rt.setGroupPhoto))

	// Messages
	r.GET("/users/:id/conversations/:conversationId/messages", rt.wrapAuth(rt.getMessages))
	r.POST("/users/:id/conversations/:conversationId/messages", rt.wrapAuth(rt.sendMessage))
	r.DELETE("/users/:id/conversations/:conversationId/messages/:messageId", rt.wrapAuth(rt.deleteMessage))
	r.POST("/users/:id/conversations/:conversationId/messages/:messageId/forward", rt.wrapAuth(rt.forwardMessage))

	// Reactions
	r.POST("/users/:id/conversations/:conversationId/messages/:messageId/reaction", rt.wrapAuth(rt.reactToMessage))
	r.DELETE("/users/:id/conversations/:conversationId/messages/:messageId/reaction/:emoji", rt.wrapAuth(rt.removeReaction))

	// Comments
	r.POST("/users/:id/conversations/:conversationId/messages/:messageId/comments", rt.wrapAuth(rt.commentMessage))
	r.DELETE("/users/:id/conversations/:conversationId/messages/:messageId/comments/:commentId", rt.wrapAuth(rt.deleteComment))

	// Contacts
	r.POST("/users/:id/contacts", rt.wrapAuth(rt.addContact))
	r.GET("/users/:id/contacts", rt.wrapAuth(rt.listContacts))
	r.DELETE("/users/:id/contacts/:contactId", rt.wrapAuth(rt.removeContact))

	// WebSocket (should this be authenticated?)
	r.GET("/ws", rt.wrap(rt.serveWs))

	return r
}
