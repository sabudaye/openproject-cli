package attachment

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

const attachmentsPayload = `{
	"_type": "AttachmentCollection",
	"total": 2,
	"count": 2,
	"_embedded": {
		"elements": [
			{
				"id": 7,
				"fileName": "report.pdf",
				"fileSize": 11,
				"status": "uploaded",
				"contentType": "application/pdf",
				"digest": {"algorithm": "md5", "hash": "64c26a8403cd796ea4cf913cda2ee4a9"},
				"createdAt": "2014-05-21T08:51:20.396Z",
				"_links": {"downloadLocation": {"href": "/api/v3/attachments/7/content"}}
			},
			{
				"id": 8,
				"fileName": "cat.png",
				"fileSize": 24,
				"status": "uploaded",
				"contentType": "image/png",
				"_links": {"downloadLocation": {"href": "https://example.s3/cat.png"}}
			}
		]
	}
}`

func initAttachmentTest(t *testing.T, serverURL string) *printer.TestingPrinter {
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

func TestListAttachmentsTextOutput(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/api/v3/work_packages/42/attachments" {
			t.Errorf("request path = %s, want /api/v3/work_packages/42/attachments", request.URL.Path)
		}
		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write([]byte(attachmentsPayload))
	}))
	t.Cleanup(server.Close)

	testingPrinter := initAttachmentTest(t, server.URL)
	listWpId = "42"
	t.Cleanup(func() { listWpId = "" })

	if err := listAttachments(nil, nil); err != nil {
		t.Fatalf("listAttachments error = %v", err)
	}

	expected := printer.Red("#7") + " " +
		printer.Cyan("report.pdf") + " " +
		"(application/pdf, 11 B) " +
		"[" + printer.Yellow("uploaded") + "]\n" +
		printer.Red("#8") + " " +
		printer.Cyan("cat.png   ") + " " +
		"(image/png, 24 B) " +
		"[" + printer.Yellow("uploaded") + "]\n"
	if testingPrinter.Result != expected {
		t.Errorf("stdout = %q, want %q", testingPrinter.Result, expected)
	}
}

func TestListAttachmentsJsonOutput(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write([]byte(attachmentsPayload))
	}))
	t.Cleanup(server.Close)

	testingPrinter := initAttachmentTest(t, server.URL)
	if err := printer.InitRenderer("json"); err != nil {
		t.Fatal(err)
	}

	listWpId = "42"
	t.Cleanup(func() { listWpId = "" })

	if err := listAttachments(nil, nil); err != nil {
		t.Fatalf("listAttachments error = %v", err)
	}

	for _, fragment := range []string{`"id": 7`, `"file_name": "report.pdf"`, `"file_size": 11`, `"content_type": "application/pdf"`} {
		if !strings.Contains(testingPrinter.Result, fragment) {
			t.Errorf("stdout should contain %s, got %q", fragment, testingPrinter.Result)
		}
	}
}

func TestListAttachmentsReportsFetchErrorsOnce(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		http.Error(response, "attachment failure", http.StatusInternalServerError)
	}))
	t.Cleanup(server.Close)

	testingPrinter := initAttachmentTest(t, server.URL)
	listWpId = "42"
	t.Cleanup(func() { listWpId = "" })

	err := listAttachments(nil, nil)
	if !errors.Is(err, openerrors.ErrHandled) {
		t.Fatalf("listAttachments error = %v, want ErrHandled", err)
	}
	if count := strings.Count(testingPrinter.ErrResult, "[ERROR]"); count != 1 {
		t.Errorf("error diagnostic count = %d, want 1; stderr: %q", count, testingPrinter.ErrResult)
	}
	if testingPrinter.Result != "" {
		t.Errorf("stdout = %q, want empty", testingPrinter.Result)
	}
}

func TestListAttachmentsRejectsInvalidWorkPackageIdentifier(t *testing.T) {
	testingPrinter := &printer.TestingPrinter{}
	printer.Init(testingPrinter)

	listWpId = "not-an-id"
	t.Cleanup(func() { listWpId = "" })

	err := listAttachments(nil, nil)
	if !errors.Is(err, openerrors.ErrHandled) {
		t.Fatalf("listAttachments error = %v, want ErrHandled", err)
	}
	if !strings.Contains(testingPrinter.ErrResult, "invalid work package identifier") {
		t.Errorf("stderr should explain the invalid identifier, got %q", testingPrinter.ErrResult)
	}
}
