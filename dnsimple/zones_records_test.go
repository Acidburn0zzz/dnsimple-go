package dnsimple

import (
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestZonesService_zoneRecordPath(t *testing.T) {
	if want, got := "/1010/zones/example.com/records", zoneRecordPath("1010", "example.com", 0); want != got {
		t.Errorf("contactPath(%v,  ) = %v, want %v", "1010", got, want)
	}

	if want, got := "/1010/zones/example.com/records/1", zoneRecordPath("1010", "example.com", 1); want != got {
		t.Errorf("contactPath(%v, 1) = %v, want %v", "1010", got, want)
	}
}

func TestZonesService_ListRecords(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/listZoneRecords/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	recordsResponse, err := client.Zones.ListRecords("1010", "example.com", nil)
	if err != nil {
		t.Fatalf("Zones.ListRecords() returned error: %v", err)
	}

	if want, got := (&Pagination{CurrentPage: 1, PerPage: 30, TotalPages: 1, TotalEntries: 5}), recordsResponse.Pagination; !reflect.DeepEqual(want, got) {
		t.Errorf("Zones.ListRecords() pagination expected to be %v, got %v", want, got)
	}

	records := recordsResponse.Data
	if want, got := 5, len(records); want != got {
		t.Errorf("Zones.ListRecords() expected to return %v contacts, got %v", want, got)
	}

	if want, got := 64779, records[0].ID; want != got {
		t.Fatalf("Zones.ListRecords() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "", records[0].Name; want != got {
		t.Fatalf("Zones.ListRecords() returned Name expected to be `%v`, got `%v`", want, got)
	}
}

func TestZonesService_ListRecords_WithOptions(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/listZoneRecords/success.http")

		testQuery(t, r, url.Values{"page": []string{"2"}, "per_page": []string{"20"}})

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	_, err := client.Zones.ListRecords("1010", "example.com", &ListOptions{Page: 2, PerPage: 20})
	if err != nil {
		t.Fatalf("Zones.ListRecords() returned error: %v", err)
	}
}

func TestZonesService_CreateRecord(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/createZoneRecord/created.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"name": "foo", "content": "192.168.0.10", "type": "A"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	recordValues := ZoneRecord{Name: "foo", Content: "192.168.0.10", Type: "A"}

	recordResponse, err := client.Zones.CreateRecord(accountID, "example.com", recordValues)
	if err != nil {
		t.Fatalf("Zones.CreateRecord() returned error: %v", err)
	}

	record := recordResponse.Data
	if want, got := 64784, record.ID; want != got {
		t.Fatalf("Zones.CreateRecord() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "www", record.Name; want != got {
		t.Fatalf("Zones.CreateRecord() returned Name expected to be `%v`, got `%v`", want, got)
	}
}

func TestZonesService_CreateRecord_BlankName(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/createZoneRecord/created_apex.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"name": "", "content": "192.168.0.10", "type": "A"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	recordValues := ZoneRecord{Name: "", Content: "192.168.0.10", Type: "A"}

	recordResponse, err := client.Zones.CreateRecord(accountID, "example.com", recordValues)
	if err != nil {
		t.Fatalf("Zones.CreateRecord() returned error: %v", err)
	}

	record := recordResponse.Data
	if want, got := 64784, record.ID; want != got {
		t.Fatalf("Zones.CreateRecord() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "", record.Name; want != got {
		t.Fatalf("Zones.CreateRecord() returned Name expected to be `%v`, got `%v`", want, got)
	}
}

func TestZonesService_GetRecord(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records/1539", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/getZoneRecord/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"

	recordResponse, err := client.Zones.GetRecord(accountID, "example.com", 1539)
	if err != nil {
		t.Fatalf("Zones.GetRecord() returned error: %v", err)
	}

	record := recordResponse.Data
	wantSingle := &ZoneRecord{
		ID:           64784,
		ZoneID:       "example.com",
		ParentID:     0,
		Type:         "A",
		Name:         "www",
		Content:      "127.0.0.1",
		TTL:          600,
		Priority:     0,
		SystemRecord: false,
		CreatedAt:    "2016-01-07T17:45:13.653Z",
		UpdatedAt:    "2016-01-07T17:45:13.653Z"}

	if !reflect.DeepEqual(record, wantSingle) {
		t.Fatalf("Zones.GetRecord() returned %+v, want %+v", record, wantSingle)
	}
}

func TestZonesService_UpdateRecord(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records/2", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/updateZoneRecord/success.http")

		testMethod(t, r, "PATCH")
		testHeaders(t, r)

		want := map[string]interface{}{"content": "192.168.0.10", "name": "bar"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	recordValues := ZoneRecord{Name: "bar", Content: "192.168.0.10"}

	recordResponse, err := client.Zones.UpdateRecord(accountID, "example.com", 2, recordValues)
	if err != nil {
		t.Fatalf("Zones.UpdateRecord() returned error: %v", err)
	}

	record := recordResponse.Data
	if want, got := 64784, record.ID; want != got {
		t.Fatalf("Zones.UpdateRecord() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "www", record.Name; want != got {
		t.Fatalf("Zones.UpdateRecord() returned Label expected to be `%v`, got `%v`", want, got)
	}
}

func TestZonesService_DeleteRecord(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/zones/example.com/records/2", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/deleteZoneRecord/success.http")

		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"

	_, err := client.Zones.DeleteRecord(accountID, "example.com", 2)
	if err != nil {
		t.Fatalf("Zones.DeleteRecord() returned error: %v", err)
	}
}
