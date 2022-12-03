package api

import (
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getRawPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, actx reqcontext.AuthRequestContext) {
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

	// We may be banned from seeing this!
	if rt.db.CheckBan(photodetails.Author, actx.ReqUserHandle) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	photofile, err := rt.db.GetBlobPhoto(intphotoid)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't retrieve photo blob from db")
		return
	}

	contenttype := http.DetectContentType(photofile[:512])

	w.Header().Set("Content-Type", contenttype)
	_, err = w.Write(photofile)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't write photo in response")
		return
	}
}
