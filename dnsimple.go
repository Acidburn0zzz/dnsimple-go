// Package dnsimple implements a client for the DNSimple API.
//
// In order to use this package you will need a DNSimple account and your API Token.
package dnsimple

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	libraryVersion = "0.1"
	defaultBaseURL = "https://api.dnsimple.com/"
	userAgent      = "go-dnsimple/" + libraryVersion

	apiVersion = "v1"
)

type DNSimpleClient struct {
	// HTTP client used to communicate with the API.
	HttpClient *http.Client

	// API Token for the DNSimple account you want to use.
	ApiToken string

	// Email associated with the provided DNSimple API Token.
	Email string

	// Domain Token to be used for authentication
	// as an alternative to the DNSimple API Token for some domain-scoped operations.
	DomainToken string

	// Base URL for API requests.
	// Defaults to the public DNSimple API, but can be set to a different endpoint (e.g. the sandbox).
	// BaseURL should always be specified with a trailing slash.
	BaseURL string

	// User agent used when communicating with the DNSimple API.
	UserAgent string
}

func NewClient(apiToken, email string) *DNSimpleClient {
	return &DNSimpleClient{ApiToken: apiToken, Email: email, HttpClient: &http.Client{}, BaseURL: defaultBaseURL, UserAgent: userAgent}
}

func (client *DNSimpleClient) get(path string, val interface{}) error {
	body, _, err := client.sendRequest("GET", path, nil)
	if err != nil {
		return err
	}

	if err = json.Unmarshal([]byte(body), &val); err != nil {
		return err
	}

	return nil
}

func (client *DNSimpleClient) postOrPut(method, path string, payload, val interface{}) (int, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}

	body, status, err := client.sendRequest(method, path, strings.NewReader(string(jsonPayload)))
	if err != nil {
		return 0, err
	}

	if err = json.Unmarshal([]byte(body), &val); err != nil {
		return 0, err
	}

	return status, nil
}

func (client *DNSimpleClient) put(path string, payload, val interface{}) (int, error) {
	return client.postOrPut("PUT", path, payload, val)
}

func (client *DNSimpleClient) post(path string, payload, val interface{}) (int, error) {
	return client.postOrPut("POST", path, payload, val)
}

func (client *DNSimpleClient) makeRequest(method, path string, body io.Reader) (*http.Request, error) {
	url := client.BaseURL + fmt.Sprintf("%s/%s", apiVersion, path)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", client.UserAgent)
	req.Header.Add("X-DNSimple-Token", fmt.Sprintf("%s:%s", client.Email, client.ApiToken))

	return req, nil
}

func (client *DNSimpleClient) sendRequest(method, path string, body io.Reader) (string, int, error) {
	req, err := client.makeRequest(method, path, body)
	if err != nil {
		return "", 0, err
	}

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	return string(responseBytes), resp.StatusCode, nil
}
