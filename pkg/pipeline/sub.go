package pipeline

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go.nanomsg.org/mangos/v3/protocol"
	"time"

	"github.com/sunet/tq/pkg/message"
	"github.com/sunet/tq/pkg/utils"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/sub"
	"github.com/avast/retry-go"
	_ "go.nanomsg.org/mangos/v3/transport/all"
)

func MakeSubscribePipeline(args ...string) Pipeline {
	if len(args) < 1 {
		Log.Panicf("sub requires at least one argument - the subscription URL")
	}

	url := args[0]
	var topic []byte
	delay := "3s"
	if len(args) > 1 {
		topic = []byte(args[1])
	} else {
		topic = []byte("")
	}
	if len(args) > 2 {
		delay = args[2]
	}

	sleepTime, err := time.ParseDuration(delay)
	if err != nil {
		Log.Panicf("Unable to parse duration %s", delay)
	}

	var sock protocol.Socket

	err = retry.Do(func() error {

		sock, err := sub.NewSocket()
		if err != nil {
			Log.Warnf("can't create sub socket: %s", err.Error())
		}
		err = sock.Dial(url)
		if err != nil {
			Log.Warnf("can't dial %s on socket: %s", url, err.Error())
		}
		err = sock.SetOption(mangos.OptionSubscribe, topic)
		if err != nil {
			Log.Warnf("cannot subscribe (to '%s'): %s", string(topic), err.Error())
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
