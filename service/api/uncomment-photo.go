package api

import (
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, actx reqcontext.AuthRequestContext) {

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

	// TODO-Philosophical: can we remove a comment when we are banned by the photo author?
	// We may be banned from the author!
	if rt.db.CheckBan(photodetails.Author, actx.ReqUserName) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	commentid := ps.ByName("comment-id")

	intcommentid, err := strconv.Atoi(commentid)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	comment, err := rt.db.GetComment(intphotoid, intcommentid)

	// Probably bad comment id used
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Are we the author of the comment?
	if actx.ReqUserName != comment.Author {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = rt.db.RemoveComment(intphotoid, intcommentid)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't remove comment")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
