package virtual

import (
	"github.com/spf13/cobra"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListHostCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
	Name                 string
	Datacenter           string
	Owner                string
	Order                int
}

func NewListHostCommand(sl *metadata.SoftlayerCommand) (cmd *ListHostCommand) {
	thisCmd := &ListHostCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
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
	return thisCmd
}

func (cmd *ListHostCommand) Run(args []string) error {

	outputFormat := cmd.GetOutputFlag()

	hosts, err := cmd.VirtualServerManager.ListDedicatedHost(cmd.Name, cmd.Datacenter, cmd.Owner, cmd.Order)
	if err != nil {
		return slErrors.NewAPIError(T("Failed to list dedicated hosts on your account.\n"), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, hosts)
	}

	if len(hosts) == 0 {
		cmd.UI.Print(T("No dedicated hosts are found."))
		return nil
	}

	table := cmd.UI.Table([]string{T("id"), T("name"), T("datacenter"), T("router"), T("cpu (allocated/total)"), T("memory (allocated/total)"), T("disk (allocated/total)"), T("guests")})
	for _, host := range hosts {
		table.Add(
			utils.FormatIntPointer(host.Id),
			utils.FormatStringPointer(host.Name),
			utils.FormatStringPointer(host.Datacenter.Name),
			utils.FormatStringPointer(host.BackendRouter.Hostname),
			utils.FormatIntPointer(host.AllocationStatus.CpuAllocated)+"/"+utils.FormatIntPointer(host.AllocationStatus.CpuCount),
			utils.FormatIntPointer(host.AllocationStatus.MemoryAllocated)+"/"+utils.FormatIntPointer(host.AllocationStatus.MemoryCapacity),
			utils.FormatIntPointer(host.AllocationStatus.DiskAllocated)+"/"+utils.FormatIntPointer(host.AllocationStatus.DiskCapacity),
			utils.FormatUIntPointer(host.GuestCount),
		)
	}
	table.Print()
	return nil
}
