package pipeline

import (
	"fmt"
	"go.nanomsg.org/mangos/v3/protocol"
	"time"

	"github.com/sunet/tq/pkg/message"
	"github.com/sunet/tq/pkg/utils"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/pull"
	"github.com/avast/retry-go"
	_ "go.nanomsg.org/mangos/v3/transport/all"
)

func MakePullPipeline(args ...string) Pipeline {
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
		sock, err = pull.NewSocket()
		if err != nil {
			Log.Warnf("can't create pull socket: %s", err.Error())
		}
		err = sock.Dial(url)
		if err != nil {
			Log.Warnf("can't dial %s on socket: %s", url, err.Error())
		}
		_, err = sock.GetOption(mangos.OptionTLSConfig)
		if err == nil {
			err = sock.SetOption(mangos.OptionTLSConfig, utils.GetTLSConfig())
			if err != nil {
				Log.Warnf("cannot set TLS op: %s", err.Error())
			}
		}

		if err != nil {
			Log.Warnf("Retrying in %v...", sleepTime)
			time.Sleep(sleepTime)
		}
		return err
	})

	if err != nil {
		Log.Panicf("Giving up!")
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
