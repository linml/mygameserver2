package routes

import (

	"github.com/gin-gonic/gin"
	"algameserver/core/wsrouter"
	"algameserver/controller"
)

// WebsocketRoutes for websocket action
func WebsocketRoutes(r *wsrouter.Router) {
	r.AddRoute("echo", controller.EchoHandler)
	r.AddRoute("initial", controller.Initial)
	r.AddRoute("bet", controller.BetHandler)

	r.Hub.AddRoute(wsrouter.BroadcastCurrentPeriod, controller.BroadcastCurrentPeriod)

}

// WebsocketURL for connect websocket
func WebsocketURL(router *gin.Engine, r *wsrouter.Router) {
	router.GET("/ws", r.HandshakeAndRun)
}
