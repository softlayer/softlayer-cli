package globalip

import (
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
}

func NewListCommand(ui terminal.UI, networkManager managers.NetworkManager) (cmd *ListCommand) {
	return &ListCommand{
		UI:             ui,
		NetworkManager: networkManager,
	}
}

func (cmd *ListCommand) Run(c *cli.Context) error {
	if c.IsSet("v4") && c.IsSet("v6") {
		return errors.NewInvalidUsageError(T("[--v4] is not allowed with [--v6]."))
	}

	version := 0
	if c.IsSet("v4") {
		version = 4
	}
	if c.IsSet("v6") {
		version = 6
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	ips, err := cmd.NetworkManager.ListGlobalIPs(version, c.Int("order"))
	if err != nil {
		return cli.NewExitError(T("Failed to list global IPs on your account.\n")+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, ips)
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
	table.Print()
	return nil
}

func GlobalIpListMetaData() cli.Command {
	return cli.Command{
		Category:    "globalip",
		Name:        "list",
		Description: T("List all global IPs on your account"),
		Usage: T(`${COMMAND_NAME} sl globalip list [OPTIONS]

EXAMPLE:
    ${COMMAND_NAME} sl globalip list --v4 
	This command lists all IPv4 addresses on the current account.`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "v4",
				Usage: T("Display IPv4 IPs only"),
			},
			cli.BoolFlag{
				Name:  "v6",
				Usage: T("Display IPv6 IPs only"),
			},
			cli.IntFlag{
				Name:  "order",
				Usage: T("Filter by the ID of order that purchased this IP address"),
			},
			metadata.OutputFlag(),
		},
	}
}
