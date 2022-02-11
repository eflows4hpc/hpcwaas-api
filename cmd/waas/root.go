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
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/eflows4hpc/hpcwaas-api/api"
)

func main() {
	Execute()
}

var cfgFile string
var output string
var clientConfig *api.Configuration

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "waas",
	Short: "A CLI for managing HPC Workflows as a Service",
	Long: `A Command Line Interface (CLI) for managing HPC Workflows as a Service.

This is the user side of the HPC Workflows as a Service (HPCWAAS) API.
It interacts with the API to manage the execution of HPC Workflows.`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		clientConfig = &api.Configuration{}
		err := viper.Unmarshal(clientConfig)
		if err != nil {
			return errors.Wrapf(err, "failed to parse config file: %q", viper.ConfigFileUsed())
		}
		if output != JSONOutput && output != TextOutput {
			return usageError{Mes: "Only \"json\" and \"text\" output formats are supported."}
		}

		if output == JSONOutput {
			pterm.DisableOutput()
		}
		if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) == 0 {
			// not a TTY
			pterm.DisableStyling()
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Print(err)
		}

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

const apiURLFlagName = "api_url"
const skipTLSVerifyFlagName = "skip_tls_verify"
const keyFileFlagName = "key_file"
const certFileFlagName = "cert_file"
const caFileFlagName = "ca_file"
const caPathFlagName = "ca_path"

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", DefaultDisplayOutput, "Output format either \"text\" or \"json\".")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hpcwaas-api.yaml)")
	rootCmd.PersistentFlags().StringP(apiURLFlagName, "", api.DefaultAPIAddress, "The default URL used to connect to the API")
	rootCmd.PersistentFlags().Bool(skipTLSVerifyFlagName, false, "Either or not to skip SSL certificates validation")
	rootCmd.PersistentFlags().String(caFileFlagName, "", "CA File to use to validate SSL certificates")
	rootCmd.PersistentFlags().String(caPathFlagName, "", "directory path to CA Files to use to validate SSL certificates")
	rootCmd.PersistentFlags().String(keyFileFlagName, "", "TLS key file to use to authenticate to the API")
	rootCmd.PersistentFlags().String(certFileFlagName, "", "TLS cert file to use to authenticate to the API")

	// Global flags/config binding
	viper.BindPFlag(apiURLFlagName, rootCmd.PersistentFlags().Lookup(apiURLFlagName))
	viper.BindPFlag(skipTLSVerifyFlagName, rootCmd.PersistentFlags().Lookup(skipTLSVerifyFlagName))
	viper.BindPFlag(caFileFlagName, rootCmd.PersistentFlags().Lookup(caFileFlagName))
	viper.BindPFlag(caPathFlagName, rootCmd.PersistentFlags().Lookup(caPathFlagName))
	viper.BindPFlag(keyFileFlagName, rootCmd.PersistentFlags().Lookup(keyFileFlagName))
	viper.BindPFlag(certFileFlagName, rootCmd.PersistentFlags().Lookup(certFileFlagName))

	//Environment Variables
	viper.SetEnvPrefix("HW") // HW == HpcWaas
	viper.AutomaticEnv()     // read in environment variables that match

	// Global env
	viper.BindEnv(apiURLFlagName)
	viper.BindEnv(skipTLSVerifyFlagName)
	viper.BindEnv(caFileFlagName)
	viper.BindEnv(caPathFlagName)
	viper.BindEnv(keyFileFlagName)
	viper.BindEnv(certFileFlagName)

	// Global defaults
	viper.SetDefault(apiURLFlagName, api.DefaultAPIAddress)

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

		// Search config in home directory with name ".waas" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".waas")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
