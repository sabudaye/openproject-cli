package attachment

import (
	"github.com/spf13/cobra"

	openerrors "github.com/sabudaye/openproject-cli/components/errors"
	"github.com/sabudaye/openproject-cli/components/printer"
	"github.com/sabudaye/openproject-cli/components/resources/attachments"
	"github.com/sabudaye/openproject-cli/components/resources/work_packages"
)

var listWpId string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists attachments of a work package",
	Long:  `Get a list of attachments of the work package given by --work-package.`,
	RunE:  listAttachments,
}

func listAttachments(_ *cobra.Command, _ []string) error {
	if err := work_packages.ValidateIdentifier(listWpId); err != nil {
		printer.ErrorText(err.Error())
		return openerrors.ErrHandled
	}

	list, err := attachments.ByWorkPackage(listWpId)
	if err != nil {
		printer.Error(err)
		return openerrors.ErrHandled
	}

	printer.Attachments(list)
	return nil
}
