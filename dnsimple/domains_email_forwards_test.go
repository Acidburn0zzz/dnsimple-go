package dnsimple

import (
	"io"
	"net/http"
	//"reflect"
	"regexp"
	"testing"
)

var regexpEmail = regexp.MustCompile(`.+@.+`)

func TestDomainsService_EmailForwardsList(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/domains/example.com/email_forwards", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/listEmailForwards/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	forwardsResponse, err := client.Domains.ListEmailForwards("1010", "example.com")
	if err != nil {
		t.Fatalf("Domains.ListEmailForwards() returned error: %v", err)
	}

	forwards := forwardsResponse.Data
	if want, got := 2, len(forwards); want != got {
		t.Errorf("Domains.ListEmailForwards() expected to return %v contacts, got %v", want, got)
	}

	if want, got := 17702, forwards[0].ID; want != got {
		t.Fatalf("Domains.ListEmailForwards() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if !regexpEmail.MatchString(forwards[0].From) {
		t.Errorf("Domains.ListEmailForwards() From expected to be an email, got %v", forwards[0].From)
	}
}
