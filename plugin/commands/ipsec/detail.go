package ipsec

import (
	"bytes"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type DetailCommand struct {
	*metadata.SoftlayerCommand
	IPSECManager managers.IPSECManager
	Command      *cobra.Command
	Include      []string
}

func NewDetailCommand(sl *metadata.SoftlayerCommand) (cmd *DetailCommand) {
	thisCmd := &DetailCommand{
		SoftlayerCommand: sl,
		IPSECManager:     managers.NewIPSECManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "detail " + T("CONTEXT_ID"),
		Short: T("List IPSec VPN tunnel context details"),
		Long: T(`${COMMAND_NAME} sl ipsec detail CONTEXT_ID [OPTIONS]

  List IPSEC VPN tunnel context details.

  Additional resources can be joined using multiple instances of the include
  option, for which the following choices are available.

  at: address translations
  is: internal subnets
  rs: remote subnets
  sr: statically routed subnets
  ss: service subnets`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringSliceVarP(&thisCmd.Include, "include", "i", []string{}, T("Include extra resources. Options are: at,is,rs,sr,ss"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *DetailCommand) Run(args []string) error {
	args0 := args[0]
	contextId, err := strconv.Atoi(args0)
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Context ID")
	}

	outputFormat := cmd.GetOutputFlag()

	at, is, rs, sr, ss := false, false, false, false, false
	includes := cmd.Include
	if i := utils.StringInSlice("at", includes); i >= 0 {
		at = true
	}
	if i := utils.StringInSlice("is", includes); i >= 0 {
		is = true
	}
	if i := utils.StringInSlice("rs", includes); i >= 0 {
		rs = true
	}
	if i := utils.StringInSlice("sr", includes); i >= 0 {
		sr = true
	}
	if i := utils.StringInSlice("ss", includes); i >= 0 {
		ss = true
	}
	mask := GetTunnelContextMask(at, is, rs, sr, ss)
	context, err := cmd.IPSECManager.GetTunnelContext(contextId, mask)
	if err != nil {
		return errors.NewAPIError(T("Failed to get IPSec with ID {{.ID}}.\n", map[string]interface{}{"ID": contextId}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, context)
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table = printTunnelContext(table, &context)
	if at {
		table = printAddressTransaltion(table, T("Address Translations"), context.AddressTranslations)
	}
	if is {
		table = printSubnetTable(table, T("Internal Subnets"), context.InternalSubnets)
	}
	if rs {
		table = printCustomerSubnetTable(table, T("Remote Subnets"), context.CustomerSubnets)
	}
	if sr {
		table = printSubnetTable(table, T("Static Subnets"), context.StaticRouteSubnets)
	}
	if ss {
		table = printSubnetTable(table, T("Service Subnets"), context.ServiceSubnets)
	}
	table.Print()
	return nil
}

// Yields a mask for a tunnel context
// All exposed properties on the tunnel context service are included inthe constructed mask. Additional joins may be requested.
// addressTranslations: Whether to join the context's address translation entries.
// internalSubnets: Whether to join the context's internal subnet associations.
// remoteSubnets: Whether to join the context's remote subnet associations.
// staticSubnets: Whether to join the context's statically routed subnet associations.
// serviceSubnets: Whether to join the SoftLayer service network subnets.
func GetTunnelContextMask(addressTranslation, internalSubnets, remoteSubnets, statusSubnets, serviceSubnets bool) string {
	mask := "id,accountId,advancedConfigurationFlag,createDate,customerPeerIpAddress,modifyDate,name,friendlyName,internalPeerIpAddress" +
		",phaseOneAuthentication,phaseOneDiffieHellmanGroup,phaseOneEncryption,phaseOneKeylife" +
		",phaseTwoAuthentication,phaseTwoDiffieHellmanGroup,phaseTwoEncryption,phaseTwoKeylife" +
		",phaseTwoPerfectForwardSecrecy,presharedKey"
	if addressTranslation {
		mask = mask + ",addressTranslations[internalIpAddressRecord[ipAddress],customerIpAddressRecord[ipAddress]]"
	}
	if internalSubnets {
		mask = mask + ",internalSubnets"
	}
	if remoteSubnets {
		mask = mask + ",customerSubnets"
	}
	if statusSubnets {
		mask = mask + ",staticRouteSubnets"
	}
	if serviceSubnets {
		mask = mask + ",serviceSubnets"
	}
	return mask
}

func printTunnelContext(table terminal.Table, context *datatypes.Network_Tunnel_Module_Context) terminal.Table {
	table.Add(T("ID"), utils.FormatIntPointer(context.Id))
	table.Add(T("Name"), utils.FormatStringPointer(context.Name))
	table.Add(T("Friendly name"), utils.FormatStringPointer(context.FriendlyName))
	table.Add(T("Internal peer IP address"), utils.FormatStringPointer(context.InternalPeerIpAddress))
	table.Add(T("Remote peer IP address"), utils.FormatStringPointer(context.CustomerPeerIpAddress))
	table.Add(T("Advanced configuration flag"), utils.FormatIntPointer(context.AdvancedConfigurationFlag))
	table.Add(T("Preshared key"), utils.FormatStringPointer(context.PresharedKey))
	table.Add(T("Phase 1 authentication"), utils.FormatStringPointer(context.PhaseOneAuthentication))
	table.Add(T("Phase 1 diffie hellman group"), utils.FormatIntPointer(context.PhaseOneDiffieHellmanGroup))
	table.Add(T("Phase 1 encryption"), utils.FormatStringPointer(context.PhaseOneEncryption))
	table.Add(T("Phase 1 key life"), utils.FormatIntPointer(context.PhaseOneKeylife))
	table.Add(T("Phase 2 authentication"), utils.FormatStringPointer(context.PhaseTwoAuthentication))
	table.Add(T("Phase 2 diffie hellman group"), utils.FormatIntPointer(context.PhaseTwoDiffieHellmanGroup))
	table.Add(T("Phase 2 encryption"), utils.FormatStringPointer(context.PhaseTwoEncryption))
	table.Add(T("Phase 2 key life"), utils.FormatIntPointer(context.PhaseTwoKeylife))
	table.Add(T("Phase 2 perfect forward secrecy"), utils.FormatIntPointer(context.PhaseTwoPerfectForwardSecrecy))
	table.Add(T("Created"), utils.FormatSLTimePointer(context.CreateDate))
	table.Add(T("Modified"), utils.FormatSLTimePointer(context.ModifyDate))
	return table
}

func printSubnetTable(table terminal.Table, header string, subnets []datatypes.Network_Subnet) terminal.Table {
	if len(subnets) == 0 {
		table.Add(header, T("None"))
	} else {
		buf := new(bytes.Buffer)
		snTable := terminal.NewTable(buf, []string{T("ID"), T("Network identifier"), T("CIDR"), T("Note")})
		for _, sn := range subnets {
			snTable.Add(utils.FormatIntPointer(sn.Id),
				utils.FormatStringPointer(sn.NetworkIdentifier),
				utils.FormatIntPointer(sn.Cidr),
				utils.FormatStringPointer(sn.Note))
		}
		snTable.Print()
		table.Add(header, buf.String())
	}
	return table
}

func printCustomerSubnetTable(table terminal.Table, header string, subnets []datatypes.Network_Customer_Subnet) terminal.Table {
	if len(subnets) == 0 {
		table.Add(header, T("None"))
	} else {
		buf := new(bytes.Buffer)
		snTable := terminal.NewTable(buf, []string{T("ID"), T("Network identifier"), T("CIDR"), T("Note")})
		for _, sn := range subnets {
			snTable.Add(utils.FormatIntPointer(sn.Id),
				utils.FormatStringPointer(sn.NetworkIdentifier),
				utils.FormatIntPointer(sn.Cidr),
				"")
		}
		snTable.Print()
		table.Add(header, buf.String())
	}
	return table
}

func printAddressTransaltion(table terminal.Table, header string, translations []datatypes.Network_Tunnel_Module_Context_Address_Translation) terminal.Table {
	if len(translations) == 0 {
		table.Add(header, T("None"))
	} else {
		buf := new(bytes.Buffer)
		atTable := terminal.NewTable(buf, []string{T("ID"), T("Static IP address"), T("Static IP address ID"), T("Remote IP address"), T("Remote IP address ID"), T("Note")})
		for _, at := range translations {
			atTable.Add(utils.FormatIntPointer(at.Id),
				utils.FormatStringPointer(at.InternalIpAddressRecord.IpAddress),
				utils.FormatIntPointer(at.InternalIpAddressId),
				utils.FormatStringPointer(at.CustomerIpAddressRecord.IpAddress),
				utils.FormatIntPointer(at.CustomerIpAddressId),
				utils.FormatStringPointer(at.Notes))
		}
		atTable.Print()
		table.Add(header, buf.String())
	}

	return table
}
