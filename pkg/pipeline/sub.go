package pipeline

import (
	"fmt"

	"github.com/SUNET/tq/pkg/message"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/sub"
	_ "go.nanomsg.org/mangos/v3/transport/all"
)

func MakeSubscribePipeline(url string) Pipeline {
	var err error
	var sock mangos.Socket
	if sock, err = sub.NewSocket(); err != nil {
		Log.Panicf("can't create sub socket: %s", err.Error())
	}
	if err = sock.Dial(url); err != nil {
		Log.Panicf("can't dial %s on socket: %s", url, err.Error())
	}
	err = sock.SetOption(mangos.OptionSubscribe, []byte(""))
	if err != nil {
		Log.Panicf("cannot subscribe: %s", err.Error())
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
