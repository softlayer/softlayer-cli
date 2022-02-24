package ipsec

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	UI           terminal.UI
	IPSECManager managers.IPSECManager
}

func NewListCommand(ui terminal.UI, ipsecManager managers.IPSECManager) (cmd *ListCommand) {
	return &ListCommand{
		UI:           ui,
		IPSECManager: ipsecManager,
	}
}

func (cmd *ListCommand) Run(c *cli.Context) error {
	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}
	contexts, err := cmd.IPSECManager.GetTunnelContexts(c.Int("order"), "")
	if err != nil {
		return cli.NewExitError(T("Failed to get IPSec on your account.")+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, contexts)
	}

	if len(contexts) == 0 {
		cmd.UI.Print(T("No IPSec was found."))
	} else {
		table := cmd.UI.Table([]string{T("ID"), T("Name"), T("FriendlyName"), T("Internal Peer IP Address"), T("Customer Peer IP Address"), T("Created")})
		for _, c := range contexts {
			table.Add(utils.FormatIntPointer(c.Id),
				utils.FormatStringPointer(c.Name),
				utils.FormatStringPointer(c.FriendlyName),
				utils.FormatStringPointer(c.InternalPeerIpAddress),
				utils.FormatStringPointer(c.CustomerPeerIpAddress),
				utils.FormatSLTimePointer(c.CreateDate))
		}
		table.Print()
	}
	return nil
}

func IpsecListMetaData() cli.Command {
	return cli.Command{
		Category:    "ipsec",
		Name:        "list",
		Description: T("List IPSec VPN tunnel contexts"),
		Usage:       "${COMMAND_NAME} sl ipsec list [OPTIONS]",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "order",
				Usage: T("Filter by ID of the order that purchased the IPSec"),
			},
			metadata.OutputFlag(),
		},
	}
}
