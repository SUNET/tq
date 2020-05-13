package pipeline

import (
	"github.com/sunet/tq/pkg/message"
	"go.nanomsg.org/mangos/v3"
	_ "go.nanomsg.org/mangos/v3/transport/all"
)

func recvMessage(sock mangos.Socket) (message.Message, error) {
	data, err := sock.Recv()
	if err != nil {
		Log.Errorf("Cannot recv: %s", err.Error())
		return nil, err
	}
	o, err := message.ToJson(data)
	if err != nil {
		Log.Errorf("Cannot parse json: %s", err.Error())
		return nil, err
	}
	return o, nil
}

func sendMessage(sock mangos.Socket, o message.Message) (message.Message, error) {
	data, err := message.FromJson(o)
	if err != nil {
		Log.Errorf("Error serializing json: %s", err)
		return nil, err
	} else {
		sock.Send(data)
		return o, nil
	}
}
