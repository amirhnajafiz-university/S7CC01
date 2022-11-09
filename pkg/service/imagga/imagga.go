package imagga

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Imagga
// handles the image tagging.
type Imagga struct {
	Cfg Config
}

// Process
// sends one http request to Imagga website.
func (i *Imagga) Process(url string) (*Response, error) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", "https://api.imagga.com/v2/tags?image_url="+url, nil)
	req.SetBasicAuth(i.Cfg.ApiKey, i.Cfg.ApiSecret)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	var response Response

	if err := json.Unmarshal(respBody, &response); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	return &response, nil
}
