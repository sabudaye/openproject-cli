package attachment

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"

	openerrors "github.com/sabudaye/openproject-cli/components/errors"
)

const fileContent = "hello world"

func attachmentDetails(downloadHref string) string {
	return `{
		"_type": "Attachment",
		"id": 7,
		"fileName": "report.pdf",
		"fileSize": ` + strconv.Itoa(len(fileContent)) + `,
		"status": "uploaded",
		"contentType": "application/pdf",
		"_links": {"downloadLocation": {"href": "` + downloadHref + `"}}
	}`
}

func TestDownloadWritesFileWithInstanceAuth(t *testing.T) {
	var receivedAuth string
	var authMu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/api/v3/attachments/7":
			response.Header().Set("Content-Type", "application/json")
			_, _ = response.Write([]byte(attachmentDetails("/api/v3/attachments/7/content")))
		case "/api/v3/attachments/7/content":
			authMu.Lock()
			receivedAuth = request.Header.Get("Authorization")
			authMu.Unlock()
			response.Header().Set("Content-Type", "application/octet-stream")
			_, _ = response.Write([]byte(fileContent))
		default:
			t.Errorf("unexpected request path %s", request.URL.Path)
			http.NotFound(response, request)
		}
	}))
	t.Cleanup(server.Close)

	initAttachmentTest(t, server.URL)

	output := filepath.Join(t.TempDir(), "report.pdf")
	downloadOutput = output
	t.Cleanup(func() { downloadOutput = "" })

	if err := downloadAttachment(nil, []string{"7"}); err != nil {
		t.Fatalf("downloadAttachment error = %v", err)
	}

	content, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("reading downloaded file: %v", err)
	}
	if string(content) != fileContent {
		t.Errorf("file content = %q, want %q", content, fileContent)
	}

	// The instance host must receive the API token as basic auth.
	if !strings.Contains(receivedAuth, "Basic ") {
		t.Errorf("instance content request should carry basic auth, got %q", receivedAuth)
	}
}

func TestDownloadRemoteLocationSendsNoAuth(t *testing.T) {
	remoteAuths := []string{}
	var authMu sync.Mutex

	remote := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		authMu.Lock()
		remoteAuths = append(remoteAuths, request.Header.Get("Authorization"))
		authMu.Unlock()
		_, _ = response.Write([]byte(fileContent))
	}))
	t.Cleanup(remote.Close)

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write([]byte(attachmentDetails(remote.URL + "/cat.png")))
	}))
	t.Cleanup(server.Close)

	initAttachmentTest(t, server.URL)

	output := filepath.Join(t.TempDir(), "cat.png")
	downloadOutput = output
	t.Cleanup(func() { downloadOutput = "" })

	if err := downloadAttachment(nil, []string{"7"}); err != nil {
		t.Fatalf("downloadAttachment error = %v", err)
	}

	authMu.Lock()
	defer authMu.Unlock()
	if len(remoteAuths) != 1 {
		t.Fatalf("remote location hit %d times, want 1", len(remoteAuths))
	}
	if remoteAuths[0] != "" {
		t.Errorf("remote download location must not receive the API token, got %q", remoteAuths[0])
	}
}

func TestDownloadFailureRemovesPartialFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/api/v3/attachments/7" {
			response.Header().Set("Content-Type", "application/json")
			_, _ = response.Write([]byte(attachmentDetails("/api/v3/attachments/7/content")))
			return
		}
		http.Error(response, "content failure", http.StatusInternalServerError)
	}))
	t.Cleanup(server.Close)

	testingPrinter := initAttachmentTest(t, server.URL)

	output := filepath.Join(t.TempDir(), "report.pdf")
	downloadOutput = output
	t.Cleanup(func() { downloadOutput = "" })

	err := downloadAttachment(nil, []string{"7"})
	if !errors.Is(err, openerrors.ErrHandled) {
		t.Fatalf("downloadAttachment error = %v, want ErrHandled", err)
	}
	if _, statErr := os.Stat(output); !os.IsNotExist(statErr) {
		t.Errorf("partial file should be removed, stat error = %v", statErr)
	}
	if testingPrinter.Result != "" {
		t.Errorf("stdout = %q, want empty", testingPrinter.Result)
	}
}

func TestDownloadDefaultsOutputToFileName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/api/v3/attachments/7" {
			response.Header().Set("Content-Type", "application/json")
			_, _ = response.Write([]byte(attachmentDetails("/api/v3/attachments/7/content")))
			return
		}
		_, _ = response.Write([]byte(fileContent))
	}))
	t.Cleanup(server.Close)

	initAttachmentTest(t, server.URL)

	dir := t.TempDir()
	t.Chdir(dir)

	if err := downloadAttachment(nil, []string{"7"}); err != nil {
		t.Fatalf("downloadAttachment error = %v", err)
	}

	content, err := os.ReadFile("report.pdf")
	if err != nil {
		t.Fatalf("default output file missing: %v", err)
	}
	if string(content) != fileContent {
		t.Errorf("file content = %q, want %q", content, fileContent)
	}
}
