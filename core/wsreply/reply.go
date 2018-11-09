package wsreply

import "time"

//
const (
	ResultSuccess = "success"
	ResultFailure = "failure"
)

// Reply struct for websocket reply
type Reply struct {
	Action    string      `json:"action"`
	Result    string      `json:"result"`
	Timestamp int64       `json:"timestamp"`
	Code      int         `json:"code"`
	Data      interface{} `json:"data"`
}

// NewSuccessReply create Reply instance with success result
func NewSuccessReply(action string, data interface{}) Reply {
	return Reply{
		Action:    action,
		Result:    ResultSuccess,
		Timestamp: time.Now().Unix(),
		Data:      data,
	}
}

// NewFailureReply create Reply instance with failure result
func NewFailureReply(action string, code int) Reply {
	return Reply{
		Action:    action,
		Result:    ResultFailure,
		Timestamp: time.Now().Unix(),
		Code:      code,
	}
}
