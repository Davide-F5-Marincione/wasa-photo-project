package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) findUserHandle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := r.URL.Query().Get("user-name")

	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	namebase := r.URL.Query().Get("name-base")
	handlebase := r.URL.Query().Get("handle-base")

	if namebase != "" && handlebase == "" || namebase == "" && handlebase != "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	namesandhandles, err := rt.db.FindSimilar(username, handlebase, namebase)

	// Maybe empty result may throw an error here? Will see.
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rt.baseLogger.WithError(err).Error("can't get username search results from db")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(namesandhandles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rt.baseLogger.WithError(err).Error("can't encode user profile response")
		return
	}
}
