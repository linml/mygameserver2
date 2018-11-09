package controller

import (
	"github.com/gorilla/websocket"
	"algameserver/core/services"
	"game_lib/bet"
	"encoding/json"
	"game_lib/session"
)

// BetHandler websocket bet
func BetHandler(s *session.Session, p []byte) (err error) {

	in := &bet.BetRequest{

	}

	json.Unmarshal(p, in)

	result, err := services.BetService.Bet(in)

	jsonBytes, err := json.Marshal(result)

	return s.Conn.WriteMessage(websocket.TextMessage, jsonBytes)
}
