package api

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"regexp"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

type loginResponse struct {
	name   string `json:"resp-username"`
	handle string `json:"resp-userhandle"`
	auth   int64  `json:"resp-authtoken"`
}

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userhandle := r.URL.Query().Get("user-handle")
	check, err := regexp.Match("^.{1,27}#\\d{4}$", []byte(userhandle))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("can't regexp user-handle")
		return
	}
	if !check {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	details, err := rt.db.GetUserDetails(userhandle)
	status := http.StatusOK
	if err != nil {
		// Hey, this guy doesn't exist in the db, let's register him
		rand_auth := rand.Int63()
		for !rt.db.CheckAuthFree(rand_auth) {
			rand_auth = rand.Int63()
		}
		details = database.UserDetails{Handle: userhandle, Name: userhandle, Auth: rand_auth}
		err = rt.db.InsertUser(details)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("can't register new user")
			return
		}
		status = http.StatusCreated
	}

	resp := loginResponse{name: details.Name, handle: details.Handle, auth: details.Auth}

	// TODO: return stuff
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("can't encode login response")
		return
	}

	w.WriteHeader(status)
	println("out")
}

//// This stuff is kept here because.

// // getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
// func (rt *_router) getHelloWorld(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	w.Header().Set("content-type", "text/plain")
// 	_, _ = w.Write([]byte("Hello World!"))
// }

// // getContextReply is an example of HTTP endpoint that returns "Hello World!" as a plain text. The signature of this
// // handler accepts a reqcontext.RequestContext (see httpRouterHandler).
// func (rt *_router) getContextReply(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
// 	w.Header().Set("content-type", "text/plain")
// 	_, _ = w.Write([]byte("Hello World!"))
// }
