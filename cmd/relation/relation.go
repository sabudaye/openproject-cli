package relation

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "relation [verb]",
	Short: "Manage work package relations",
	Long:  "List relations of a work package (--work-package).",
}

func init() {
	listCmd.Flags().StringVarP(
		&listWpId,
		"work-package",
		"w",
		"",
		"Work package ID or identifier to list relations for",
	)
	_ = listCmd.MarkFlagRequired("work-package")

	RootCmd.AddCommand(listCmd)
}
