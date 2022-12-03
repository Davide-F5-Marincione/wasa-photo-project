package api

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// Part of this function got inspired by https://stackoverflow.com/questions/25959386/how-to-check-if-a-file-is-a-valid-image
// just want to be fair.
func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, actx reqcontext.AuthRequestContext) {
	err := r.ParseMultipartForm(52428800 + 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	title := r.PostFormValue("title")
	photofile, header, err := r.FormFile("photo") // Why doesn't there exist PostFormFile?
	defer photofile.Close()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if header.Size >= 52428800 {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	photobuff := make([]byte, 512)
	if _, err = photofile.Read(photobuff); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't read initial buffer")
		return
	}

	contenttype := http.DetectContentType(photobuff)

	if !strings.HasPrefix(contenttype, "image/") {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	contenttype = strings.TrimPrefix(contenttype, "image/")

	// Everything else should be alright
	if contenttype == "x-icon" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	check, err := regexp.Match("^.{4,64}$", []byte(title))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't regexp title")
		return
	}
	if !check {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	remainingsize := header.Size - 512 //maybe we've read the whole image already
	if remainingsize > 0 {
		leftbuff := make([]byte, remainingsize)

		if _, err = photofile.Read(leftbuff); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			actx.Logger.WithError(err).Error("can't read photo buffer")
			return
		}

		photobuff = append(photobuff, leftbuff...)
	}

	id, err := rt.db.InsertPhoto(actx.ReqUserHandle, title, photobuff)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't insert photo in db")
	}

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't encode photo id response")
		return
	}
}
