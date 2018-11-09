package wsrouter

import (
	"fmt"

	"github.com/gorilla/websocket"
	"game_lib/logging"
	"algameserver/core/wsreply"

	"game_lib/session"
	)

// Broadcast action name
const (
	BroadcastMessage       = "message"
	BroadcastCurrentPeriod = "current_period"
	BroadcastLatestDraw    = "latest_draw"
	BroadcastTouchTable    = "touch_table"
	BroadcastTotalOnline   = "total_online"
)

// WsBroadcastHandlerFunc handler for websocket broadcast
type WsBroadcastHandlerFunc func(s *session.Session, action *Action) (err error)

// HubRoute for websocket broadcast
type HubRoute struct {
	Action  string
	Handler WsBroadcastHandlerFunc
}

// Hub collects active session
type Hub struct {
	sessions        map[*session.Session]bool
	broadcast       chan []byte
	broadcastAction chan *Action
	register        chan *session.Session
	unregister      chan *session.Session
	Routes          []HubRoute
}

// NewHub create Hub instance
func NewHub() *Hub {
	return &Hub{
		sessions:        make(map[*session.Session]bool),
		broadcast:       make(chan []byte),
		broadcastAction: make(chan *Action),
		register:        make(chan *session.Session),
		unregister:      make(chan *session.Session),
	}
}

// AddRoute add broadcast route
func (h *Hub) AddRoute(act string, handler WsBroadcastHandlerFunc) {
	h.Routes = append(h.Routes, HubRoute{
		Action:  act,
		Handler: handler,
	})
}

// FindRoute find broadcast route
func (h *Hub) FindRoute(act string) (WsBroadcastHandlerFunc, error) {
	for _, v := range h.Routes {
		if v.Action == act {
			return v.Handler, nil
		}
	}

	return nil, ErrRouteNotFound
}

// RunRoute find route and run
func (h *Hub) RunRoute(act string, action *Action) error {
	handler, err := h.FindRoute(act)
	if err != nil {
		return err
	}

	for s := range h.sessions {
		err = handler(s, action)
		if err != nil {
			return err
		}
	}

	return nil
}

// Run hub
// always operate sessions here
func (h *Hub) Run() {
	for {
		select {
		case session := <-h.register:
			h.sessions[session] = true

		case session := <-h.unregister:

			if _, ok := h.sessions[session]; ok {
				delete(h.sessions, session)
			}
		case message := <-h.broadcast:
			for s := range h.sessions {
				s.Conn.WriteMessage(websocket.TextMessage, message)
			}
		case action := <-h.broadcastAction:
			switch action.Name {
			case BroadcastCurrentPeriod:
				handler, err := h.FindRoute(BroadcastCurrentPeriod)
				if err != nil {
					logging.S().Error(BroadcastCurrentPeriod, err.Error())
					continue
				}

				for s := range h.sessions {
					err = handler(s, action)
					if err != nil {
						logging.L().Error(err.Error())
						continue
					}
				}
			case BroadcastLatestDraw:
				err := h.RunRoute(BroadcastLatestDraw, action)
				if err != nil {
					logging.L().Error(err.Error())
					continue
				}
			case BroadcastTouchTable:
			default:
				fmt.Println("Hub run default")
				for s := range h.sessions {
					broadcastMessage(s, action)
				}
			}
		}
	}
}

func broadcastMessage(s *session.Session, action *Action) {
	s.Conn.WriteJSON(wsreply.NewSuccessReply(action.Name, action.Data))
}

func (h *Hub) Sessions() map[*session.Session]bool {
	return h.sessions
}
