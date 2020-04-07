package pipeline

import (
	"github.com/SUNET/tq/pkg/message"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/pub"
	_ "go.nanomsg.org/mangos/v3/transport/all"
)

var newPubSink = func(url string) Pipeline {
	var err error
	var sock mangos.Socket
	var data []byte
	if sock, err = pub.NewSocket(); err != nil {
		Log.Panicf("can't create pub socket: %s", err.Error())
	}
	if err = sock.Listen(url); err != nil {
		Log.Panicf("can't listen to pub %s on socket: %s", url, err.Error())
	}

	return NewPipelineFromHandler(func(o message.Message) message.Message {
		data, err = message.FromJson(o)
		if err != nil {
			Log.Errorf("Error serializing json: %s", err)
		} else {
			sock.Send(data)
		}
		return o
	})
}

func init() {
	Register("pub", newPubSink)
}
