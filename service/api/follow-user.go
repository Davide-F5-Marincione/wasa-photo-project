package api

import (
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, actx reqcontext.AuthRequestContext) {
	givenhandle := ps.ByName("user-handle")

	resuser, err := rt.db.GetUserDetails(givenhandle)

	// Probably bad user handle used
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// The given authorization is for another user!
	if resuser.Handle != actx.ReqUserHandle {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	otherhandle := ps.ByName("other-handle")

	resuser, err = rt.db.GetUserDetails(otherhandle)

	// Probably bad user handle used
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Check if trying to follow self
	if otherhandle == actx.ReqUserHandle {
		w.WriteHeader(http.StatusConflict)
		return
	}

	// Check if already banned
	if rt.db.CheckFollow(actx.ReqUserHandle, otherhandle) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = rt.db.InsertFollow(actx.ReqUserHandle, otherhandle)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't follow user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
