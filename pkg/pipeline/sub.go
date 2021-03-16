package pipeline

import (
	"fmt"
	"github.com/sirupsen/logrus"

	"github.com/sunet/tq/pkg/message"
	"github.com/sunet/tq/pkg/utils"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/sub"
	_ "go.nanomsg.org/mangos/v3/transport/all"
)

func MakeSubscribePipeline(args ...string) Pipeline {
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

	sock, err := sub.NewSocket()
	if err != nil {
		Log.Panicf("can't create sub socket: %s", err.Error())
	}
	err = sock.Dial(url)
	if err != nil {
		Log.Panicf("can't dial %s on socket: %s", url, err.Error())
	}
	err = sock.SetOption(mangos.OptionSubscribe, topic)
	if err != nil {
		Log.Panicf("cannot subscribe (to '%s'): %s", string(topic), err.Error())
	}
	_, err = sock.GetOption(mangos.OptionTLSConfig)
	if err == nil {
		err = sock.SetOption(mangos.OptionTLSConfig, utils.GetTLSConfig())
		if err != nil {
			Log.Panicf("cannot set TLS op: %s", err.Error())
		}
	}

	return func(...*message.MessageChannel) *message.MessageChannel {
		out := message.NewMessageChannel(fmt.Sprintf("sub [%s]", url))
		go func() {
			for {
				o, err := recvMessage(sock)
				if Log.IsLevelEnabled(logrus.TraceLevel) {
					s, err := o.String()
					if err == nil {
						Log.Tracef("[SUB] %v <= %v", args, s)
					}
				}
				if err != nil {
					out.Send(o)
				}
			}
		}()
		return out
	}
}
