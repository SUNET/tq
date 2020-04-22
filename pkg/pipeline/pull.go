package pipeline

import (
	"fmt"

	"github.com/sunet/tq/pkg/message"
	"github.com/sunet/tq/pkg/utils"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/pull"
	_ "go.nanomsg.org/mangos/v3/transport/all"
)

func MakePullPipeline(url string) Pipeline {
	var err error
	var sock mangos.Socket

	if sock, err = pull.NewSocket(); err != nil {
		Log.Panicf("can't create pull socket: %s", err.Error())
	}
	if err = sock.Dial(url); err != nil {
		Log.Panicf("can't dial %s on socket: %s", url, err.Error())
	}
	_, err = sock.GetOption(mangos.OptionTLSConfig)
	if err == nil {
		err = sock.SetOption(mangos.OptionTLSConfig, utils.GetTLSConfig())
		if err != nil {
			Log.Panicf("cannot set TLS op: %s", err.Error())
		}
	}

	return func(...*message.MessageChannel) *message.MessageChannel {
		out := message.NewMessageChannel(fmt.Sprintf("pull [%s]", url))
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
