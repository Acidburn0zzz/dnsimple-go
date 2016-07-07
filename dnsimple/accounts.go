package dnsimple

import (
	"fmt"
)

type AccountsService struct {
	client *Client
}

// Account represents a DNSimple account.
type Account struct {
	ID             int    `json:"id,omitempty"`
	Email          string `json:"email,omitempty"`
	PlanIdentifier string `json:"plan_identifier,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
}

func accountsPath() string {
	return fmt.Sprintf("/accounts")
}

// AccountResponse represents a response from an API method that returns an Account struct.
type AccountResponse struct {
	Response
	Data *Account `json:"data"`
}

// AccountsResponse represents a response from an API method that returns a collection of Account struct.
type AccountsResponse struct {
	Response
	Data []Account `json:"data"`
}

// ListAccounts list the accounts for an user.
//
// See https://developer.dnsimple.com/v2/accounts/#list
func (s *AccountsService) ListAccounts(options *ListOptions) (*AccountsResponse, error) {
	path := versioned(accountsPath())
	accountsResponse := &AccountsResponse{}

	path, err := addURLQueryOptions(path, options)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.get(path, accountsResponse)
	if err != nil {
		return accountsResponse, err
	}

	accountsResponse.HttpResponse = resp
	return accountsResponse, nil
}
