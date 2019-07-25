package message

import "encoding/json"

var messageMap map[int64]string

func init() {
	messageMap = map[int64]string{
		404001: "非集群模式不可操作",
	}
}

type res struct {
	Code    int64       `json:"code"`
	Message interface{} `json:"message"`
}

// Res ...
func Res(code int64, message interface{}) string {
	var msg interface{}
	switch message.(type) {
	case string:
		if messageMap[code] != "" {
			msg = messageMap[code]
		} else {
			msg = message
		}
	default:
		msg = message
	}
	bytes, _ := json.Marshal(&res{
		Code:    code,
		Message: msg,
	})
	return string(bytes)
}
