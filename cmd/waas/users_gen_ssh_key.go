package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eflows4hpc/hpcwaas-api/api"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func init() {
	usersKeyGenCmd := &cobra.Command{
		Use:     "key-gen <user-id>",
		Aliases: []string{"gen", "gen-key", "ssh-key", "g"},
		Short:   "Generate a new SSH key for a user",
		Long: `Generate a new SSH key for a user

This will create an SSH Key Pair for a given user.
The private key is stored in HashiCorp Vault and *only* the public key is returned to the user.
The public key can not be seen again so take note of it.
Once the user has the public key, he should add it to the authorized_keys file on the systems he want to run his workflows.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.GetClient(*clientConfig)
			if err != nil {
				errExit(err)
			}
			err = generateSSHKeyPair(client, args[0], output)
			return handleUsageError(cmd, err)
		},
	}

	usersCmd.AddCommand(usersKeyGenCmd)
}

func generateSSHKeyPair(client api.HTTPClient, userName, output string) error {
	sshKey, err := client.Users().GenerateSSHKey(context.Background(), userName)
	if err != nil {
		return err
	}
	if output == JSONOutput {
		data, err := json.Marshal(sshKey)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	pterm.Info.Println("Below is your newly generated SSH public key.")
	pterm.Info.Println("Take note of it as you will not see it again.")
	pterm.Info.Println("You are responsible for adding it to the authorized_keys file on the systems you want to run your workflows.")
	pterm.Info.Printf("SSH Public key: %s", sshKey.PublicKey)
	return nil
}
