package pipeline

import (
	"github.com/SUNET/tq/pkg/message"
	"github.com/qntfy/kazaam"
)

func MakeKazaamPipeline(spec string) Pipeline {

	kz, err := kazaam.NewKazaam(spec)
	if err != nil {
		Log.Fatalf("Unable to create kazaam from %s: %s", spec, err.Error())
	}

	return func(cs ...*message.MessageChannel) *message.MessageChannel {
		return message.ProcessChannels(func(o message.Message) (message.Message, error) {
			j, err := message.FromJson(o)
			if err != nil {
				return nil, err
			}

			j, err = kz.Transform(j)
			if err != nil {
				return nil, err
			}

			o, err = message.ToJson(j)
			if err != nil {
				return nil, err
			}

			return o, nil

		}, cs...)
	}
}
