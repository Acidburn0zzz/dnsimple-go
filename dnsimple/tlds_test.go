package dnsimple

import (
	"io"
	"net/http"
	"testing"
)

func TestTldsService_ListTlds(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/tlds", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/listTlds/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	tldsResponse, err := client.Tlds.ListTlds()
	if err != nil {
		t.Fatalf("Tlds.ListTlds() returned error: %v", err)
	}

	tlds := tldsResponse.Data
	if want, got := 2, len(tlds); want != got {
		t.Errorf("Tlds.ListTlds() expected to return %v TLDs, got %v", want, got)
	}

	if want, got := "ac", tlds[0].Tld; want != got {
		t.Fatalf("Tlds.ListTlds() returned Tld expected to be `%v`, got `%v`", want, got)
	}
}
func TestTldsService_GetTld(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/tlds/com", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/getTld/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	tldResponse, err := client.Tlds.GetTld("com")
	if err != nil {
		t.Fatalf("Tlds.GetTlds() returned error: %v", err)
	}

	tld := tldResponse.Data
	if want, got := "com", tld.Tld; want != got {
		t.Fatalf("Tlds.GetTlds() returned Tld expected to be `%v`, got `%v`", want, got)
	}
}
