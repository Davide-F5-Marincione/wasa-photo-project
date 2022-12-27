package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

type photoPostResponse struct {
	Title    string                     `json:"photo-title"`
	Author   string                     `json:"photo-author"`
	Comments []database.CommentShow     `json:"comments-running-batch"`
	Likes    []database.UserAndDatetime `json:"likes-running-batch"`
	Date     string                     `json:"photo-date"`
	Liked    bool                       `json:"liked"`
}

func (rt *_router) getPhotoPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params, actx reqcontext.AuthRequestContext) {
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

	// We may be banned from seeing this!
	if rt.db.CheckBan(photodetails.Author, actx.ReqUserName) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	commentslimit := r.URL.Query().Get("comments-limit")
	intcommentlimit := 0
	if len(commentslimit) > 0 {
		intcommentlimit, err = strconv.Atoi(commentslimit)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	comments, err := rt.db.GetPhotoComments(intphotoid, intcommentlimit)

	// Maybe empty result may throw an error here? Will see.
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't get comments for photo show")
		return
	}

	likesbase := r.URL.Query().Get("likes-base")
	likes, err := rt.db.GetPhotoLikes(intphotoid, likesbase)

	// Maybe empty result may throw an error here? Will see.
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't get likers for photo show")
		return
	}

	res := photoPostResponse{Title: photodetails.Title, Author: photodetails.Author, Date: photodetails.UploadDate, Comments: comments, Likes: likes, Liked: rt.db.CheckLike(actx.ReqUserName, intphotoid)}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't encode photo show response")
		return
	}
}
