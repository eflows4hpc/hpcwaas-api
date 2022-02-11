package main

import "github.com/spf13/cobra"

var executionsCmd = &cobra.Command{
	Use:     "executions",
	Aliases: []string{"execution", "exec", "e"},
	Short:   "perform operations on executions",
	Long: `perform operations on executions:

This is the root of all executions commands.
	`,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(executionsCmd)
}
