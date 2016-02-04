package dnsimple

import (
	"fmt"
	"net/http"
)

type Record struct {
	ID        int    `json:"id,omitempty"`
	ZoneID    string `json:"zone_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Content   string `json:"content,omitempty"`
	TTL       int    `json:"ttl,omitempty"`
	Priority  int    `json:"priority,omitempty"`
	Type      string `json:"type,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type recordsWrapper struct {
	Records []Record `json:"data"`
}
type recordWrapper struct {
	Record *Record `json:"data"`
}

// recordPath generates the resource path for given record that belongs to a domain.
func recordPath(accountID string, domain interface{}, record interface{}) string {
	path := fmt.Sprintf("/%v/zones/%v/records", accountID, domainIDentifier(domain))

	if record != nil {
		path += fmt.Sprintf("/%v", record)
	}

	return path
}

// List the zone records.
//
// See https://developer.dnsimple.com/v2/zones/#list
func (s *ZonesService) ListRecords(accountID string, domain interface{}) ([]Record, *http.Response, error) {
	path := recordPath(accountID, domain, nil)
	data := recordsWrapper{}

	res, err := s.client.get(path, &data)
	if err != nil {
		return []Record{}, res, err
	}

	return data.Records, res, nil
}

// CreateRecord creates a zone record.
//
// See https://developer.dnsimple.com/v2/zones/#create
func (s *ZonesService) CreateRecord(accountID string, domain interface{}, recordAttributes Record) (*Record, *http.Response, error) {
	path := recordPath(accountID, domain, nil)
	data := recordWrapper{}

	res, err := s.client.post(path, recordAttributes, &data)
	if err != nil {
		return &Record{}, res, err
	}

	return data.Record, res, nil
}

// GetRecord gets the zone record.
//
// See https://developer.dnsimple.com/v2/zones/#get
func (s *ZonesService) GetRecord(accountID string, domain interface{}, recordID int) (*Record, *http.Response, error) {
	path := recordPath(accountID, domain, recordID)
	data := recordWrapper{}

	res, err := s.client.get(path, &data)
	if err != nil {
		return &Record{}, res, err
	}

	return data.Record, res, nil
}

// UpdateRecord updates a zone record.
//
// See https://developer.dnsimple.com/v2/zones/#update
func (s *ZonesService) UpdateRecord(accountID string, domain interface{}, recordID int, recordAttributes Record) (*Record, *http.Response, error) {
	path := recordPath(accountID, domain, recordID)
	data := recordWrapper{}

	res, err := s.client.patch(path, recordAttributes, &data)
	if err != nil {
		return &Record{}, res, err
	}

	return data.Record, res, nil
}

// DeleteRecord deletes a zone record.
//
// See https://developer.dnsimple.com/v2/zones/#delete
func (s *ZonesService) DeleteRecord(accountID string, domain interface{}, recordID int) (*http.Response, error) {
	path := recordPath(accountID, domain, recordID)

	return s.client.delete(path, nil, nil)
}
