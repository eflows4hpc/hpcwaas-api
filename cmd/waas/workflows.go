package main

import "github.com/spf13/cobra"

var workflowsCmd = &cobra.Command{
	Use:     "workflows",
	Aliases: []string{"workflow", "wf", "w"},
	Short:   "perform operations on workflows",
	Long: `perform operations on workflows:

This is the root of all workflows commands. But this is also an alias for the "workflows list" command.
	`,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return workflowsListCmd.RunE(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(workflowsCmd)
}
