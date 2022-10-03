package testhelpers

import (
	"github.com/spf13/cobra"
)



func RunCobraCommand(cmd *cobra.Command, args ...string) error {
	// If we do cmd.SetArgs(args) with no args, Cobra will try to read them from the actual command line
	// which breaks unit tests when using -ginkgo.focus (or other) flags.
	if len(args) == 0 {
		cmd.SetArgs([]string{})	
	} else {
		cmd.SetArgs(args)
	}
	
	
	_, err := cmd.ExecuteC()
	return err
}
