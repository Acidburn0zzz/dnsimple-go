package dnsimple

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestContacts_contactPath(t *testing.T) {
	var pathTests = []struct {
		input    interface{}
		expected string
	}{
		{nil, "contacts"},
		{1, "contacts/1"},
	}

	for _, pt := range pathTests {
		actual := contactPath(pt.input)
		if actual != pt.expected {
			t.Errorf("contactPath(%+v): expected %s, actual %s", pt.input, pt.expected)
		}
	}
}

func TestContactsService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/contacts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"contact":{"id":1,"user_id":21,"label":"Simone","first_name":"Simone","last_name":"Carletti","job_title":"Underwater Programmer","organization_name":"DNSimple","email_address":"simone.carletti@dnsimple.com","phone":"+1 111 4567890","fax":"+1 222 4567890","address1":"Awesome Street","address2":"c/o Someone","city":"Rome","state_province":"RM","postal_code":"00171","country":"IT"}},{"contact":{"id":2,"user_id":22,"label":"Simone","first_name":"Simone","last_name":"Carletti","job_title":"Underwater Programmer","organization_name":"DNSimple","email_address":"simone.carletti@dnsimple.com","phone":"+1 111 4567890","fax":"+1 222 4567890","address1":"Awesome Street","address2":"c/o Someone","city":"Rome","state_province":"RM","postal_code":"00171","country":"IT"}}]`)
	})

	contacts, _, err := client.Contacts.List()

	if err != nil {
		t.Errorf("Contacts.List returned error: %v", err)
	}

	if size, want := len(contacts), 2; size != want {
		t.Fatalf("Contacts.List returned %+v items, want %+v", size, want)
	}

	for i, item := range contacts {
		if kind, want := reflect.TypeOf(item).Name(), "Contact"; kind != want {
			t.Errorf("Contacts.List expected [%d] to be %s, got %s", i, want, kind)
		}
	}
}

func TestContactsService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/contacts/1", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "GET")
			fmt.Fprint(w, `{"contact":{"id":1,"user_id":21,"label":"Default","first_name":"Simone","last_name":"Carletti","job_title":"Underwater Programmer","organization_name":"DNSimple","email_address":"simone.carletti@dnsimple.com","phone":"+1 111 4567890","fax":"+1 222 4567890","address1":"Awesome Street","address2":"c/o Someone","city":"Rome","state_province":"RM","postal_code":"00171","country":"IT"}}`)
		})

	contact, _, err := client.Contacts.Get(1)

	if err != nil {
		t.Errorf("Contacts.Get returned error: %v", err)
	}

	want := Contact{
		Id:           1,
		Label:        "Default",
		FirstName:    "Simone",
		LastName:     "Carletti",
		JobTitle:     "Underwater Programmer",
		Organization: "DNSimple",
		Email:        "simone.carletti@dnsimple.com",
		Phone:        "+1 111 4567890",
		Fax:          "+1 222 4567890",
		Address1:     "Awesome Street",
		Address2:     "c/o Someone",
		City:         "Rome",
		Zip:          "00171",
		Country:      "IT"}
	if !reflect.DeepEqual(contact, want) {
		t.Fatalf("Contacts.Get returned %+v, want %+v", contact, want)
	}
}
