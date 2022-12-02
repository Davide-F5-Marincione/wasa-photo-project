package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.POST("/users", rt.wrap(rt.doLogin))
	rt.router.PUT("/users/:user-handle", rt.authWrap(rt.setMyUserName))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
