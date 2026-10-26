package attachment

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "attachment [verb]",
	Short: "Manage attachments",
	Long:  "List and download attachments of a work package.",
}

func init() {
	listCmd.Flags().StringVarP(
		&listWpId,
		"work-package",
		"w",
		"",
		"Work package ID or identifier to list attachments for",
	)
	_ = listCmd.MarkFlagRequired("work-package")

	downloadCmd.Flags().StringVarP(
		&downloadOutput,
		"output",
		"o",
		"",
		"Target file for the downloaded attachment (default: the attachment's file name)",
	)

	RootCmd.AddCommand(listCmd, downloadCmd)
}
