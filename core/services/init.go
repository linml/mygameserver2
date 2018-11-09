package services

import (
	"game_lib/bet"
	"game_lib/dbconn"
	_ "game_lib/sys"
	"game_lib/sys"
	"game_lib/game"
)

var BetService *bet.BetService
var Lobby *sys.GameLobby

func init() {
	BetService = bet.NewBetService(dbconn.DBWrite(), dbconn.DBWrite())

	gameRoomSettings := make([]*sys.GameRoomSetting, 0)
	gameRoomSettings = append(gameRoomSettings, &sys.GameRoomSetting{
		GameType:     game.GameType_PK10,
		MaxRoomCount: 3,
		MaxUserCount: 3,
	})
	Lobby = sys.NewGameLobby(gameRoomSettings)

}
