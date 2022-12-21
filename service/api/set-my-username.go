package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, actx reqcontext.AuthRequestContext) {
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

	err = rt.db.UpdateUsername(actx.ReqUserName, newname)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't update username")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
