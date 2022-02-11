package main

import (
	"os"

	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

const (
	JSONOutput              = "json"
	TextOutput              = "text"
	TOMLOutput              = "toml"
	DefaultDisplayPageLimit = 50
	DefaultDisplayMaxlength = 80
	DefaultDisplayOutput    = "text"
)

type usageError struct {
	Mes string
}

func (e usageError) Error() string {
	return e.Mes
}

// isUsageError checks if an error is an UsageError
func isUsageError(err error) bool {
	_, ok := errors.Cause(err).(usageError)
	return ok
}

// ErrExit allows to exit on error with exit code 1 after printing error message
func errExit(msg interface{}) {
	pterm.DefaultSection.Println("Errors:")
	pterm.Error.Println(msg)
	os.Exit(1)
}

func handleUsageError(cmd *cobra.Command, err error) error {
	if !isUsageError(err) {
		cmd.SilenceUsage = true
	}
	return err
}
