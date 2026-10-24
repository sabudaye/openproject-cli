package wptype

import (
	"github.com/spf13/cobra"

	openerrors "github.com/sabudaye/openproject-cli/components/errors"
	"github.com/sabudaye/openproject-cli/components/printer"
	"github.com/sabudaye/openproject-cli/components/resources/types"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists work package types",
	Long:  "Get a list of all work package types of the instance.",
	RunE:  listTypes,
}

func listTypes(_ *cobra.Command, _ []string) error {
	if all, err := types.All(); err == nil {
		printer.Types(all)
	} else {
		printer.Error(err)
		return openerrors.ErrHandled
	}
	return nil
}
