package attachment

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"

	openerrors "github.com/sabudaye/openproject-cli/components/errors"
	"github.com/sabudaye/openproject-cli/components/printer"
	"github.com/sabudaye/openproject-cli/components/resources/attachments"
)

var downloadOutput string

var downloadCmd = &cobra.Command{
	Use:   "download <attachment-id>",
	Short: "Downloads an attachment to a file",
	Long: `Download the attachment with the given id.
The content is written to --output, or to the attachment's
file name when the flag is omitted.`,
	Args: cobra.ExactArgs(1),
	RunE: downloadAttachment,
}

func downloadAttachment(_ *cobra.Command, args []string) error {
	attachment, err := attachments.Lookup(args[0])
	if err != nil {
		printer.Error(err)
		return openerrors.ErrHandled
	}

	output := downloadOutput
	if output == "" {
		if attachment.FileName == "" {
			printer.ErrorText("The attachment has no file name. Specify the target file with --output.")
			return openerrors.ErrHandled
		}
		output = attachment.FileName
	}

	file, err := os.Create(output)
	if err != nil {
		printer.ErrorText(err.Error())
		return openerrors.ErrHandled
	}

	printer.Info("Downloading attachment ...")
	written, err := attachments.Download(attachment, file)
	_ = file.Close()

	if err != nil {
		// Do not leave a partial file behind.
		_ = os.Remove(output)
		printer.Error(err)
		return openerrors.ErrHandled
	}

	if attachment.FileSize > 0 && written != attachment.FileSize {
		printer.Warning(
			"downloaded " + strconv.FormatInt(written, 10) + " bytes, server reports " + strconv.FormatInt(attachment.FileSize, 10) + " bytes",
		)
	}
	printer.Info("Saved " + printer.Cyan(attachment.FileName) + " to " + output)
	printer.Done()
	return nil
}
