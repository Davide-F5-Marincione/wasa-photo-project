package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

type userProfileResponse struct {
	PostIDs        []int                      `json:"photos-running-batch"`
	FollowersNames []database.UserAndDatetime `json:"followers-running-batch"`
	FollowingNames []database.UserAndDatetime `json:"following-running-batch"`
}

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, actx reqcontext.AuthRequestContext) {
	givenname := ps.ByName("user-name")

	resuser, err := rt.db.GetUserDetails(givenname)

	// Probably bad user name used
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Are we allowed to see this user?
	if rt.db.CheckBan(resuser.Name, actx.ReqUserName) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	photoslimit := r.URL.Query().Get("photos-limit")
	intphotolimit := 0
	if len(photoslimit) > 0 {
		intphotolimit, err = strconv.Atoi(photoslimit)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	var photos []int

	photos, err = rt.db.GetPhotosProfile(givenname, intphotolimit)

	// Maybe empty result may throw an error here? Will see.
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't get photo ids for profile show")
		return
	}

	followersbase := r.URL.Query().Get("followers-base")
	var followers []database.UserAndDatetime

	followers, err = rt.db.GetFollowers(givenname, followersbase)

	// Maybe empty result may throw an error here? Will see.
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't get followers for profile show")
		return
	}

	followingbase := r.URL.Query().Get("following-base")
	var following []database.UserAndDatetime

	following, err = rt.db.GetFollowing(givenname, followingbase)

	// Maybe empty result may throw an error here? Will see.
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't get following for profile show")
		return
	}

	var res userProfileResponse

	res.PostIDs = photos
	res.FollowersNames = followers
	res.FollowingNames = following

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't encode user profile response")
		return
	}
}
