package virtual

import (
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dedicatedhost"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

func NewListHostCommand(sl *metadata.SoftlayerCommand) (cmd *dedicatedhost.ListCommand) {
	thisCmd := &dedicatedhost.ListCommand{
		SoftlayerCommand:     sl,
		DedicatedHostManager: managers.NewDedicatedhostManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "host-list",
		Short: T("List dedicated hosts on your account"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	cobraCmd.Flags().StringVarP(&thisCmd.Name, "name", "n", "", T("Filter by name of the dedicated host"))
	cobraCmd.Flags().StringVarP(&thisCmd.Datacenter, "datacenter", "d", "", T("Filter by datacenter of the dedicated host"))
	cobraCmd.Flags().StringVar(&thisCmd.Owner, "owner", "", T("Filter by owner of the dedicated host"))
	cobraCmd.Flags().IntVar(&thisCmd.Order, "order", 0, T("Filter by ID of the order which purchased this dedicated host"))
	cobraCmd.Flags().StringVar(&thisCmd.SortBy, "sortby", "", T("Column to sort by (Id, Name, Datacenter, Router, Cpu, Memory, Disk, Guests)[default: Id]"))
	return thisCmd
}
