package teamweek

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var (
	client *Client
	mux    *http.ServeMux
	server *httptest.Server
)

func setup() {
	client = NewClient(nil)
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}

func teardown() {
	server.Close()
}

func TestNewClient(t *testing.T) {
	c := NewClient(nil)

	if c.BaseURL.String() != defaultBaseURL {
		t.Errorf("NewClient BaseURL = %v, want %v", c.BaseURL.String(), defaultBaseURL)
	}
	if c.UserAgent != userAgent {
		t.Errorf("NewClient UserAgent = %v, want %v", c.UserAgent, userAgent)
	}
}

func TestListAccounts(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/me/accounts.json", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[
 			{"id":1,"name":"Account 1","is_demo":false},
 			{"id":2,"name":"Account 2","is_demo":true}
    	]`)
	})

	accounts, err := client.ListAccounts()
	if err != nil {
		t.Errorf("ListAccounts returned error: %v", err)
	}

	want := []Account{
		{ID: 1, Name: "Account 1"},
		{ID: 2, Name: "Account 2", IsDemo: true},
	}

	if !reflect.DeepEqual(accounts, want) {
		t.Errorf("ListAccounts returned %+v, want %+v", accounts, want)
	}
}
