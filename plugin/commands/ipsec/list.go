package ipsec

import (
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	*metadata.SoftlayerCommand
	IPSECManager managers.IPSECManager
	Command      *cobra.Command
	Order        int
}

func NewListCommand(sl *metadata.SoftlayerCommand) (cmd *ListCommand) {
	thisCmd := &ListCommand{
		SoftlayerCommand: sl,
		IPSECManager:     managers.NewIPSECManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "list",
		Short: T("List IPSec VPN tunnel contexts"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().IntVar(&thisCmd.Order, "order", 0, T("Filter by ID of the order that purchased the IPSec"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ListCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()
	contexts, err := cmd.IPSECManager.GetTunnelContexts(cmd.Order, "")
	if err != nil {
		return errors.NewAPIError(T("Failed to get IPSec on your account."), err.Error(), 2)
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
