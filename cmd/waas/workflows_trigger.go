package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/eflows4hpc/hpcwaas-api/api"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func init() {
	var (
		// output string
		// maxlength int
		follow bool
		inputs []string
	)
	var workflowsTriggerCmd = &cobra.Command{
		Use:     "trigger <workflow-id>",
		Aliases: []string{"run", "r", "t"},
		Short:   "Triggers a given workflow",
		Long: `Triggers a given workflow:

The workflow should be referenced by its fully qualified ID.
		`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := api.GetClient(*clientConfig)
			if err != nil {
				errExit(err)
			}
			if follow {
				output = TextOutput
			}
			return triggerWorkflow(client, args[0], inputs, output, follow)
		},
	}
	workflowsTriggerCmd.Flags().StringSliceVarP(&inputs, "inputs", "i", nil, "Inputs for the workflow. The format is <key>=<value>. Can be specified multiple times.")
	workflowsTriggerCmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow workflow execution status, this implies --output=text.")
	workflowsCmd.AddCommand(workflowsTriggerCmd)
}

func triggerWorkflow(client api.HTTPClient, workflowID string, inputs []string, output string, follow bool) error {
	ctx := context.Background()
	workflowInputs := &api.WorkflowInputs{}
	for _, input := range inputs {
		parts := strings.SplitN(input, "=", 2)
		if workflowInputs.Inputs == nil {
			workflowInputs.Inputs = make(map[string]interface{})
		}
		workflowInputs.Inputs[parts[0]] = parts[1]
	}
	spinner, _ := pterm.DefaultSpinner.Start("Triggering workflow...")
	executionID, err := client.Workflows().Trigger(ctx, workflowID, workflowInputs)
	if err != nil {
		spinner.Fail("Failed to trigger workflow.")
		errExit(err)
	}
	if output == JSONOutput {
		fmt.Printf("%q", executionID)
	}
	spinner.Success("Workflow triggered with execution ID: ", executionID)

	if follow {
		return executionStatus(client, executionID, output, follow)
	}
	return nil
}
