package pipeline

import (
	"fmt"
	"go.nanomsg.org/mangos/v3/protocol"
	"time"

	"github.com/sunet/tq/pkg/message"
	"github.com/sunet/tq/pkg/utils"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/respondent"
	"github.com/avast/retry-go"
	_ "go.nanomsg.org/mangos/v3/transport/all"
)

func MakeRespondentPipeline(args ...string) Pipeline {
	url := args[0]
	delay := "3s"
	if len(args) > 1 {
		delay = args[1]
	}

	sleepTime, err := time.ParseDuration(delay)
	if err != nil {
		Log.Panicf("Unable to parse duration %s", delay)
	}

	var sock protocol.Socket

	err = retry.Do(func() error {

		if err != nil {
			Log.Warnf("Retrying in %v...", sleepTime)
			time.Sleep(sleepTime)
		}

		if sock, err = respondent.NewSocket(); err != nil {
			Log.Warnf("can't create pull socket: %s", err.Error())
			return err
		}
		if err = sock.Dial(url); err != nil {
			Log.Warnf("can't dial %s on socket: %s", url, err.Error())
			return err
		}
		_, err = sock.GetOption(mangos.OptionTLSConfig)
		if err == nil {
			err = sock.SetOption(mangos.OptionTLSConfig, utils.GetTLSConfig())
			if err != nil {
				Log.Warnf("cannot set TLS op: %s", err.Error())
				return err
			}
		}

		return nil
	})

	if err != nil {
		Log.Panicf("Giving up!")
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
