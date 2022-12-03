package api

import (
	"net/http"
	"strconv"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// httpRouterHandler is the signature for functions that accepts a reqcontext.RequestContext in addition to those
// required by the httprouter package.
type httpRouterHandler func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext)

// wrap parses the request and adds a reqcontext.RequestContext instance related to the request.
func (rt *_router) wrap(fn httpRouterHandler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		reqUUID, err := uuid.NewV4()
		if err != nil {
			rt.baseLogger.WithError(err).Error("can't generate a request UUID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var ctx = reqcontext.RequestContext{
			ReqUUID: reqUUID,
		}

		// Create a request-specific logger
		ctx.Logger = rt.baseLogger.WithFields(logrus.Fields{
			"reqid":     ctx.ReqUUID.String(),
			"remote-ip": r.RemoteAddr,
		})

		// Call the next handler in chain (usually, the handler function for the path)
		fn(w, r, ps, ctx)
	}
}

// httpAuthRouterHandler is the signature for functions that accepts a reqcontext.AuthRequestContext in addition to those
// required by the httprouter package.
type httpAuthRouterHandler func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.AuthRequestContext)

// Maybe this is a bit too much
func (rt *_router) authWrap(fn httpAuthRouterHandler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authheader := r.Header.Get("Authorization")

		if !strings.HasPrefix(authheader, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		authstr := strings.TrimPrefix(authheader, "Bearer ")

		auth, err := strconv.Atoi(authstr)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			rt.baseLogger.WithError(err).Error("Authorization couldn't be cast to int")
			return
		}

		details, err := rt.db.GetUserDetailsAuth(auth)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			rt.baseLogger.WithError(err).Error("couldn't check for authorization validity")
			return
		}

		reqUUID, err := uuid.NewV4()
		if err != nil {
			rt.baseLogger.WithError(err).Error("can't generate a request UUID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var actx = reqcontext.AuthRequestContext{
			ReqUUID:       reqUUID,
			ReqUserHandle: details.Handle,
			ReqUserAuth:   details.Auth,
		}

		actx.Logger = rt.baseLogger.WithFields(logrus.Fields{
			"user-handle": details.Handle,
			"reqid":       actx.ReqUUID.String(),
			"remote-ip":   r.RemoteAddr,
		})

		fn(w, r, ps, actx)
	}
}
