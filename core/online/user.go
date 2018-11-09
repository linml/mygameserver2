package online

import (
	"github.com/go-redis/redis"

	"game_lib/logging"
)

// new online user
type UserCount struct {
	redisConn *redis.Client
}

const USER_COUNT_KEY = "online:user_count"

// init onlin user count
func NewUserCount(conn *redis.Client) *UserCount {
	return &UserCount{
		redisConn: conn,
	}
}

// 線上人數加一
func (u *UserCount) IncrCount() error {
	cmd := u.redisConn.Incr(USER_COUNT_KEY)
	if cmd.Err() != nil {
		logging.L().Info("cmd: " + cmd.Err().Error())
		return cmd.Err()
	}
	return nil
}

// 取得線上人數數量
func (u *UserCount) Count() int {
	count := u.redisConn.Get(USER_COUNT_KEY)
	onlineCount, _ := count.Int64()
	if onlineCount < 0 {
		logging.L().Info("onlineCount < 0")
	}
	return int(onlineCount)
}

// 線上人數減一
func (u *UserCount) DecrCount() {
	count := u.redisConn.Get(USER_COUNT_KEY)
	onlineCount, _ := count.Int64()

	if onlineCount > 0 {
		u.redisConn.Decr(USER_COUNT_KEY)
	} else {
		u.redisConn.Set(USER_COUNT_KEY, 0, 0)
	}
}

// 重設線上人數數量
func (u *UserCount) ResetCount() {
	u.redisConn.Set(USER_COUNT_KEY, 0, 0)
}
