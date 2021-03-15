package message

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gopkg.in/square/go-jose.v2/json"
)

var Log = logrus.New()

type Message map[string]interface{}
type MessageHandler func(Message) (Message, error)
type MessageSender func(*MessageChannel, Message)

func ToJson(data []byte) (Message, error) {
	var o Message
	err := json.Unmarshal(data, &o)
	return o, err
}

func FromJson(o Message) ([]byte, error) {
	return json.Marshal(o)
}

func Jsonf(format string, args ...interface{}) (Message, error) {
	jsonStr := fmt.Sprintf(format, args...)
	return ToJson([]byte(jsonStr))
}

func (msg *Message) String() (string, error) {
	s, err := FromJson(*msg)
	if err == nil {
		return string(s), nil
	} else {
		return "", err
	}
}
