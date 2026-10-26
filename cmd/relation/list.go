package relation

import (
	"github.com/spf13/cobra"

	openerrors "github.com/sabudaye/openproject-cli/components/errors"
	"github.com/sabudaye/openproject-cli/components/printer"
	"github.com/sabudaye/openproject-cli/components/resources/relations"
	"github.com/sabudaye/openproject-cli/components/resources/work_packages"
)

var listWpId string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists relations of a work package",
	Long:  `Get a list of relations in which the work package given by --work-package is involved (on either side).`,
	RunE:  listRelations,
}

func listRelations(_ *cobra.Command, _ []string) error {
	if err := work_packages.ValidateIdentifier(listWpId); err != nil {
		printer.ErrorText(err.Error())
		return openerrors.ErrHandled
	}

	list, err := relations.ByWorkPackage(listWpId)
	if err != nil {
		printer.Error(err)
		return openerrors.ErrHandled
	}

	printer.Relations(list)
	return nil
}
