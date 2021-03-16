package pipeline

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sunet/tq/pkg/message"
	"github.com/PaesslerAG/jsonpath"
	"regexp"
)

func MakeFilterPipeline(test func(message.Message) bool) Pipeline {
	return func(cs ...*message.MessageChannel) *message.MessageChannel {
		return message.FilterChannels(test, fmt.Sprintf("filter"), cs...)
	}
}

func MakeEQFilter(key string, value string) Pipeline {
	return func(cs ...*message.MessageChannel) *message.MessageChannel {
		return message.FilterChannels(func(o message.Message) bool {
			if val, ok := o[key]; !ok {
				return false
			} else {
				return val == value
			}
		}, fmt.Sprintf("eq %s == %s", key, value), cs...)
	}
}

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func regCompileAll(list []string) []*regexp.Regexp {
	res := make([]*regexp.Regexp, len(list))
	for i, v := range list {
		res[i], _ = regexp.Compile(v)
	}
	return res
}

func MakeMatchAnyFilter(path string, values ...string) Pipeline {
	return func(cs ...*message.MessageChannel) *message.MessageChannel {
		return message.FilterChannels(func(o message.Message) bool {
			pathValues, err := jsonpath.Get(path, o)
			if err != nil {
				logrus.Error(err)
			} else {
				for _, v := range pathValues.([]string) {
					if stringInSlice(v, values) {
						return true
					}
				}
			}
			return false
		}, fmt.Sprintf("any %s in %s", path, values), cs...)
	}
}
