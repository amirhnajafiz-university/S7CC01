package imagga

import (
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
func (i *Imagga) Process(url string) (string, error) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", "https://api.imagga.com/v2/tags?image_url="+url, nil)
	req.SetBasicAuth(i.Cfg.ApiKey, i.Cfg.ApiSecret)

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	return string(respBody), nil
}
