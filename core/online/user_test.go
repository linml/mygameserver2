package online

import (
	"testing"

	"dgit.cc/PcEgg/gameserv/core/dbconn"
	"github.com/stretchr/testify/assert"
)

func TestUserCount(t *testing.T) {
	dbconn.RedisDial(":6379", 0)
	redisConn := dbconn.RedisClient()
	countService := NewUserCount(redisConn)
	countService.ResetCount()

	t.Run("Online User Count IncrCount", func(t *testing.T) {
		countService.IncrCount()
		assert.True(t, countService.Count() == 1)
	})

	t.Run("Online User Count DecrCount", func(t *testing.T) {
		countService.DecrCount()
		assert.True(t, countService.Count() == 0)

		countService.DecrCount()
		assert.True(t, countService.Count() == 0)
	})

	t.Run("Online User Count ResetCount", func(t *testing.T) {
		countService.IncrCount()
		countService.IncrCount()
		countService.IncrCount()
		countService.ResetCount()
		assert.True(t, countService.Count() == 0)
	})
}
