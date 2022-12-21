package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// TODO: IMPLEMENT THIS STUFF
func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, actx reqcontext.AuthRequestContext) {
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

	photoslimit := r.URL.Query().Get("photos-limit")
	intlimit := 0
	if len(photoslimit) > 0 {
		intlimit, err = strconv.Atoi(photoslimit)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	var res []int

	res, err = rt.db.GetStream(givenname, intlimit)

	// Maybe empty result may throw an error here? Will see.
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't get photo ids for stream")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't encode photo ids in response")
		return
	}
}
