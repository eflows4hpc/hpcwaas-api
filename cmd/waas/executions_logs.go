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
	var follow bool
	var fromIndex int
	var debug bool
	executionsLogsCmd := &cobra.Command{
		Use:     "logs <execution-id>",
		Aliases: []string{"log", "l"},
		Short:   "Get logs of an execution",
		Long:    `Get logs of an execution given its ID:`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.GetClient(*clientConfig)
			if err != nil {
				errExit(err)
			}
			err = executionLogs(client, args[0], fromIndex, output, follow, debug)
			return handleUsageError(cmd, err)
		},
	}

	executionsLogsCmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow execution logs, this implies --output=text.")
	executionsLogsCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Show debug logs")
	executionsLogsCmd.Flags().IntVarP(&fromIndex, "from", "", 0, "Starting index for pagination")

	executionsCmd.AddCommand(executionsLogsCmd)
}

func executionLogs(httpClient api.HTTPClient, executionsID string, fromIndex int, output string, follow, debug bool) error {

	if follow {
		output = TextOutput
	}

	if output != JSONOutput && output != TextOutput {
		return usageError{Mes: "Only \"json\" and \"text\" output formats are supported."}
	}

	ctx := context.Background()
	size := 50
	logsOpts := &api.LogsRequestOpts{
		FromIndex: &fromIndex,
		Size:      &size,
		Levels:    api.SetLogLevels(api.INFO, api.WARN, api.ERROR),
	}
	if debug {
		logsOpts.Levels = api.SetLogLevels(logsOpts.Levels, api.DEBUG)
	}
	executionLogsResult, err := httpClient.Executions().Logs(ctx, executionsID, logsOpts)
	if err != nil {
		return err
	}

	if output == JSONOutput {
		data, err := json.Marshal(executionLogsResult)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	displayExecutionLogs(executionLogsResult)
	if !follow {
		return nil
	}
	for {
		*logsOpts.FromIndex = *logsOpts.FromIndex + len(executionLogsResult.Logs)
		executionLogsResult, err = httpClient.Executions().Logs(ctx, executionsID, logsOpts)
		if err != nil {
			return err
		}
		displayExecutionLogs(executionLogsResult)
	}
}

func getLeveledPrefix(level string) pterm.PrefixPrinter {
	var p pterm.PrefixPrinter
	if level == "DEBUG" {
		p = pterm.Debug
		p.Debugger = false
		p.Prefix.Text = "DEBUG"
	} else if level == "ERROR" {
		p = pterm.Error
		p.ShowLineNumber = false
		p.Prefix.Text = "ERROR"
	} else if level == "WARN" {
		p = pterm.Warning
		p.Prefix.Text = "WARN "
	} else {
		p = pterm.Info
		p.Prefix.Text = "INFO "
	}
	return p
}

func displayExecutionLogs(executionLogs api.ExecutionLogs) {
	for _, log := range executionLogs.Logs {
		prefix := getLeveledPrefix(log.Level)
		prefix.Printf("%s %s\n", pterm.Gray("["+log.Timestamp.String()+"]"), log.Content)
	}
}
