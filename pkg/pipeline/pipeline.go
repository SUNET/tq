package pipeline

import (
	"github.com/SUNET/tq/pkg/message"
	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

type Pipeline func(...*message.MessageChannel) *message.MessageChannel

func Merge(cs ...*message.MessageChannel) *message.MessageChannel {
	return message.ProcessChannels(func(v message.Message) (message.Message, error) {
		return v, nil
	}, "merge", cs...)
}

func WaitForAll(cs ...*message.MessageChannel) {
	for _, c := range cs {
		Log.Print(c)
		c.Wait()
	}
}

func LogMessages(cs ...*message.MessageChannel) *message.MessageChannel {
	return message.ProcessChannels(func(v message.Message) (message.Message, error) {
		m, err := message.FromJson(v)
		if err != nil {
			Log.Errorf("Unable to serialize json: %s", err.Error())
			return nil, err
		} else {
			Log.Print(string(m))
			return v, nil
		}
	}, "log", cs...)
}

func RecvMessage(cs *message.MessageChannel) message.Message {
	return cs.Recv()
}

func SinkChannel(cs *message.MessageChannel) {
	cs.Sink()
}

func Run(cs ...*message.MessageChannel) {
	if len(cs) == 0 || cs == nil {
		cs = message.AllFinalChannels()
	}
	Log.Debugf("running %d final channels", len(cs))
	if len(cs) == 0 {
		SinkChannel(cs[0])
	} else {
		SinkChannel(Merge(cs...))
	}
}
