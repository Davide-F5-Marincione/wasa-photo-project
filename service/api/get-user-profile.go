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
	PostIDs          []int                      `json:"photos-running-batch"`
	FollowersHandles []database.UserAndDatetime `json:"followers-running-batch"`
	FollowingHandles []database.UserAndDatetime `json:"following-running-batch"`
}

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, actx reqcontext.AuthRequestContext) {
	givenhandle := ps.ByName("user-handle")

	resuser, err := rt.db.GetUserDetails(givenhandle)

	// Probably bad user handle used
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Are we allowed to see this user?
	if rt.db.CheckBan(resuser.Handle, actx.ReqUserHandle) {
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

	if intphotolimit < 1 {
		photos, err = rt.db.GetPhotosProfile(givenhandle)
	} else {
		photos, err = rt.db.GetPhotosProfileLimit(givenhandle, intphotolimit)
	}

	// Maybe empty result may throw an error here? Will see.
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't get photo ids for profile show")
		return
	}

	followersbase := r.URL.Query().Get("followers-base")
	var followers []database.UserAndDatetime

	if len(followersbase) < 1 {
		followers, err = rt.db.GetFollowers(givenhandle)
	} else {
		followers, err = rt.db.GetFollowersLimit(givenhandle, followersbase)
	}

	// Maybe empty result may throw an error here? Will see.
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't get followers for profile show")
		return
	}

	followingbase := r.URL.Query().Get("following-base")
	var following []database.UserAndDatetime

	if len(followingbase) < 1 {
		following, err = rt.db.GetFollowing(givenhandle)
	} else {
		following, err = rt.db.GetFollowingLimit(givenhandle, followingbase)
	}

	// Maybe empty result may throw an error here? Will see.
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't get following for profile show")
		return
	}

	var res userProfileResponse

	res.PostIDs = photos
	res.FollowersHandles = followers
	res.FollowingHandles = following

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		actx.Logger.WithError(err).Error("can't encode user profile response")
		return
	}
}
