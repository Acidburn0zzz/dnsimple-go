package dnsimple

import (
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestRegistrarService_GetDomainDelegation(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/delegation", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/getDomainDelegation/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	delegationResponse, err := client.Registrar.GetDomainDelegation("1010", "example.com")
	if err != nil {
		t.Fatalf("Registrar.GetDomainDelegation() returned error: %v", err)
	}

	delegation := delegationResponse.Data
	wantSingle := &Delegation{"ns1.dnsimple.com", "ns2.dnsimple.com", "ns3.dnsimple.com", "ns4.dnsimple.com"}

	if !reflect.DeepEqual(delegation, wantSingle) {
		t.Fatalf("Registrar.GetDomainDelegation() returned %+v, want %+v", delegation, wantSingle)
	}
}

func TestRegistrarService_ChangeDomainDelegation(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/registrar/domains/example.com/delegation", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/changeDomainDelegation/success.http")

		testMethod(t, r, "PUT")
		testHeaders(t, r)

		want := []interface{}{"ns1.dnsimple.com", "ns2.dnsimple.com"}
		testRequestJSONArray(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	newDelegation := &Delegation{"ns1.dnsimple.com", "ns2.dnsimple.com"}

	delegationResponse, err := client.Registrar.ChangeDomainDelegation("1010", "example.com", newDelegation)
	if err != nil {
		t.Fatalf("Registrar.ChangeDomainDelegation() returned error: %v", err)
	}

	delegation := delegationResponse.Data
	wantSingle := &Delegation{"ns1.dnsimple.com", "ns2.dnsimple.com", "ns3.dnsimple.com", "ns4.dnsimple.com"}

	if !reflect.DeepEqual(delegation, wantSingle) {
		t.Fatalf("Registrar.ChangeDomainDelegation() returned %+v, want %+v", delegation, wantSingle)
	}
}
