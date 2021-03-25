package pipeline

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spy16/sabre"
	"github.com/sunet/tq/pkg/message"
)

var Log = logrus.New()

type Pipeline func(...*message.MessageChannel) *message.MessageChannel

func CallPipeline(p Pipeline, cs ...*message.MessageChannel) *message.MessageChannel {
	return p(cs...)
}

func Merge(cs ...*message.MessageChannel) *message.MessageChannel {
	return message.ProcessChannels(func(v message.Message) (message.Message, error) {
		return v, nil
	}, fmt.Sprintf("merge of %v", cs), cs...)
}

func ForkAndMerge(in *message.MessageChannel, pipelines ...*sabre.Fn) *message.MessageChannel {
	inputs := make([]*message.MessageChannel, len(pipelines))
	outputs := make([]*message.MessageChannel, len(pipelines))
	for i, p := range pipelines {
		inputs[i] = message.NewMessageChannel(fmt.Sprintf("input[%d]", i))
		outputs[i] = (*p)(inputs[i])
	}
	go func(inputs ...*message.MessageChannel) {
		message.ForkChannel(in, inputs...)
	}(inputs...)
	return Merge(outputs...)
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

func Run(cs ...*message.MessageChannel) {
	if len(cs) == 0 || cs == nil {
		cs = message.AllFinalChannels()
	}
	Log.Debugf("running %d final channels: %v", len(cs), cs)
	if len(cs) == 0 {
		select {} // for some reason the user wants us to block forever...
	} else if len(cs) == 1 {
		cs[0].Sink()
	} else {
		Merge(cs...).Sink()
	}
}
