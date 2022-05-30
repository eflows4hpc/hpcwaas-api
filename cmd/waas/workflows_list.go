package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eflows4hpc/hpcwaas-api/api"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var workflowsListCmd *cobra.Command

func init() {
	workflowsListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Short:   "list available workflows",
		Long: `List available workflows:

This command lists all workflows that are tagged as usable by the waas service.
		`,
		Args: cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.GetClient(*clientConfig)
			if err != nil {
				errExit(err)
			}
			err = listWorkflows(client, output)
			return handleUsageError(cmd, err)
		},
	}

	workflowsCmd.AddCommand(workflowsListCmd)
}

func listWorkflows(httpClient api.HTTPClient, output string) error {
	ctx := context.Background()
	workflowsResult, err := httpClient.Workflows().List(ctx)
	if err != nil {
		return err
	}
	if workflowsResult.Workflows == nil || len(workflowsResult.Workflows) == 0 {
		fmt.Println("No workflows found.")
		return nil
	}
	if output == JSONOutput {
		data, err := json.Marshal(workflowsResult.Workflows)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	displayWorkflowsTable(workflowsResult.Workflows)
	return nil
}

func displayWorkflowsTable(workflows []api.Workflow) {
	data := pterm.TableData{
		{"Workflow ID", "Workflow Name", "Application", "Environment ID", "Environment Name"},
	}
	for _, workflow := range workflows {
		data = append(data, []string{workflow.ID, workflow.Name, workflow.ApplicationID, workflow.EnvironmentID, workflow.EnvironmentName})
	}
	pterm.DefaultTable.WithHasHeader().WithHeaderRowSeparator("*").WithRowSeparator("-").WithData(data).WithBoxed().Render()
}
