package ticket

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   "ticket",
		Short: T("Classic infrastructure Manage Tickets"),
		RunE:  nil,
	}
	cobraCmd.AddCommand(NewAttachDeviceTicketCommand(sl).Command)
	cobraCmd.AddCommand(NewCreateStandardTicketCommand(sl).Command)
	cobraCmd.AddCommand(NewDetachDeviceTicketCommand(sl).Command)
	cobraCmd.AddCommand(NewDetailTicketCommand(sl).Command)
	cobraCmd.AddCommand(NewListTicketCommand(sl).Command)
	cobraCmd.AddCommand(NewSubjectsTicketCommand(sl).Command)
	cobraCmd.AddCommand(NewSummaryTicketCommand(sl).Command)
	cobraCmd.AddCommand(NewUpdateTicketCommand(sl).Command)
	cobraCmd.AddCommand(NewUploadFileTicketCommand(sl).Command)
	return cobraCmd
}

func TicketNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "ticket",
		Description: T("Classic infrastructure Manage Tickets"),
	}
}
