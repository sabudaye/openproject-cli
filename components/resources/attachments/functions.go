package attachments

import (
	"fmt"
	"io"

	"github.com/sabudaye/openproject-cli/components/parser"
	"github.com/sabudaye/openproject-cli/components/paths"
	"github.com/sabudaye/openproject-cli/components/requests"
	"github.com/sabudaye/openproject-cli/dtos"
	"github.com/sabudaye/openproject-cli/models"
)

func ByWorkPackage(id string) ([]*models.Attachment, error) {
	response, err := requests.Get(paths.WorkPackageAttachments(id), nil)
	if err != nil {
		return nil, err
	}

	collection := parser.Parse[dtos.AttachmentCollectionDto](response)
	return collection.Convert(), nil
}

func Lookup(id string) (*models.Attachment, error) {
	response, err := requests.Get(paths.Attachment(id), nil)
	if err != nil {
		return nil, err
	}

	attachment := parser.Parse[dtos.AttachmentDto](response)
	return attachment.Convert(), nil
}

// Download streams the attachment content into out and returns the number
// of bytes written.
func Download(attachment *models.Attachment, out io.Writer) (int64, error) {
	if attachment.DownloadHref == "" {
		return 0, fmt.Errorf("attachment %d has no download location", attachment.Id)
	}

	return requests.DownloadFile(attachment.DownloadHref, out)
}
