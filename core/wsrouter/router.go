package wsrouter

import (
	"encoding/json"
	"errors"
	"net/http"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"game_lib/logging"
	"algameserver/core/online"
	"game_lib/dbconn"
	"dgit.cc/game-router/grpc"
	"game_lib/game"


	"game_lib/session"
		"algameserver/core/services"
)

// Debug flag
var Debug = false

// ErrRouteNotFound is returned when router can not find matching route
var ErrRouteNotFound = errors.New("type route not found")

var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WsHandlerFunc handler for websocket
type WsHandlerFunc func(s *session.Session, p []byte) (err error)

// Route websocket route
type Route struct {
	Action  string
	Handler WsHandlerFunc
}

// Router websocket router
type Router struct {
	Routes []Route
	Hub    *Hub
}

// New Router Instance
func New() *Router {
	hub := NewHub()
	go hub.Run()

	//go func() {
	//	tick := time.NewTicker(time.Second * 10)
	//	for {
	//		select {
	//		case <-tick.C:
	//			fmt.Println("...")
	//			redisConn := dbconn.RedisClient()
	//			countService := online.NewUserCount(redisConn)
	//			count := countService.Count()
	//			hub.broadcast <- []byte(strconv.Itoa(count))
	//		}
	//	}
	//}()
	router := &Router{
		Hub: hub,
	}
	go router.BroadcastAction()
	return router
}

// AddRoute add route into router
func (r *Router) AddRoute(act string, h WsHandlerFunc) {
	r.Routes = append(r.Routes, Route{
		Action:  act,
		Handler: h,
	})
}

// FindRoute find matching route
func (r *Router) FindRoute(act string) (WsHandlerFunc, error) {
	for _, v := range r.Routes {
		if v.Action == act {
			return v.Handler, nil
		}
	}
	return nil, ErrRouteNotFound
}

// PayloadAction in websocket message
type PayloadAction struct {
	Action string `json:"action"`
}

// HandshakeAndRun websocket handshake and serve
func (r *Router) HandshakeAndRun(c *gin.Context) {
	token, ok := c.GetQuery("token")
	if !ok {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	gameType, ok := c.GetQuery("gameType")
	if !ok {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	// upgrade
	conn, err := websocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logging.S().Error("Failed to set websocket upgrade:", err)
		return
	}
	defer conn.Close()

	// 線上人數統計
	redisConn := dbconn.RedisClient()
	countService := online.NewUserCount(redisConn)

	err = countService.IncrCount()
	if err != nil {
		logging.L().Error(fmt.Sprintf("countService.IncrCount() %s", err.Error()))
		return
	}
	defer countService.DecrCount()

	// get acc user info
	accUserInfo, authErr := grpc.AuthToken(token)
	if authErr != nil {
		logging.L().Error(fmt.Errorf("authToken err %v", authErr).Error())
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	if accUserInfo == nil {
		logging.L().Error(fmt.Errorf("accUserInfo is null err").Error())
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	session := session.NewSession(conn)
	session.Token = token
	session.UserID = accUserInfo.UserID
	session.Username = accUserInfo.Username
	session.AgentID = accUserInfo.Aid
	session.Domain = accUserInfo.Domain
	sysGameType, err:= game.GetGameType(gameType)
	if err!=nil{
		logging.L().Error(err.Error())
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	session.GameType =  sysGameType
	logging.L().Info("get room id")
	roomId, err := services.Lobby.GetRoomId(session.GameType, session)
	if err != nil {
		logging.L().Error(fmt.Sprintf("GetRoomId err #%v", err.Error()))
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("GetRoomId err #%v", err.Error()))
		return
	}
	if roomId == nil {
		logging.L().Error(fmt.Sprintf("GetRoomId roomId is null, token %s", session.Token))
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("GetRoomId roomId is null, token %s", session.Token))
		return
	}
	session.RoomId = roomId

	r.Hub.register <- session
	defer func() {
		logging.L().Info("remove room id")
		err := services.Lobby.RemoveSession(session.GameType, session)
		if err != nil {
			logging.L().Error(fmt.Sprintf("remove session #%v", err))
		}
		r.Hub.unregister <- session
	}()

	// keeping listen message
	for {
		_, payload, err := conn.ReadMessage()

		if err != nil {
			closeError, ok := err.(*websocket.CloseError)
			if ok {
				switch closeError.Code {
				case websocket.CloseGoingAway:
					logging.L().Info("socket closed")
					return
				}
			}
			logging.S().Error("read messag error:", err)
			break
		}

		pa := PayloadAction{}
		err = json.Unmarshal(payload, &pa)
		if err != nil {
			logging.S().Error("read json error:", err)
			er := conn.WriteJSON(map[string]string{
				"message": err.Error(),
			})
			if er != nil {
				logging.S().Error("write json error:", err)
			}
			continue
		}

		handlerFunc, err := r.FindRoute(pa.Action)
		if err != nil {
			logging.S().Error("route error:", err)
			er := conn.WriteJSON(map[string]string{
				"message": err.Error(),
			})
			if er != nil {
				logging.S().Error("write json error:", err)
			}
			continue
		}

		err = handlerFunc(session, payload)
		if err != nil {
			logging.S().Error("handle error:", err)
			continue
		}
	}
}

// Broadcast message to all session
func (r *Router) Broadcast(msg []byte) {
	r.Hub.broadcast <- msg
}

const (
	ChannelWsBroadcast              = "ws.broadcast.*"
	ChannelWsBroadcastMessaage      = "ws.broadcast.message"
	ChannelWsBroadcastCurrentPeriod = "ws.broadcast.current_period"
	ChannelWsBroadcastLatestDraw    = "ws.broadcast.latest_draw"
	ChannelWsBroudcaseTouchTable    = "ws.broadcast.touch_table"
)

// Action with broadcast action name and data
type Action struct {
	Name string
	Data string
}

// BroadcastAction run with redis pubsub
func (r *Router) BroadcastAction() {
	redisC := dbconn.RedisClient()
	pubsub := redisC.PSubscribe("ws.broadcast.*")
	ch := pubsub.Channel()

	for msg := range ch {
		switch msg.Channel {
		case ChannelWsBroadcastMessaage:
			r.Hub.broadcastAction <- &Action{
				Name: BroadcastMessage,
				Data: msg.Payload,
			}
		case ChannelWsBroadcastCurrentPeriod:
			r.Hub.broadcastAction <- &Action{
				Name: BroadcastCurrentPeriod,
				Data: msg.Payload,
			}
		case ChannelWsBroadcastLatestDraw:
			r.Hub.broadcastAction <- &Action{
				Name: BroadcastLatestDraw,
				Data: msg.Payload,
			}
		case ChannelWsBroudcaseTouchTable:
			r.Hub.broadcastAction <- &Action{
				Name: BroadcastTouchTable,
				Data: msg.Payload,
			}
		default:
			logging.S().Warn("receive no matching action")
		}
	}
}
