package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eflows4hpc/hpcwaas-api/api"
	"github.com/mitchellh/mapstructure"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func init() {
	var metadata map[string]string
	sshKeyGenCmd := &cobra.Command{
		Use:     "key-gen",
		Aliases: []string{"gen", "gen-key", "g"},
		Short:   "Generate a new SSH key pair",
		Long: `Generate a new SSH key pair

This will create an SSH Key Pair.
The private key is stored in HashiCorp Vault and *only* the public key is returned to the user along with a randomly generated
identifier for the key.
The public key and the key identifier can not be seen again so take note of it.
Once the user has the public key, he should add it to the authorized_keys file on the systems he want to run his workflows.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.GetClient(*clientConfig)
			if err != nil {
				errExit(err)
			}
			var m map[string]interface{}
			err = mapstructure.Decode(metadata, &m)
			if err != nil {
				errExit(err)
			}

			err = generateSSHKeyPair(client, m, output)
			return handleUsageError(cmd, err)
		},
	}
	sshKeyGenCmd.Flags().StringToStringVarP(&metadata, "metadata", "m", nil, "Metadata to be stored along with the SSH key")
	sshKeysCmd.AddCommand(sshKeyGenCmd)
}

func generateSSHKeyPair(client api.HTTPClient, metadata map[string]interface{}, output string) error {
	sshKey, err := client.SSHKeys().GenerateSSHKey(context.Background(), api.SSHKeyGenerationRequest{MetaData: metadata})
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
	pterm.Info.Printf("SSH key ID:     %s\n", sshKey.ID)
	pterm.Info.Printf("SSH Public key: %s\n", sshKey.PublicKey)
	return nil
}
