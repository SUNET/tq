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
	sock, err := pull.NewSocket()
	if err != nil {
		Log.Panicf("can't create pull socket: %s", err.Error())
	}
	err = sock.Dial(url)
	if err != nil {
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
			for {
				o, err := recvMessage(sock)
				if err != nil {
					out.Send(o)
				}
			}
		}()
		return out
	}
}
