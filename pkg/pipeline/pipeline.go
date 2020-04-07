package pipeline

import (
	"github.com/SUNET/tq/pkg/message"
	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

type Pipeline func(src <-chan message.Message) <-chan message.Message
type PipelineFactory func(s string) Pipeline
type Handler func(o message.Message) message.Message

var PipelineFactories = map[string]PipelineFactory{}

func Register(name string, pf PipelineFactory) {
	PipelineFactories[name] = pf
}

func IsPipelineFactory(name string) bool {
	if _, ok := PipelineFactories[name]; ok {
		return true
	} else {
		return false
	}
}

func NewPipelineFromHandler(h Handler) Pipeline {
	return func(in <-chan message.Message) <-chan message.Message {
		out := make(chan message.Message)
		go func() {
			o := h(<-out)
			if o != nil {
				out <- o
			}
			close(out)
		}()
		return out
	}
}
