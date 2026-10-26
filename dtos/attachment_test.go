package dtos_test

import (
	"encoding/json"
	"testing"

	"github.com/sabudaye/openproject-cli/dtos"
)

func TestAttachmentDtoConvertPrefersFileSize(t *testing.T) {
	var dto dtos.AttachmentDto
	payload := `{"id": 7, "fileName": "report.pdf", "fileSize": 11, "filesize": 99}`
	if err := json.Unmarshal([]byte(payload), &dto); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if got := dto.Convert().FileSize; got != 11 {
		t.Errorf("FileSize = %d, want 11 (fileSize must win over filesize)", got)
	}
}

// The APIv3 spec declares the attribute as "fileSize" but its own examples
// serialize "filesize"; instances follow either spelling.
func TestAttachmentDtoConvertFallsBackToFilesize(t *testing.T) {
	var dto dtos.AttachmentDto
	payload := `{"id": 7, "fileName": "report.pdf", "filesize": 42}`
	if err := json.Unmarshal([]byte(payload), &dto); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if got := dto.Convert().FileSize; got != 42 {
		t.Errorf("FileSize = %d, want 42 (legacy filesize must be picked up)", got)
	}
}

func TestAttachmentDtoConvertKeepsZeroFileSize(t *testing.T) {
	var dto dtos.AttachmentDto
	payload := `{"id": 7, "fileName": "empty.bin", "fileSize": 0}`
	if err := json.Unmarshal([]byte(payload), &dto); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if got := dto.Convert().FileSize; got != 0 {
		t.Errorf("FileSize = %d, want 0", got)
	}
}

func TestAttachmentDtoConvertMapsDownloadLocation(t *testing.T) {
	var dto dtos.AttachmentDto
	payload := `{
		"id": 7,
		"fileName": "report.pdf",
		"_links": {"downloadLocation": {"href": "https://remote.example/file"}}
	}`
	if err := json.Unmarshal([]byte(payload), &dto); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if got := dto.Convert().DownloadHref; got != "https://remote.example/file" {
		t.Errorf("DownloadHref = %q, want %q", got, "https://remote.example/file")
	}
}
