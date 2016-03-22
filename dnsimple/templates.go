package dnsimple

import (
	"fmt"
)

// TemplatesService handles communication with the template related
// methods of the DNSimple API.
//
// See https://developer.dnsimple.com/v2/templates/
type TemplatesService struct {
	client *Client
}

// Template represents a Template in DNSimple.
type Template struct {
	ID          int    `json:"id,omitempty"`
	AccountID   int    `json:"account_id,omitempty"`
	Name        string `json:"name,omitempty"`
	ShortName   string `json:"short_name,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

func templatePath(accountID string, templateID int) string {
	if templateID != 0 {
		return fmt.Sprintf("/%v/templates/%v", accountID, templateID)
	}

	return fmt.Sprintf("/%v/templates", accountID)
}

// TemplatesResponse represents a response from an API method that returns a collection of Template struct.
type TemplatesResponse struct {
	Response
	Data []Template `json:"data"`
}

// ListTemplates list the templates for an account.
//
// See https://developer.dnsimple.com/v2/templates/#list
func (s *TemplatesService) ListTemplates(accountID string, options *ListOptions) (*TemplatesResponse, error) {
	path := versioned(templatePath(accountID, 0))
	templatesResponse := &TemplatesResponse{}

	path, err := addURLQueryOptions(path, options)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.get(path, templatesResponse)
	if err != nil {
		return templatesResponse, err
	}

	templatesResponse.HttpResponse = resp
	return templatesResponse, nil
}
