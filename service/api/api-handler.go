package api

import (
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// httpRouterHandler converts a handler that needs reqcontext.RequestContext into a standard httprouter.Handle
func (rt *_router) httpRouterHandler(fn func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := reqcontext.NewRequestContext(rt.baseLogger, r)
		fn(w, r, ps, ctx)
	}
}

func (rt *_router) Handler() http.Handler {
	r := rt.router

	// Session and users
	r.POST("/session", rt.wrap(rt.doLogin))
	r.GET("/users", rt.wrap(rt.listUsers))

	// User specific routes
	r.PUT("/users/:id", rt.wrap(rt.setMyUserName))
	r.PUT("/users/:id/photo", rt.wrap(rt.setMyPhoto))
	r.GET("/users/:id/context", rt.wrap(rt.getContextReply))

	// Conversations
	r.POST("/users/:id/conversations", rt.wrap(rt.createConversation))
	r.GET("/users/:id/conversations", rt.wrap(rt.getMyConversations))
	r.GET("/users/:id/conversations/:convId", rt.wrap(rt.getConversation))
	r.POST("/users/:id/conversations/:convId/members", rt.wrap(rt.addtoGroup))
	r.DELETE("/users/:id/conversations/:convId/members", rt.wrap(rt.leaveGroup))
	r.PUT("/users/:id/conversations/:convId/name", rt.wrap(rt.setGroupName))
	r.PUT("/users/:id/conversations/:convId/photo", rt.wrap(rt.setGroupPhoto))

	// Messages
	r.POST("/users/:id/conversations/:convId/messages", rt.wrap(rt.sendMessage))
	r.DELETE("/users/:id/conversations/:convId/messages/:msgId", rt.wrap(rt.deleteMessage))
	r.POST("/users/:id/conversations/:convId/messages/:msgId/forward", rt.wrap(rt.forwardMessage))

	// Reactions
	r.POST("/users/:id/conversations/:convId/messages/:msgId/reaction", rt.wrap(rt.reactToMessage))
	r.DELETE("/users/:id/conversations/:convId/messages/:msgId/reaction/:emoji", rt.wrap(rt.removeReaction))

	// Comments
	r.POST("/users/:id/conversations/:convId/messages/:msgId/comments", rt.wrap(rt.commentMessage))
	r.DELETE("/users/:id/conversations/:convId/messages/:msgId/comments/:commentId", rt.wrap(rt.deleteComment))

	// Contacts
	r.POST("/users/:id/contacts", rt.wrap(rt.addContact))
	r.GET("/users/:id/contacts", rt.wrap(rt.listContacts))
	r.DELETE("/users/:id/contacts/:contactId", rt.wrap(rt.removeContact))

	// WebSocket and health check
	r.GET("/ws", rt.wrap(rt.serveWs))
	r.GET("/liveness", rt.wrap(rt.liveness))

	return r
}