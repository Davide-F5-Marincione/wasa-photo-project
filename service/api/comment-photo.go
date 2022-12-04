package api

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, actx reqcontext.AuthRequestContext) {

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

	var content string
	err = json.NewDecoder(r.Body).Decode(&content)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	check, err := regexp.Match("^.+$", []byte(content))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't regexp comment")
		return
	}
	if !check {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// A comment's content can't be empty or greater than 256 characters!
	if len(content) > 256 || len(content) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := rt.db.InsertComment(actx.ReqUserHandle, content, intphotoid)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't insert comment")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't encode comment id response")
		return
	}
}
