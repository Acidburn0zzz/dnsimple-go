package dnsimple

import (
	"errors"
	"fmt"
	"net/url"
)

// RecordsService handles communication with the record related
// methods of the DNSimple API.
//
// DNSimple API docs: http://developer.dnsimple.com/domains/records/
type RecordsService struct {
	client *DNSimpleClient
}

type Record struct {
	Id         int    `json:"id,omitempty"`
	DomainId   int    `json:"domain_id,omitempty"`
	Name       string `json:"name,omitempty"`
	Content    string `json:"content,omitempty"`
	TTL        int    `json:"ttl,omitempty"`
	Priority   int    `json:"prio,omitempty"`
	RecordType string `json:"record_type,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
}

type recordWrapper struct {
	Record Record `json:"record"`
}

// recordPath generates the resource path for given record that belongs to a domain.
func recordPath(domain interface{}, record interface{}) string {
	path := fmt.Sprintf("domains/%s/records", domainIdentifier(domain))

	if record != nil {
		path += fmt.Sprintf("/%d", record)
	}

	return path
}

// List the records for a domain that belongs to the authenticated user.
//
// DNSimple API docs: http://developer.dnsimple.com/domains/records/#list-records-for-a-domain
func (s *RecordsService) List(domain interface{}, name, recordType string) ([]Record, error) {
	reqStr := recordPath(domain, nil)
	v := url.Values{}

	if name != "" {
		v.Add("name", name)
	}

	if recordType != "" {
		v.Add("type", recordType)
	}

	reqStr += "?" + v.Encode()

	wrappedRecords := []recordWrapper{}

	if err := s.client.get(reqStr, &wrappedRecords); err != nil {
		return []Record{}, err
	}

	records := []Record{}
	for _, record := range wrappedRecords {
		records = append(records, record.Record)
	}

	return records, nil
}

// Create a new record for the specified domain.
//
// DNSimple API docs: http://developer.dnsimple.com/domains/records/#create-a-record
func (s *RecordsService) Create(domain interface{}, record Record) (Record, error) {
	// pre-validate the Record?
	wrappedRecord := recordWrapper{Record: record}
	returnedRecord := recordWrapper{}

	status, err := s.client.post(recordPath(domain, nil), wrappedRecord, &returnedRecord)
	if err != nil {
		return Record{}, err
	}

	if status == 400 {
		return Record{}, errors.New("Invalid Record")
	}

	return returnedRecord.Record, nil
}

// Get fetches a record.
//
// DNSimple API docs: http://developer.dnsimple.com/domains/records/#get-a-record
func (s *RecordsService) Get(domain interface{}, recordID int) (Record, error) {
	wrappedRecord := recordWrapper{}

	if err := s.client.get(recordPath(domain, recordID), &wrappedRecord); err != nil {
		return Record{}, err
	}

	return wrappedRecord.Record, nil
}

func (record *Record) Update(client *DNSimpleClient, recordAttributes Record) (Record, error) {
	// pre-validate the Record?
	// name, content, ttl, prio - only things allowed
	wrappedRecord := recordWrapper{Record: Record{
		Name:     recordAttributes.Name,
		Content:  recordAttributes.Content,
		TTL:      recordAttributes.TTL,
		Priority: recordAttributes.Priority}}

	returnedRecord := recordWrapper{}

	status, err := client.put(recordPath(record.DomainId, record.Id), wrappedRecord, &returnedRecord)
	if err != nil {
		return Record{}, err
	}

	if status == 400 {
		return Record{}, errors.New("Invalid Record")
	}

	return returnedRecord.Record, nil
}

func (record *Record) Delete(client *DNSimpleClient) error {
	response, err := client.sendRequest("DELETE", recordPath(record.DomainId, record.Id), nil, nil)
	if err != nil {
		return err
	}

	if response.StatusCode == 200 {
		return nil
	}

	return errors.New("Failed to delete domain")
}

func (record *Record) UpdateIP(client *DNSimpleClient, IP string) error {
	newRecord := Record{Content: IP, Name: record.Name}
	_, err := record.Update(client, newRecord)
	return err
}
