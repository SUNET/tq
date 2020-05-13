package pipeline

import (
	"fmt"

	"github.com/sunet/tq/pkg/message"
	"github.com/sunet/tq/pkg/utils"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/push"
	_ "go.nanomsg.org/mangos/v3/transport/all"
)

func MakePushPipeline(url string) Pipeline {
	sock, err := push.NewSocket()
	if err != nil {
		Log.Panicf("can't create push socket: %s", err.Error())
	}
	err = sock.Listen(url)
	if err != nil {
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
			return sendMessage(sock, o)
		}, fmt.Sprintf("push %s", url), cs...)
	}
}
