package printer_test

import (
	"net/url"
	"os"
	"testing"

	"github.com/sabudaye/openproject-cli/components/printer"
	"github.com/sabudaye/openproject-cli/components/routes"
)

var testingPrinter = &printer.TestingPrinter{}

func TestMain(m *testing.M) {
	apiUrl, _ := url.Parse("https://example.com")
	routes.Init(apiUrl)
	printer.Init(testingPrinter)

	os.Exit(m.Run())
}
