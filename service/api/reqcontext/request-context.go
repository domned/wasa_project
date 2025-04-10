/*
Package reqcontext contains the request context. Each request will have its own instance of RequestContext filled by the
middleware code in the api-context-wrapper.go (parent package).

Each value here should be assumed valid only per request only, with some exceptions like the logger.
*/
package reqcontext

import (
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

// RequestContext is the context of the request, for request-dependent parameters
type RequestContext struct {
	// ReqUUID is the request unique ID
	ReqUUID uuid.UUID

	// Logger is a custom field logger for the request
	Logger logrus.FieldLogger
}

// NewRequestContext creates a new request context with the given logger and request
func NewRequestContext(logger logrus.FieldLogger, r *http.Request) RequestContext {
	reqUUID, err := uuid.NewV4()
	if err != nil {
		logger.WithError(err).Error("can't generate request UUID")
	}

	return RequestContext{
		ReqUUID: reqUUID,
		Logger: logger.WithFields(logrus.Fields{
			"reqid":     reqUUID.String(),
			"remote-ip": r.RemoteAddr,
		}),
	}
}
