package dnsimple

import (
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestContacts_contactPath(t *testing.T) {
	if want, got := "/1010/contacts", contactPath("1010", nil); want != got {
		t.Errorf("contactPath(%v,  ) = %v, want %v", "1010", got, want)
	}

	if want, got := "/1010/contacts/1", contactPath("1010", 1); want != got {
		t.Errorf("contactPath(%v, 1) = %v, want %v", "1010", got, want)
	}
}

func TestContactsService_List(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/contacts", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/listContacts/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"

	contactsResponse, err := client.Contacts.ListContacts(accountID)
	if err != nil {
		t.Fatalf("Contacts.ListContacts() returned error: %v", err)
	}

	contacts := contactsResponse.Data
	if want, got := 2, len(contacts); want != got {
		t.Errorf("Contacts.ListContacts() expected to return %v contacts, got %v", want, got)
	}

	if want, got := 1, contacts[0].ID; want != got {
		t.Fatalf("Contacts.ListContacts() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "Default", contacts[0].Label; want != got {
		t.Fatalf("Contacts.ListContacts() returned Label expected to be `%v`, got `%v`", want, got)
	}
}

func TestContactsService_Create(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/contacts", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/createContact/created.http")

		testMethod(t, r, "POST")
		testHeaders(t, r)

		want := map[string]interface{}{"label": "Default"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	contactAttributes := Contact{Label: "Default"}

	contactResponse, err := client.Contacts.CreateContact(accountID, contactAttributes)
	if err != nil {
		t.Fatalf("Contacts.CreateDomain() returned error: %v", err)
	}

	contact := contactResponse.Data
	if want, got := 1, contact.ID; want != got {
		t.Fatalf("Contacts.CreateDomain() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "Default", contact.Label; want != got {
		t.Fatalf("Contacts.CreateDomain() returned Label expected to be `%v`, got `%v`", want, got)
	}
}

func TestContactsService_Get(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/contacts/1", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/getContact/success.http")

		testMethod(t, r, "GET")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	contactID := 1

	contactResponse, err := client.Contacts.GetContact(accountID, contactID)
	if err != nil {
		t.Fatalf("Contacts.GetDomain() returned error: %v", err)
	}

	contact := contactResponse.Data
	wantSingle := &Contact{
		ID:            1,
		AccountID:     1010,
		Label:         "Default",
		FirstName:     "First",
		LastName:      "User",
		JobTitle:      "CEO",
		Organization:  "Awesome Company",
		Address1:      "Italian Street, 10",
		City:          "Roma",
		StateProvince: "RM",
		PostalCode:    "00100",
		Country:       "IT",
		Phone:         "+18001234567",
		Fax:           "+18011234567",
		Email:         "first@example.com",
		CreatedAt:     "2016-01-19T20:50:26.066Z",
		UpdatedAt:     "2016-01-19T20:50:26.066Z"}

	if !reflect.DeepEqual(contact, wantSingle) {
		t.Fatalf("Contacts.GetDomain() returned %+v, want %+v", contact, wantSingle)
	}
}

func TestContactsService_Update(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/contacts/1", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/updateContact/success.http")

		testMethod(t, r, "PATCH")
		testHeaders(t, r)

		want := map[string]interface{}{"label": "Default"}
		testRequestJSON(t, r, want)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	contactAttributes := Contact{Label: "Default"}
	accountID := "1010"
	contactID := 1

	contactResponse, err := client.Contacts.UpdateContact(accountID, contactID, contactAttributes)
	if err != nil {
		t.Fatalf("Contacts.UpdateDomain() returned error: %v", err)
	}

	contact := contactResponse.Data
	if want, got := 1, contact.ID; want != got {
		t.Fatalf("Contacts.UpdateDomain() returned ID expected to be `%v`, got `%v`", want, got)
	}
	if want, got := "Default", contact.Label; want != got {
		t.Fatalf("Contacts.UpdateDomain() returned Label expected to be `%v`, got `%v`", want, got)
	}
}

func TestContactsService_Delete(t *testing.T) {
	setupMockServer()
	defer teardownMockServer()

	mux.HandleFunc("/v2/1010/contacts/1", func(w http.ResponseWriter, r *http.Request) {
		httpResponse := httpResponseFixture(t, "/deleteContact/success.http")

		testMethod(t, r, "DELETE")
		testHeaders(t, r)

		w.WriteHeader(httpResponse.StatusCode)
		io.Copy(w, httpResponse.Body)
	})

	accountID := "1010"
	contactID := 1

	_, err := client.Contacts.DeleteContact(accountID, contactID)
	if err != nil {
		t.Fatalf("Contacts.DeleteDomain() returned error: %v", err)
	}
}
