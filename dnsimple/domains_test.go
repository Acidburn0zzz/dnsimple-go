package dnsimple

import (
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestDomains_domainPath(t *testing.T) {
	actual := domainPath("1", nil)
	expected := "/1/domains"

	if actual != expected {
		t.Errorf("domainPath(\"1\", nil): actual %s, expected %s", actual, expected)
	}

	actual = domainPath("1", "example.com")
	expected = "/1/domains/example.com"

	if actual != expected {
		t.Errorf("domainPath(\"1\", \"example.com\", nil): actual %s, expected %s", actual, expected)
	}

	actual = domainPath("1", 1)
	expected = "/1/domains/1"

	if actual != expected {
		t.Errorf("domainPath(\"1\", 1, nil): actual %s, expected %s", actual, expected)
	}
}

func TestDomainsService_List(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/listDomains/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"

	domainsResponse, err := client.Domains.ListDomains(accountID)
	if err != nil {
		t.Fatalf("Domains.ListDomains() returned error: %v", err)
	}

	domains := domainsResponse.Data
	if want, got := 2, len(domains); want != got {
		t.Errorf("Domains.ListDomains() expected to return %v contacts, got %v", want, got)
	}

	if want, got := 1, domains[0].ID; want != got {
		t.Fatalf("Domains.ListDomains() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "example-alpha.com", domains[0].Name; want != got {
		t.Fatalf("Domains.ListDomains() returned Name expected to be `%v`, got `%v`", want, got)
	}
}

func TestDomainsService_Create(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1/domains", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/createDomain/created.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"name": "example.com"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1"
	domainAttributes := Domain{Name: "example.com"}

	domainResponse, err := client.Domains.CreateDomain(accountID, domainAttributes)
	if err != nil {
		t.Fatalf("Domains.Create() returned error: %v", err)
	}

	domain := domainResponse.Data
	if want, got := 1, domain.ID; want != got {
		t.Fatalf("Domains.Create() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "example-alpha.com", domain.Name; want != got {
		t.Fatalf("Domains.Create() returned Name expected to be `%v`, got `%v`", want, got)
	}
}

func TestDomainsService_Get(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/getDomain/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"

	domainResponse, err := client.Domains.GetDomain(accountID, "example.com")
	if err != nil {
		t.Errorf("Domains.Get() returned error: %v", err)
	}

	domain := domainResponse.Data
	wantSingle := &Domain{
		ID:           1,
		AccountID:    1010,
		RegistrantID: 0,
		Name:         "example-alpha.com",
		UnicodeName:  "example-alpha.com",
		Token:        "domain-token",
		State:        "hosted",
		PrivateWhois: false,
		ExpiresOn:    "",
		CreatedAt:    "2014-12-06T15:56:55.573Z",
		UpdatedAt:    "2015-12-09T00:20:56.056Z"}

	if !reflect.DeepEqual(domain, wantSingle) {
		t.Fatalf("Domains.Get() returned %+v, want %+v", domain, wantSingle)
	}
}

func TestDomainsService_Delete(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/deleteDomain/success.http")

		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"

	_, err := client.Domains.DeleteDomain(accountID, "example.com")
	if err != nil {
		t.Fatalf("Domains.Delete() returned error: %v", err)
	}
}

func TestDomainsService_ResetDomainToken(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/token", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/resetDomainToken/success.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"

	domainResponse, err := client.Domains.ResetDomainToken(accountID, "example.com")
	if err != nil {
		t.Fatalf("Domains.ResetDomainToken() returned error: %v", err)
	}

	domain := domainResponse.Data
	if want, got := 1, domain.ID; want != got {
		t.Fatalf("Domains.ResetDomainToken() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "example-alpha.com", domain.Name; want != got {
		t.Fatalf("Domains.ResetDomainToken() returned Name expected to be `%v`, got `%v`", want, got)
	}
}
