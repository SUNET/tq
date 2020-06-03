package message

import (
	"bytes"
	"io/ioutil"
	"mime"
	"net/http"
)

func PostHandler(url string, o Message) (Message, error) {
	data, err := FromJson(o)
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

	Log.Debugf("POST to %s got response status: %v", url, resp.Status)
	ct := resp.Header.Get("content-type")
	mt, _, err := mime.ParseMediaType(ct)
	body, _ := ioutil.ReadAll(resp.Body)
	if mt == "application/json" {
		return ToJson(body)
	} else {
		return Jsonf("{\"status\": %v,\"body\": \"%s\"}", resp.Status, string(body))
	}
}
