package imagga

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// Imagga
// handles the image tagging.
type Imagga struct {
	Cfg Config
}

// Process
// sends one http request to Imagga website.
func (i *Imagga) Process(address string) (*Response, error) {
	// creating a new client
	client := &http.Client{}

	// creating a new get request
	req, _ := http.NewRequest(
		"GET",
		"https://api.imagga.com/v2/tags?image_url="+url.QueryEscape(address),
		nil,
	)
	// set the auth
	req.SetBasicAuth(i.Cfg.ApiKey, i.Cfg.ApiSecret)

	log.Printf("sending request to imagga:\n\t%s\n", req.URL)

	// do the http request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// return response
	respBody, _ := ioutil.ReadAll(resp.Body)

	var response Response

	if resp.StatusCode != http.StatusOK {
		log.Printf("imagga response: %s\n", resp.Status)
		log.Printf("\t%s\n", string(respBody))
	} else {
		if err := json.Unmarshal(respBody, &response); err != nil {
			return nil, err
		}
	}

	return &response, nil
}
