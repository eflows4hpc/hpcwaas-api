package main

import (
	"context"
	"fmt"

	"github.com/eflows4hpc/hpcwaas-api/api"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func init() {
	executionsCancelCmd := &cobra.Command{
		Use:     "cancel <execution-id>...",
		Aliases: []string{"delete", "del", "d", "c"},
		Short:   "Cancel one or more executions",
		Long: `Cancel one or more executions given their IDs


		`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.GetClient(*clientConfig)
			if err != nil {
				errExit(err)
			}
			err = executionsCancel(client, args, output)
			return handleUsageError(cmd, err)
		},
	}

	executionsCmd.AddCommand(executionsCancelCmd)
}

func executionsCancel(client api.HTTPClient, executionsIDs []string, output string) error {
	errs := &api.Errors{}

	progress, _ := pterm.DefaultProgressbar.WithTotal(len(executionsIDs)).WithTitle("Cancelling executions").Start()
	progress.RemoveWhenDone = true
	for _, executionID := range executionsIDs {
		progress.UpdateTitle(fmt.Sprintf("Cancelling execution %q", executionID))
		progress.Increment()

		err := client.Executions().Cancel(context.Background(), executionID)
		if err != nil {
			if e, ok := err.(*api.Errors); ok {
				errs.Errors = append(errs.Errors, e.Errors...)
			} else {
				errs.Errors = append(errs.Errors, &api.Error{ID: "internal error", Title: "Internal Error", Detail: err.Error()})
			}
			pterm.Error.Printf("Failed to cancel execution %q...\n", executionID)
			continue
		}
		pterm.Success.Printf("Cancelled execution %q\n", executionID)
	}
	progress.Stop()

	if len(errs.Errors) > 0 {
		errExit(errs)
	}

	pterm.Success.Printf("Successfully cancelled %d executions...\n", len(executionsIDs))

	return nil
}
