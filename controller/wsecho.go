package controller

import (
	"github.com/gorilla/websocket"
	"game_lib/session"
	"fmt"
)

// EchoHandler websocket echo
func EchoHandler(s *session.Session, p []byte) (err error) {

	if s.RoomId ==nil{
		fmt.Println("roomid is null")
		return s.Conn.WriteMessage(websocket.TextMessage, []byte("room id is null"))
	}
	return s.Conn.WriteMessage(websocket.TextMessage, []byte(s.RoomId.String()))
}
