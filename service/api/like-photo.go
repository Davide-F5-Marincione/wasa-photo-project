package api

import (
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, actx reqcontext.AuthRequestContext) {
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

	// We may be banned from the author!
	if rt.db.CheckBan(photodetails.Author, actx.ReqUserHandle) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// We cannot like our own photo!
	if resuser.Handle == photodetails.Author {
		w.WriteHeader(http.StatusConflict)
		return
	}

	//Check if already liked
	if rt.db.CheckLike(actx.ReqUserHandle, intphotoid) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = rt.db.InsertLike(actx.ReqUserHandle, intphotoid)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't like photo")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
