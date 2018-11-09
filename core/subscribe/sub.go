package subscribe

import (
	"github.com/go-redis/redis"
)

// WsPushChannel is channel name for redis subscribe
const WsPushChannel = "ws:push:channel"

var ch <-chan *redis.Message

func Sub(client *redis.Client) {
	pubsub := client.Subscribe(WsPushChannel)
	ch = pubsub.Channel()
}

func ReciveMessage(fn func([]byte)) {
	for msg := range ch {
		fn([]byte(msg.Payload))
	}
}
