package pipeline

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/SUNET/tq/pkg/message"
)

var newRESTSource = func(u string) Pipeline {
	url_parsed, err := url.Parse(u)
	if err != nil {
		Log.Panicf("Unable to parse url: %s", err.Error())
	}

	return func(<-chan message.Message) <-chan message.Message {
		out := make(chan message.Message)
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
				out <- o
			}
			http.HandleFunc(url_parsed.Path, pubHandler)
			http.ListenAndServe(url_parsed.Host, nil)
		}()
		return out
	}
}

func init() {
	Register("rest", newRESTSource)
}
