package api

import (
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, actx reqcontext.AuthRequestContext) {
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

	othername := ps.ByName("other-name")

	_, err = rt.db.GetUserDetails(othername)

	// Probably bad user name used
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Check if already unfollowed
	if !rt.db.CheckFollow(actx.ReqUserName, othername) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = rt.db.RemoveFollow(actx.ReqUserName, othername)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't unfollow user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
