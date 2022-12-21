package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) findUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := r.URL.Query().Get("user-name")

	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	namebase := r.URL.Query().Get("name-base")

	names, err := rt.db.FindSimilar(username, namebase)

	// Maybe empty result may throw an error here? Will see.
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rt.baseLogger.WithError(err).Error("can't get username search results from db")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(names)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rt.baseLogger.WithError(err).Error("can't encode user profile response")
		return
	}
}
