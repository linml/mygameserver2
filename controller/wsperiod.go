package controller

import (
	"algameserver/core/wsrouter"
	"algameserver/core/wsreply"
	"game_lib/session"
)

// BroadcastCurrentPeriod broadcast current period and odds
func BroadcastCurrentPeriod(s *session.Session, action *wsrouter.Action) (err error) {

	return s.Conn.WriteJSON(wsreply.NewSuccessReply(action.Name, struct {
	}{}))
}
