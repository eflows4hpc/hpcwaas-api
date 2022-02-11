package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/eflows4hpc/hpcwaas-api/api"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func init() {
	var follow bool
	executionsStatusCmd := &cobra.Command{
		Use:     "status <execution-id>",
		Aliases: []string{"stat", "get", "g", "s"},
		Short:   "Get status of an execution",
		Long: `Get status of an execution given its ID:


		`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.GetClient(*clientConfig)
			if err != nil {
				errExit(err)
			}
			err = executionStatus(client, args[0], output, follow)
			return handleUsageError(cmd, err)
		},
	}

	executionsStatusCmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow execution status, this implies --output=text.")

	executionsCmd.AddCommand(executionsStatusCmd)
}

func executionStatus(httpClient api.HTTPClient, executionsID, output string, follow bool) error {

	if follow {
		output = TextOutput
	}

	if output != JSONOutput && output != TextOutput {
		return usageError{Mes: "Only \"json\" and \"text\" output formats are supported."}
	}
	var executionResult api.Execution
	var err error

	ctx := context.Background()
	executionResult, err = httpClient.Executions().Status(ctx, executionsID)
	if err != nil {
		return err
	}

	if output == JSONOutput {
		data, err := json.Marshal(executionResult)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	return displayExecution(httpClient, executionResult, follow)
}

func displayExecution(httpClient api.HTTPClient, execution api.Execution, follow bool) error {
	prefix := pterm.PrefixPrinter{
		MessageStyle: &pterm.ThemeDefault.InfoMessageStyle,
		Prefix: pterm.Prefix{
			Style: &pterm.ThemeDefault.InfoPrefixStyle,
			Text:  "Execution",
		},
	}
	prefix.Printf("ID: %s\n", execution.ID)
	prefix.Printf("Status: %s", execution.Status)

	if follow {
		// Clear line
		pterm.Printo(strings.Repeat(" ", pterm.GetTerminalWidth()))
		return followExecution(httpClient, execution)
	}
	pterm.Println()
	return nil
}

func followExecution(httpClient api.HTTPClient, execution api.Execution) error {
	spinner, _ := pterm.DefaultSpinner.Start("Status: ", execution.Status)
	for strings.HasSuffix(strings.ToLower(execution.Status), "ing") {
		time.Sleep(5 * time.Second)
		var err error

		ctx := context.Background()
		execution, err = httpClient.Executions().Status(ctx, execution.ID)
		if err != nil {
			spinner.Fail("Fail to get execution status")
			return err
		}
		spinner.UpdateText(fmt.Sprintf("Status: %s", execution.Status))
	}

	switch strings.ToLower(execution.Status) {
	case "failed", "cancelled":
		spinner.Fail(fmt.Sprintf("Status: %s", execution.Status))
	default:
		spinner.Success(fmt.Sprintf("Status: %s", execution.Status))
	}

	return nil
}
