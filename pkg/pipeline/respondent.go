package pipeline

import (
	"fmt"

	"github.com/sunet/tq/pkg/message"
	"github.com/sunet/tq/pkg/utils"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/respondent"
	_ "go.nanomsg.org/mangos/v3/transport/all"
)

func MakeRespondentPipeline(url string) Pipeline {
	var err error
	var sock mangos.Socket

	if sock, err = respondent.NewSocket(); err != nil {
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

	return func(cs ...*message.MessageChannel) *message.MessageChannel {
		if len(cs) > 0 {
			return message.ProcessChannels(func(o message.Message) (message.Message, error) {
				return sendMessage(sock, o)
			}, fmt.Sprintf("respondent (out) %s", url), cs...)
		} else {
			out := message.NewMessageChannel(fmt.Sprintf("respondent (in) [%s]", url))
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
}
