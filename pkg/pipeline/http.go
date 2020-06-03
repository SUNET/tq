package pipeline

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

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

func MakeRESTPipeline(path string) Pipeline {
	return func(...*message.MessageChannel) *message.MessageChannel {
		out := message.NewMessageChannel(fmt.Sprintf("http [%s]", path))
		go func() {
			pubHandler := func(w http.ResponseWriter, req *http.Request) {
				var err error
				var data []byte
				var o message.Message
				defer req.Body.Close()
				data, err = ioutil.ReadAll(req.Body)
				Log.Infof("got data: %s", data)
				io.WriteString(w, "ok\n")
				if o, err = message.ToJson(data); err != nil {
					Log.Errorf("Cannot parse json: %s", err.Error())
				} else {
					out.Send(o)
				}
			}
			api.Handler.HandleFunc(path, pubHandler).Methods("PUT", "POST")
		}()
		return out
	}
}
