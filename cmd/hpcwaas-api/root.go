/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/eflows4hpc/hpcwaas-api/pkg/managers/a4c"
	"github.com/eflows4hpc/hpcwaas-api/pkg/managers/vault"
	"github.com/eflows4hpc/hpcwaas-api/pkg/rest"
)

func main() {
	Execute()
}

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hpcwaas-api",
	Short: "A REST API for managing HPC Workflows as a Service",
	Long: `A REST API for managing HPC Workflows as a Service.

This is the server side of the HPC Workflows as a Service (HPCWAAS) API.
It interacts with Hashicorp Vault to manage credentials and with Alien4Cloud to manage the execution of HPC Workflows.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		conf := &rest.Config{}
		err := viper.Unmarshal(conf)
		if err != nil {
			return errors.Wrapf(err, "failed to parse config file: %q", viper.ConfigFileUsed())
		}
		s := &rest.Server{Config: conf}
		return s.StartServer()

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

const listenAddressFlagName = "listen_address"
const a4cAddressFlagName = "a4c_address"
const a4cAddressConfigName = "alien_config.address"
const a4cUserFlagName = "a4c_user"
const a4cUserConfigName = "alien_config.user"
const a4cPasswordFlagName = "a4c_password"
const a4cPasswordConfigName = "alien_config.password"
const a4cCAFileFlagName = "a4c_ca_file"
const a4cCAFileConfigName = "alien_config.ca_file"
const a4cSkipSecureFlagName = "a4c_skip_secure"
const a4cSkipSecureConfigName = "alien_config.skip_secure"
const vaultAddressFlagName = "vault_address"
const vaultAddressConfigName = "vault_config.address"
const vaultRoleIDFlagName = "vault_role_id"
const vaultRoleIDConfigName = "vault_config.role_id"
const vaultSecretIDFlagName = "vault_secret_id"
const vaultSecretIDConfigName = "vault_config.secret_id"
const vaultIsSecretWrappedFlagName = "vault_is_secret_wrapped"
const vaultIsSecretWrappedConfigName = "vault_config.is_secret_wrapped"

/*
	Address    string `mapstructure:"address"`
	User       string `mapstructure:"user"`
	Password   string `mapstructure:"password"`
	CaFile     string `mapstructure:"ca_file"`
	SkipSecure bool   `mapstructure:"skip_secure"`
*/

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hpcwaas-api.yaml)")
	rootCmd.PersistentFlags().StringP(listenAddressFlagName, "l", rest.DefaultListenAddress, "Listening address for the HTTP REST API")

	// Alien4Cloud config flag
	rootCmd.PersistentFlags().String(a4cAddressFlagName, a4c.DefaultAddress, "Address of the Alien4Cloud HTTP REST API")
	rootCmd.PersistentFlags().String(a4cUserFlagName, a4c.DefaultUser, "User to connect to Alien4Cloud")
	rootCmd.PersistentFlags().String(a4cPasswordFlagName, a4c.DefaultPassword, "Password to connect to Alien4Cloud")
	rootCmd.PersistentFlags().String(a4cCAFileFlagName, "", "CA File to use to validate SSL certificates")
	rootCmd.PersistentFlags().Bool(a4cSkipSecureFlagName, false, "Either or not to skip SSL certificates validation")

	// Vault config flag
	rootCmd.PersistentFlags().String(vaultAddressFlagName, vault.DefaultAddress, "Address of the Vault REST API")
	rootCmd.PersistentFlags().String(vaultRoleIDFlagName, "", "RoleID to connect to Vault")
	rootCmd.PersistentFlags().String(vaultSecretIDFlagName, "", "SecretID to connect to Vault")
	rootCmd.PersistentFlags().Bool(vaultIsSecretWrappedFlagName, false, "Either or not the provided secret ID is wrapped")

	// Global flags/config binding
	viper.BindPFlag(listenAddressFlagName, rootCmd.PersistentFlags().Lookup(listenAddressFlagName))

	// Alien4Cloud config
	viper.BindPFlag(a4cAddressConfigName, rootCmd.PersistentFlags().Lookup(a4cAddressFlagName))
	viper.BindPFlag(a4cUserConfigName, rootCmd.PersistentFlags().Lookup(a4cUserFlagName))
	viper.BindPFlag(a4cPasswordConfigName, rootCmd.PersistentFlags().Lookup(a4cPasswordFlagName))
	viper.BindPFlag(a4cCAFileConfigName, rootCmd.PersistentFlags().Lookup(a4cCAFileFlagName))
	viper.BindPFlag(a4cSkipSecureConfigName, rootCmd.PersistentFlags().Lookup(a4cSkipSecureFlagName))

	// Vault config
	viper.BindPFlag(vaultAddressConfigName, rootCmd.PersistentFlags().Lookup(vaultAddressFlagName))
	viper.BindPFlag(vaultRoleIDConfigName, rootCmd.PersistentFlags().Lookup(vaultRoleIDFlagName))
	viper.BindPFlag(vaultSecretIDConfigName, rootCmd.PersistentFlags().Lookup(vaultSecretIDFlagName))
	viper.BindPFlag(vaultIsSecretWrappedConfigName, rootCmd.PersistentFlags().Lookup(vaultIsSecretWrappedFlagName))

	//Environment Variables
	viper.SetEnvPrefix("HWA") // HWA == HpcWaasApi
	viper.AutomaticEnv()      // read in environment variables that match

	// Global env
	viper.BindEnv(listenAddressFlagName)

	// Alien4Cloud env
	viper.BindEnv(a4cAddressConfigName)
	viper.BindEnv(a4cUserConfigName)
	viper.BindEnv(a4cPasswordConfigName)
	viper.BindEnv(a4cCAFileConfigName)
	viper.BindEnv(a4cSkipSecureConfigName)
	// Vault env
	viper.BindEnv(vaultAddressConfigName)
	viper.BindEnv(vaultRoleIDConfigName)
	viper.BindEnv(vaultSecretIDConfigName)
	viper.BindEnv(vaultIsSecretWrappedConfigName)

	// Global defaults
	viper.SetDefault(listenAddressFlagName, rest.DefaultListenAddress)

	// Alien4Cloud defaults
	viper.SetDefault(a4cAddressConfigName, a4c.DefaultAddress)
	viper.SetDefault(a4cUserConfigName, a4c.DefaultUser)
	viper.SetDefault(a4cPasswordConfigName, a4c.DefaultPassword)

	// Vault defaults
	viper.SetDefault(vaultAddressConfigName, vault.DefaultAddress)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".hpcwaas-api" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".hpcwaas-api")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
