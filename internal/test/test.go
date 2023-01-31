package test

import (
	"github.com/ceit-aut/S7CC01/internal/test/rabbit"

	"github.com/spf13/cobra"
)

// GetTestCommands
// returns the commands of test.
func GetTestCommands() *cobra.Command {
	testCmd := cobra.Command{
		Use: "test",
	}

	testCmd.AddCommand(rabbit.GetCommand())

	return &testCmd
}
