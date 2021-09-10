package pipeline

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/kjk/betterguid"
	"github.com/sunet/tq/pkg/api"
	"github.com/sunet/tq/pkg/message"
)

func MakePOSTPipeline(url string) Pipeline {
	return func(cs ...*message.MessageChannel) *message.MessageChannel {
		return message.ProcessChannels(func(o message.Message) (message.Message, error) {
			return message.PostHandler(url, o)
		}, fmt.Sprintf("post %s", url), cs...)
	}
}

func MakeHTTPSourcePipeline(path string) Pipeline {
	return MakeHTTPipeline(path, func(out *message.MessageChannel, o message.Message) (message.Message, error) {
		out.Send(o)
		return nil, nil
	})
}

func MakeHTTPSourceResponsePipeline(path string, sub_url string) Pipeline {
	resp := MakeSubscribePipeline(sub_url)
	return MakeHTTPipeline(path, func(out *message.MessageChannel, o message.Message) (message.Message, error) {
		o["RequestID"] = betterguid.New()
		out.Send(o)
		for ro := range resp(out).C {
			if ro["RequestID"] == o["RequestID"] {
				return ro, nil
			}
		}
		return nil, nil
	})
}

func MakeHTTPipeline(path string, sender func(*message.MessageChannel, message.Message) (message.Message, error)) Pipeline {
	return func(...*message.MessageChannel) *message.MessageChannel {
		out := message.NewMessageChannel(fmt.Sprintf("httpd [%s]", path))
		go func() {
			pubHandler := func(w http.ResponseWriter, req *http.Request) {
				var err error
				var data []byte
				var o message.Message
				defer req.Body.Close()
				data, err = ioutil.ReadAll(req.Body)
				Log.Infof("got data: %s", data)

				if o, err = message.ToJson(data); err != nil {
					Log.Errorf("Cannot parse json: %s", err.Error())
				} else {
					if ro, err := sender(out, o); err != nil {
						Log.Errorf("Unable to send: %s", err.Error())
					} else {
						if ro != nil {
							s, _ := ro.String()
							io.WriteString(w, s)
						} else {
							io.WriteString(w, "Message sent\n")
						}
					}
				}
			}
			api.Handler.HandleFunc(path, pubHandler).Methods("PUT", "POST")
		}()
		return out
	}
}
