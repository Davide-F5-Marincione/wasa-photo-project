package api

import (
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, actx reqcontext.AuthRequestContext) {
	photoid := ps.ByName("id")

	intphotoid, err := strconv.Atoi(photoid)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	photodetails, err := rt.db.GetPhotoDetails(intphotoid)

	// Probably id used
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Only the author can delete, are we the author?
	if photodetails.Author != actx.ReqUserHandle {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = rt.db.RemovePhoto(intphotoid)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't remove photo")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
