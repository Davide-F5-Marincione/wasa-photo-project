package api

import (
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, actx reqcontext.AuthRequestContext) {
	givenname := ps.ByName("user-name")

	resuser, err := rt.db.GetUserDetails(givenname)

	// Probably bad user name used
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// The given authorization is for another user!
	if resuser.Name != actx.ReqUserName {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	photoid := ps.ByName("photo-id")

	intphotoid, err := strconv.Atoi(photoid)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	photodetails, err := rt.db.GetPhotoDetails(intphotoid)

	// Probably bad id used
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// TODO-Philosophical: can we unlike when we are banned by the author?
	// We may be banned from the author!
	if rt.db.CheckBan(photodetails.Author, actx.ReqUserName) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if already unliked
	if !rt.db.CheckLike(actx.ReqUserName, intphotoid) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = rt.db.RemoveLike(actx.ReqUserName, intphotoid)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't unlike photo")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
