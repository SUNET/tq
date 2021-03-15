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
		ticker := time.NewTicker(d)
		go func() {
			defer out.Close()
			for {
				select {
				case <-out.Quit:
					Log.Debug("stopping timer")
					ticker.Stop()
					return
				case now := <-ticker.C:
					tBytes, err := now.MarshalJSON()
					if err != nil {
						Log.Errorf("Unable to create json: %s", err.Error())
					}
					o, err := message.Jsonf("{\"timestamp\": %s}", string(tBytes))
					if err != nil {
						Log.Errorf("Unable to create json: %s", err.Error())
					}
					out.Send(o)
				}
			}
		}()
		return out
	}
}
