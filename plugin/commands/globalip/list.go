package globalip

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
	V4             bool
	V6             bool
	Order          int
}

func NewListCommand(sl *metadata.SoftlayerCommand) *ListCommand {
	thisCmd := &ListCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "list",
		Short: T("List all global IPs on your account."),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().BoolVar(&thisCmd.V4, "v4", false, T("Display IPv4 IPs only"))
	cobraCmd.Flags().BoolVar(&thisCmd.V6, "v6", false, T("Display IPv6 IPs only"))
	cobraCmd.Flags().IntVar(&thisCmd.Order, "order", 0, T("Filter by the ID of order that purchased this IP address"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ListCommand) Run(args []string) error {
	if cmd.V4 && cmd.V6 {
		return errors.NewInvalidUsageError(T("[--v4] is not allowed with [--v6]."))
	}

	version := 0
	if cmd.V4 {
		version = 4
	}
	if cmd.V6 {
		version = 6
	}

	outputFormat := cmd.GetOutputFlag()

	ips, err := cmd.NetworkManager.ListGlobalIPs(version, cmd.Order)
	if err != nil {
		return errors.NewAPIError(T("Failed to list global IPs on your account."), err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("ID"), T("ip"), T("assigned"), T("target")})
	for _, ip := range ips {
		ipAddress := ""
		assigned := T("No")
		target := T("None")
		if ip.IpAddress != nil {
			ipAddress = utils.FormatStringPointer(ip.IpAddress.IpAddress)
		}
		if ip.DestinationIpAddress != nil {
			dest := ip.DestinationIpAddress
			assigned = T("Yes")
			target = utils.FormatStringPointer(ip.DestinationIpAddress.IpAddress)
			if vs := dest.VirtualGuest; vs != nil {
				target += fmt.Sprintf("(%s)", utils.FormatStringPointer(vs.FullyQualifiedDomainName))
			} else if hw := dest.Hardware; hw != nil {
				target += fmt.Sprintf("(%s)", utils.FormatStringPointer(hw.FullyQualifiedDomainName))
			}
		}
		table.Add(utils.FormatIntPointer(ip.Id), ipAddress, assigned, target)
	}
	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}
