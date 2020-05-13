package pipeline

import (
	"fmt"
	"time"

	"github.com/sunet/tq/pkg/message"
	"github.com/sunet/tq/pkg/utils"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/surveyor"
	_ "go.nanomsg.org/mangos/v3/transport/all"
)

func MakeSurveyorPipeline(args ...string) Pipeline {
	var err error
	var sock mangos.Socket
	var url string
	var duration string = "5m"

	url = args[0]
	if len(args) > 1 {
		duration = args[1]
	}

	d, err := time.ParseDuration(duration)
	if err != nil {
		Log.Panicf("unable to parse duration '%s': %s", duration, err.Error())
	}

	if sock, err = surveyor.NewSocket(); err != nil {
		Log.Panicf("can't create surveyor socket: %s", err.Error())
	}
	if err = sock.Listen(url); err != nil {
		Log.Panicf("can't dial %s on socket: %s", url, err.Error())
	}
	err = sock.SetOption(mangos.OptionSurveyTime, d)
	if err != nil {
		Log.Panicf("can't set ttl on surveyor socket: %s", err.Error())
	}
	_, err = sock.GetOption(mangos.OptionTLSConfig)
	if err == nil {
		err = sock.SetOption(mangos.OptionTLSConfig, utils.GetTLSConfig())
		if err != nil {
			Log.Panicf("cannot set TLS op: %s", err.Error())
		}
	}

	return func(cs ...*message.MessageChannel) *message.MessageChannel {
		return message.ConsumeChannels(func(o message.Message) (*message.MessageChannel, error) {
			_, err := sendMessage(sock, o)
			if err != nil {
				return message.ProcessChannels(func(message.Message) (message.Message, error) {
					return recvMessage(sock)
				}, fmt.Sprintf("surveyor (in) %s", url)), nil
			} else {
				return nil, err
			}
		}, fmt.Sprintf("surveyor (out) %s", url), cs...)
	}
}
