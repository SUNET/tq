package message

import (
	"github.com/yalp/jsonpath"
)

func (msg *Message) GetPath(jp string) (interface{}, error) {
	return jsonpath.Read(msg, jp)
}
