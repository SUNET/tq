package pipeline

import (
	"github.com/SUNET/tq/pkg/message"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/sub"
	_ "go.nanomsg.org/mangos/v3/transport/all"
)

var newSubSource = func(url string) Pipeline {
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

	return func(<-chan message.Message) <-chan message.Message {
		out := make(chan message.Message)
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
				out <- o
			}
			close(out)
		}()
		return out
	}
}

func init() {
	Register("sub", newSubSource)
}
