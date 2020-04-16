package pipeline

import (
	"github.com/SUNET/tq/pkg/message"
	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

type Pipeline func(...*message.MessageChannel) *message.MessageChannel

func Merge(cs ...*message.MessageChannel) *message.MessageChannel {
	return message.ProcessChannels(func(v message.Message) message.Message {
		return v
	}, cs...)
}

func WaitForAll(cs ...*message.MessageChannel) {
	for _, c := range cs {
		Log.Print(c)
		c.Wait()
	}
}

func LogMessages(cs ...*message.MessageChannel) *message.MessageChannel {
	return message.ProcessChannels(func(v message.Message) message.Message {
		m, err := message.FromJson(v)
		if err != nil {
			Log.Errorf("Unable to serialize json: %s", err.Error())
		} else {
			Log.Print(string(m))
		}
		return v
	}, cs...)
}

func RecvMessage(cs *message.MessageChannel) message.Message {
	return cs.Recv()
}
