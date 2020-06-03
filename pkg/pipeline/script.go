package pipeline

import (
	"fmt"
	"strings"

	"github.com/sunet/tq/pkg/message"
)

func MakeScriptPipeline(cmdline ...string) Pipeline {
	return func(cs ...*message.MessageChannel) *message.MessageChannel {
		return message.ProcessChannels(func(o message.Message) (message.Message, error) {
			return message.ScriptHandler(cmdline, o)
		}, fmt.Sprintf("script %s", strings.Join(cmdline, " ")), cs...)
	}
}
