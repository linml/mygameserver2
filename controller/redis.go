package controller

import (
	"net/http"

	"fmt"


	"github.com/gin-gonic/gin"
	"game_lib/dbconn"
	"game_lib/logging"
)

// RedisServer response redis info server section
func RedisServer(c *gin.Context) {
	value := dbconn.RedisClient().Get("key1")
	if value != nil {

	}
	text, err1 := value.Result()
	if err1 != nil {
		logging.S().Error(fmt.Sprintf("get key err: %s", err1.Error()))
	}
	logging.L().Info(fmt.Sprintf("key1: %s", text))
	result, err := dbconn.RedisClient().Info("server").Result()
	if err != nil {
		logging.L().Error(err.Error())
	}

	c.JSON(http.StatusOK, result)
}
