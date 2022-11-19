package imagga

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
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

// readFile
// reads a file to submit over imagga.
func readFile(filePath string) (*bytes.Buffer, *multipart.Writer, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("image", file.Name())
	if err != nil {
		return nil, nil, err
	}

	n, err := part.Write(fileContents)
	if err != nil {
		return nil, nil, err
	}

	log.Printf("read file: %d bytes\n", n)

	return body, writer, nil
}

// Upload
// files to imagga service and get url.
func (i *Imagga) Upload(filePath string) (string, error) {
	// reading the file
	body, writer, err := readFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// creating a new http client
	client := &http.Client{}

	// creating a new request
	req, _ := http.NewRequest(
		"POST",
		"https://api.imagga.com/v2/uploads",
		body,
	)

	// set the headers
	req.SetBasicAuth(i.Cfg.ApiKey, i.Cfg.ApiSecret)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// make http call
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	// return response
	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed: %s", string(respBody))
	}

	return string(respBody), nil
}
