package main

import "github.com/spf13/cobra"

var usersCmd = &cobra.Command{
	Use:     "users",
	Aliases: []string{"user", "u"},
	Short:   "perform operations on users",
	Long: `perform operations on users:

This is the root of all users commands.
	`,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(usersCmd)
}
