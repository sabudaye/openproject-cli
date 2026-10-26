package printer

import (
	"fmt"
	"strings"

	"github.com/sabudaye/openproject-cli/models"
)

func Attachment(a *models.Attachment) {
	activeRenderer.Attachment(a)
}

func Attachments(as []*models.Attachment) {
	activeRenderer.Attachments(as)
}

func printAttachment(a *models.Attachment, maxIdLength, maxNameLength int) {
	diff := maxIdLength - idLength(a.Id)
	parts := []string{Red(fmt.Sprintf("%s#%d", indent(diff), a.Id))}

	nameDiff := maxNameLength - len(a.FileName)
	parts = append(parts, Cyan(a.FileName+indent(nameDiff)))

	meta := fmt.Sprintf("(%s", a.ContentType)
	if a.FileSize > 0 {
		meta += ", " + formatFileSize(a.FileSize)
	}
	meta += ")"
	parts = append(parts, meta)

	if a.Status != "" {
		parts = append(parts, "["+Yellow(a.Status)+"]")
	}

	activePrinter.Println(strings.Join(parts, " "))
}

func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
