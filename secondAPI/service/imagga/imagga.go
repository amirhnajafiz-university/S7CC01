package imagga

import (
	"bytes"
	"fmt"
	"net/http"
)

// NewRequest
// sends one http request to Imagga website.
func NewRequest(cfg Config, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, cfg.URI, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("make http reqeust failed: %w", err)
	}

	client := &http.Client{}

	return client.Do(req)
}
