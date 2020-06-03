package pipeline

import (
	"github.com/sunet/tq/pkg/message"
)

func MakeKazaamPipeline(spec string) Pipeline {

	kz := message.NewKazaam(spec)

	return func(cs ...*message.MessageChannel) *message.MessageChannel {
		return message.ProcessChannels(func(o message.Message) (message.Message, error) {
			return message.KazaamHandler(kz, o)
		}, "kazaam", cs...)
	}
}
