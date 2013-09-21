package dnsimple

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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

func (client *DNSimpleClient) sendRequest(method, url string, body io.Reader) (string, error) {
	req, err := client.makeRequest(method, url, body)
	if err != nil {
		return "", err
	}

	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(responseBody), nil
}
