package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.POST("/users", rt.wrap(rt.doLogin))
	rt.router.PUT("/users/:user-handle", rt.authWrap(rt.setMyUserName))
	rt.router.PUT("/users/:user-handle/bans/:other-handle", rt.authWrap(rt.banUser))
	rt.router.DELETE("/users/:user-handle/bans/:other-handle", rt.authWrap(rt.unbanUser))
	rt.router.PUT("/users/:user-handle/follows/:other-handle", rt.authWrap(rt.followUser))
	rt.router.DELETE("/users/:user-handle/follows/:other-handle", rt.authWrap(rt.unfollowUser))
	rt.router.POST("/photos", rt.authWrap(rt.uploadPhoto))
	rt.router.GET("/photos/:id/raw", rt.authWrap(rt.getRawPhoto))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
