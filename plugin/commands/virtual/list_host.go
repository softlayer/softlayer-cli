package virtual

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
)

type ListHostCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewListHostCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *ListHostCommand) {
	return &ListHostCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *ListHostCommand) Run(c *cli.Context) error {

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	hosts, err := cmd.VirtualServerManager.ListDedicatedHost(c.String("name"), c.String("datacenter"), c.String("owner"), c.Int("order"))
	if err != nil {
		return cli.NewExitError(T("Failed to list dedicated hosts on your account.\n")+err.Error(), 2)
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
