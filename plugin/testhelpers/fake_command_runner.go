package testhelpers

import (
	"flag"
	"fmt"
	"log"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	faketerminal "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"

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



type CMD struct {
	UI terminal.UI
}

func NewCommand(ui terminal.UI) *CMD {
	return &CMD{
		UI: ui,
	}
}
