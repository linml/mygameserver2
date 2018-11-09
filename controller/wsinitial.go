package controller

import (
	"algameserver/core/wsreply"
	"game_lib/session"
)

// Initial action: Initial
func Initial(s *session.Session, p []byte) (err error) {

	resp := wsreply.NewSuccessReply(
		"initial",
		map[string]interface{}{
			"msg": "initial",
		},
	)
	return s.Conn.WriteJSON(resp)
}
