package main

import "github.com/spf13/cobra"

var sshKeysCmd = &cobra.Command{
	Use:     "ssh_keys",
	Aliases: []string{"keys", "key", "k"},
	Short:   "perform operations on ssh keys",
	Long: `perform operations on ssh keys:

This is the root of all ssh keys commands.
	`,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(sshKeysCmd)
}
