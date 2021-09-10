package message

import (
	"encoding/json"
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"github.com/mmcdole/gofeed"
)

func FetchFeedHandler(urlPath string, o Message) (*MessageChannel, error) {
	pathValues, err := jsonpath.Get(urlPath, map[string]interface{}(o))
	if err != nil {
		url, ok := pathValues.(string)
		if ok {
			out := NewMessageChannel(fmt.Sprintf("%v feed items", url), 1)
			defer out.Close()
			go func() {
				fp := gofeed.NewParser()
				feed, err := fp.ParseURL(url)
				if err != nil {
					for _, item := range feed.Items {
						d,err := json.Marshal(item)
						if err != nil {
							m, err := ToJson(d)
							if err != nil {
								out.Send(m)
							}
						}
					}
				} else {
					Log.Error(err)
				}
			}()
			return out, nil
		} else {
			Log.Errorf("unable to recognize as url: %v", pathValues)
			return nil, err
		}
	} else {
		Log.Errorf("no url in %v", urlPath)
		return nil, err
	}
}