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
	Name string `json:"resp-username"`
	Auth int    `json:"resp-authtoken"`
}

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var username string
	err := json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("can't decode user-name")
		return
	}

	check, err := regexp.Match("^.{4,32}$", []byte(username))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("can't regexp user-name")
		return
	}
	if !check {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	details, err := rt.db.GetUserDetails(username)
	if err != nil {
		// Hey, this guy doesn't exist in the db, let's register him
		randauth := rand.Int63()
		for rt.db.CheckAuth(int(randauth)) != "" {
			randauth = rand.Int63()
		}
		details = database.UserDetails{Name: username, Auth: int(randauth)}
		err = rt.db.InsertUser(details)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("can't register new user")
			return
		}
		w.WriteHeader(http.StatusCreated)
	}

	resp := loginResponse{Name: details.Name, Auth: details.Auth}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("can't encode login response")
		return
	}
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
