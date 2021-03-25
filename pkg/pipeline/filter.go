package pipeline

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"github.com/sunet/tq/pkg/message"
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

func stringInSlice(a string, list []*regexp.Regexp) bool {
    for _, b := range list {
        if b.Match([]byte(a)) {
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
	regs := regCompileAll(values)
	Log.Debugf("regs: %v", regs)
	return func(cs ...*message.MessageChannel) *message.MessageChannel {
		return message.FilterChannels(func(o message.Message) bool {
			Log.Debugf("filtering ----- %v", o)
			pathValues, err := jsonpath.Get(path, o)
			Log.Debugf("%v", pathValues)
			if err != nil {
				Log.Fatal(err)
				Log.Exit(1)
			} else {
				for _, v := range pathValues.([]string) {
					Log.Debug(v)
					if stringInSlice(v, regs) {
						Log.Debugf("%v", v)
						return true
					}
				}

			}
			return false
		}, fmt.Sprintf("any %s in %s", path, values), cs...)
	}
}
