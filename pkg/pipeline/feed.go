package pipeline

import (
	"fmt"
	"github.com/sunet/tq/pkg/message"
)

func MakeFeedFetchPipeline(urlPath string) Pipeline {
	return func(cs ...*message.MessageChannel) *message.MessageChannel {
		return message.ConsumeChannels(func(o message.Message) (*message.MessageChannel, error) {
			return message.FetchFeedHandler(urlPath, o)
		}, fmt.Sprintf("fetch feed [%v]", urlPath), cs...)
	}
}