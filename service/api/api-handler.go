package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.POST("/users", rt.wrap(rt.doLogin))
	rt.router.GET("/users", rt.findUser)
	rt.router.PUT("/users/:user-name", rt.authWrap(rt.setMyUserName))

	rt.router.PUT("/users/:user-name/bans/:other-name", rt.authWrap(rt.banUser))
	rt.router.DELETE("/users/:user-name/bans/:other-name", rt.authWrap(rt.unbanUser))
	rt.router.GET("/users/:user-name/bans/:other-name", rt.authWrap(rt.checkBan))

	rt.router.PUT("/users/:user-name/follows/:other-name", rt.authWrap(rt.followUser))
	rt.router.DELETE("/users/:user-name/follows/:other-name", rt.authWrap(rt.unfollowUser))
	rt.router.GET("/users/:user-name/follows/:other-name", rt.authWrap(rt.checkFollow))

	rt.router.GET("/users/:user-name", rt.authWrap(rt.getUserProfile))
	rt.router.GET("/users/:user-name/stream", rt.authWrap(rt.getMyStream))

	rt.router.POST("/photos", rt.authWrap(rt.uploadPhoto))
	rt.router.GET("/photos/:photo-id", rt.authWrap(rt.getPhotoPost))
	rt.router.GET("/photos/:photo-id/raw", rt.authWrap(rt.getRawPhoto))
	rt.router.DELETE("/photos/:photo-id", rt.authWrap(rt.deletePhoto))

	rt.router.PUT("/photos/:photo-id/likes/:user-name", rt.authWrap(rt.likePhoto))
	rt.router.DELETE("/photos/:photo-id/likes/:user-name", rt.authWrap(rt.unlikePhoto))

	rt.router.POST("/photos/:photo-id/comments", rt.authWrap(rt.commentPhoto))
	rt.router.DELETE("/photos/:photo-id/comments/:comment-id", rt.authWrap(rt.uncommentPhoto))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
