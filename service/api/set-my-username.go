package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, actx reqcontext.AuthRequestContext) {
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

	var newname string
	err = json.NewDecoder(r.Body).Decode(&newname)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	check, err := regexp.Match("^.{4,32}$", []byte(newname))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't regexp username")
		return
	}
	if !check {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.UpdateUsername(actx.ReqUserHandle, newname)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't update username")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
