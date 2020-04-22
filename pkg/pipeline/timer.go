package pipeline

import (
	"fmt"
	"time"

	"github.com/sunet/tq/pkg/message"
)

func MakeTimer(duration string) Pipeline {

	d, err := time.ParseDuration(duration)
	if err != nil {
		Log.Panicf("Unable to parse duration: %s", err.Error())
	}

	return func(...*message.MessageChannel) *message.MessageChannel {
		out := message.NewMessageChannel(fmt.Sprintf("timer [%s]", duration))
		go func() {
			for now := range time.Tick(d) {
				t_bytes, err := now.MarshalJSON()
				if err != nil {
					Log.Errorf("Unable to create json: %s", err.Error())
				}
				o, err := message.Jsonf("{\"timestamp\": %s}", string(t_bytes))
				if err != nil {
					Log.Errorf("Unable to create json: %s", err.Error())
				}
				out.Send(o)
			}
		}()
		return out
	}
}
