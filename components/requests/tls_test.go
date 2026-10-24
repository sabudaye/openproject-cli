package requests_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/sabudaye/openproject-cli/components/printer"
	"github.com/sabudaye/openproject-cli/components/requests"
)

// httptest.NewTLSServer issues a self-signed certificate that is not
// trusted by the system store. The request client must still succeed,
// because certificate verification is disabled for all API requests.
func TestClientSkipsTLSVerification(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write([]byte(`{"type":"Root"}`))
	}))
	t.Cleanup(server.Close)

	host, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	requests.Init(host, "", false)
	printer.Init(&printer.TestingPrinter{})

	body, err := requests.Get("/api/v3", nil)
	if err != nil {
		t.Fatalf("Get against self-signed TLS server returned error: %v", err)
	}
	if string(body) != `{"type":"Root"}` {
		t.Errorf("unexpected response body: %s", body)
	}
}
