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

type DNSimpleClient struct {
	ApiToken    string
	Email       string
	DomainToken string
	HttpClient  *http.Client
}

func NewClient(apiToken, email string) *DNSimpleClient {
	return &DNSimpleClient{ApiToken: apiToken, Email: email, HttpClient: &http.Client{}}
}

func (client *DNSimpleClient) get(url string, val interface{}) error {
	body, _, err := client.sendRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if err = json.Unmarshal([]byte(body), &val); err != nil {
		return err
	}

	return nil
}

func (client *DNSimpleClient) postOrPut(method, url string, payload, val interface{}) (int, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}

	body, status, err := client.sendRequest(method, url, strings.NewReader(string(jsonPayload)))
	if err != nil {
		return 0, err
	}

	if err = json.Unmarshal([]byte(body), &val); err != nil {
		return 0, err
	}

	return status, nil
}

func (client *DNSimpleClient) put(url string, payload, val interface{}) (int, error) {
	return client.postOrPut("PUT", url, payload, val)
}

func (client *DNSimpleClient) post(url string, payload, val interface{}) (int, error) {
	return client.postOrPut("POST", url, payload, val)
}

func (client *DNSimpleClient) makeRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	req.Header.Add("X-DNSimple-Token", fmt.Sprintf("%s:%s", client.Email, client.ApiToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	if err != nil {
		return nil, err
	}
	return req, nil
}

func (client *DNSimpleClient) sendRequest(method, url string, body io.Reader) (string, int, error) {
	req, err := client.makeRequest(method, url, body)
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

func domainIdentifier(value interface{}) string {
	switch value := value.(type) {
	case string:
		return value
	case int:
		return fmt.Sprintf("%d", value)
	case Domain:
		return fmt.Sprintf("%d", value.Id)
	case Record:
		return fmt.Sprintf("%d", value.DomainId)
	}
	return ""
}
