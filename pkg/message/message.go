package message

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/square/go-jose.v2/json"
)

var Log = logrus.New()

type Message map[string]interface{}

func ToJson(data []byte) (Message, error) {
	var o Message
	err := json.Unmarshal(data, &o)
	return o, err
}

func FromJson(o Message) ([]byte, error) {
	return json.Marshal(o)
}
