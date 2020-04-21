package pipeline

import (
	"fmt"

	"github.com/SUNET/tq/pkg/message"
	"github.com/SUNET/tq/pkg/utils"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/sub"
	_ "go.nanomsg.org/mangos/v3/transport/all"
)

func MakeSubscribePipeline(args ...string) Pipeline {
	var err error
	var sock mangos.Socket

	if len(args) < 1 {
		Log.Panicf("sub requires at least one argument - the subscription URL")
	}

	url := args[0]
	var topic []byte
	if len(args) > 1 {
		topic = []byte(args[1])
	} else {
		topic = []byte("")
	}

	if sock, err = sub.NewSocket(); err != nil {
		Log.Panicf("can't create sub socket: %s", err.Error())
	}
	if err = sock.Dial(url); err != nil {
		Log.Panicf("can't dial %s on socket: %s", url, err.Error())
	}
	err = sock.SetOption(mangos.OptionSubscribe, topic)
	if err != nil {
		Log.Panicf("cannot subscribe (to '%s'): %s", string(topic), err.Error())
	}
	err = sock.SetOption(mangos.OptionTLSConfig, utils.GetTLSConfig())
	if err != nil {
		Log.Panicf("cannot set TLS op: %s", err.Error())
	}

	return func(...*message.MessageChannel) *message.MessageChannel {
		out := message.NewMessageChannel(fmt.Sprintf("sub [%s]", url))
		go func() {
			var err error
			var data []byte
			var o message.Message
			for {
				if data, err = sock.Recv(); err != nil {
					Log.Errorf("Cannot recv: %s", err.Error())
				}
				if o, err = message.ToJson(data); err != nil {
					Log.Errorf("Cannot parse json: %s", err.Error())
				}
				out.Send(o)
			}
		}()
		return out
	}
}
