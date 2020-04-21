package pipeline

import (
	"fmt"

	"github.com/SUNET/tq/pkg/message"
	"github.com/SUNET/tq/pkg/utils"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/push"
	_ "go.nanomsg.org/mangos/v3/transport/all"
)

func MakePushPipeline(url string) Pipeline {
	var err error
	var sock mangos.Socket
	var data []byte
	if sock, err = push.NewSocket(); err != nil {
		Log.Panicf("can't create push socket: %s", err.Error())
	}
	if err = sock.Listen(url); err != nil {
		Log.Panicf("can't listen to push %s on socket: %s", url, err.Error())
	}
	_, err = sock.GetOption(mangos.OptionTLSConfig)
	if err == nil {
		err = sock.SetOption(mangos.OptionTLSConfig, utils.GetTLSConfig())
		if err != nil {
			Log.Panicf("cannot set TLS op: %s", err.Error())
		}
	}

	return func(cs ...*message.MessageChannel) *message.MessageChannel {
		return message.ProcessChannels(func(o message.Message) (message.Message, error) {
			data, err = message.FromJson(o)
			if err != nil {
				Log.Errorf("Error serializing json: %s", err)
				return nil, err
			} else {
				sock.Send(data)
				return o, nil
			}
		}, fmt.Sprintf("push %s", url), cs...)
	}
}
