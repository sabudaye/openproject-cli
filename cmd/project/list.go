package project

import (
	"github.com/spf13/cobra"

	openerrors "github.com/sabudaye/openproject-cli/components/errors"
	"github.com/sabudaye/openproject-cli/components/printer"
	"github.com/sabudaye/openproject-cli/components/resources/projects"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists projects",
	Long: `Get a list of visible projects.
The list can get filtered further.`,
	RunE: listProjects,
}

func listProjects(_ *cobra.Command, _ []string) error {
	if all, err := projects.All(); err == nil {
		printer.Projects(all)
	} else {
		printer.Error(err)
		return openerrors.ErrHandled
	}
	return nil
}
