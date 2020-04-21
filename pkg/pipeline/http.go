package pipeline

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/SUNET/tq/pkg/api"
	"github.com/SUNET/tq/pkg/message"
)

func MakeRESTPipeline(u string) Pipeline {
	url_parsed, err := url.Parse(u)
	if err != nil {
		Log.Panicf("Unable to parse url: %s", err.Error())
	}

	return func(...*message.MessageChannel) *message.MessageChannel {
		out := message.NewMessageChannel(fmt.Sprintf("http [%s]", u))
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
				}
				out.Send(o)
			}
			api.Handler.HandleFunc(url_parsed.Path, pubHandler).Methods("POST")
		}()
		return out
	}
}
