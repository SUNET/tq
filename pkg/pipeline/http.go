package pipeline

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"

	"github.com/SUNET/tq/pkg/api"
	"github.com/SUNET/tq/pkg/message"
)

func MakePOSTPipeline(url string) Pipeline {
	return func(cs ...*message.MessageChannel) *message.MessageChannel {
		return message.ProcessChannels(func(o message.Message) (message.Message, error) {
			data, err := message.FromJson(o)
			if err != nil {
				Log.Errorf("unable to serialize to json: %s", err.Error())
				return nil, err
			}
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
			if err != nil {
				Log.Errorf("unable to create POST request: %s", err.Error())
				return nil, err
			}
			req.Header.Set("Content-Type", "application/json")
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				Log.Errorf("unable to send request: %s", err.Error())
				return nil, err
			}
			defer resp.Body.Close()

			Log.Debugf("POST to %s got response status: %d", url, resp.Status)
			ct := resp.Header.Get("content-type")
			mt, _, err := mime.ParseMediaType(ct)
			body, _ := ioutil.ReadAll(resp.Body)
			if mt == "application/json" {
				return message.ToJson(body)
			} else {
				return message.Jsonf("{\"status\": %d,\"body\": \"%s\"}", resp.Status, string(body))
			}
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
