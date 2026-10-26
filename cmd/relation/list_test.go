package relation

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	openerrors "github.com/sabudaye/openproject-cli/components/errors"
	"github.com/sabudaye/openproject-cli/components/printer"
	"github.com/sabudaye/openproject-cli/components/requests"
)

const relationsPayload = `{
	"_type": "Collection",
	"total": 2,
	"count": 2,
	"pageSize": -1,
	"offset": 0,
	"_embedded": {
		"elements": [
			{
				"id": 650,
				"name": "duplicated by",
				"type": "duplicated",
				"reverseType": "duplicates",
				"lag": 2,
				"description": "Those are the same...",
				"_links": {
					"from": {"href": "/api/v3/work_packages/108", "title": "Destroy Rebel Base"},
					"to": {"href": "/api/v3/work_packages/78", "title": "Annihilate Yavin 4"}
				}
			},
			{
				"id": 651,
				"name": "relates to",
				"type": "relates",
				"reverseType": "relates",
				"_embedded": {
					"from": {"id": 42},
					"to": {"id": 77}
				}
			}
		]
	}
}`

func initRelationTest(t *testing.T, serverURL string) *printer.TestingPrinter {
	t.Helper()

	host, err := url.Parse(serverURL)
	if err != nil {
		t.Fatal(err)
	}
	requests.Init(host, "secret-token", false)

	testingPrinter := &printer.TestingPrinter{}
	printer.Init(testingPrinter)
	t.Cleanup(testingPrinter.Reset)

	if err := printer.InitRenderer("text"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = printer.InitRenderer("text") })

	return testingPrinter
}

func TestListRelationsTextOutput(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/v3/relations" {
			t.Errorf("request path = %s, want /api/v3/relations", request.URL.Path)
		}

		filters := request.URL.Query().Get("filters")
		if !strings.Contains(filters, `"involved"`) || !strings.Contains(filters, `"operator":"="`) || !strings.Contains(filters, `"42"`) {
			t.Errorf("filters should select the involved work package 42, got %q", filters)
		}

		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write([]byte(relationsPayload))
	}))
	t.Cleanup(server.Close)

	testingPrinter := initRelationTest(t, server.URL)
	listWpId = "42"
	t.Cleanup(func() { listWpId = "" })

	if err := listRelations(nil, nil); err != nil {
		t.Fatalf("listRelations error = %v", err)
	}

	expected := printer.Red("#650") + " " +
		printer.Red("#108") + " " +
		printer.Cyan("Destroy Rebel Base") + " " +
		"[" + printer.Yellow("duplicated by") + "] " +
		printer.Red("#78") + " " +
		printer.Cyan("Annihilate Yavin 4") + " " +
		printer.Green("(lag: 2d)") + " " +
		printer.Blue(`"Those are the same..."`) + "\n" +
		printer.Red("#651") + " " +
		printer.Red("#42") + " " +
		"[" + printer.Yellow("relates to") + "] " +
		printer.Red("#77") + "\n"
	if testingPrinter.Result != expected {
		t.Errorf("stdout = %q, want %q", testingPrinter.Result, expected)
	}
}

func TestListRelationsJsonOutput(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write([]byte(relationsPayload))
	}))
	t.Cleanup(server.Close)

	testingPrinter := initRelationTest(t, server.URL)
	if err := printer.InitRenderer("json"); err != nil {
		t.Fatal(err)
	}

	listWpId = "42"
	t.Cleanup(func() { listWpId = "" })

	if err := listRelations(nil, nil); err != nil {
		t.Fatalf("listRelations error = %v", err)
	}

	for _, fragment := range []string{
		`"id": 650`,
		`"name": "duplicated by"`,
		`"type": "duplicated"`,
		`"reverse_type": "duplicates"`,
		`"lag": 2`,
		`"id": 108`,
		`"title": "Destroy Rebel Base"`,
		`"id": 78`,
		`"title": "Annihilate Yavin 4"`,
		`"id": 42`,
		`"id": 77`,
	} {
		if !strings.Contains(testingPrinter.Result, fragment) {
			t.Errorf("stdout should contain %s, got %q", fragment, testingPrinter.Result)
		}
	}
}

func TestListRelationsReportsFetchErrorsOnce(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		http.Error(response, "relation failure", http.StatusInternalServerError)
	}))
	t.Cleanup(server.Close)

	testingPrinter := initRelationTest(t, server.URL)
	listWpId = "42"
	t.Cleanup(func() { listWpId = "" })

	err := listRelations(nil, nil)
	if !errors.Is(err, openerrors.ErrHandled) {
		t.Fatalf("listRelations error = %v, want ErrHandled", err)
	}
	if count := strings.Count(testingPrinter.ErrResult, "[ERROR]"); count != 1 {
		t.Errorf("error diagnostic count = %d, want 1; stderr: %q", count, testingPrinter.ErrResult)
	}
	if testingPrinter.Result != "" {
		t.Errorf("stdout = %q, want empty", testingPrinter.Result)
	}
}

func TestListRelationsRejectsInvalidWorkPackageIdentifier(t *testing.T) {
	testingPrinter := &printer.TestingPrinter{}
	printer.Init(testingPrinter)
	t.Cleanup(testingPrinter.Reset)

	listWpId = "not-an-id"
	t.Cleanup(func() { listWpId = "" })

	err := listRelations(nil, nil)
	if !errors.Is(err, openerrors.ErrHandled) {
		t.Fatalf("listRelations error = %v, want ErrHandled", err)
	}
	if !strings.Contains(testingPrinter.ErrResult, "invalid work package identifier") {
		t.Errorf("stderr should explain the invalid identifier, got %q", testingPrinter.ErrResult)
	}
}
