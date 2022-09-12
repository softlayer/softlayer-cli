package securitygroup

import (
	"bytes"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
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
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
}

func NewDetailCommand(sl *metadata.SoftlayerCommand) (cmd *DetailCommand) {
	thisCmd := &DetailCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "detail " + T("SECURITYGROUP_ID"),
		Short: T("Get details about a security group"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *DetailCommand) Run(args []string) error {
	groupID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Security group ID")
	}

	outputFormat := cmd.GetOutputFlag()

	group, err := cmd.NetworkManager.GetSecurityGroup(groupID, "")
	if err != nil {
		return errors.NewAPIError(T("Failed to get security group {{.ID}}.\n", map[string]interface{}{"ID": groupID}), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, group)
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("ID"), utils.FormatIntPointer(group.Id))
	table.Add(T("Name"), utils.FormatStringPointer(group.Name))
	table.Add(T("Description"), utils.FormatStringPointer(group.Description))
	buf := new(bytes.Buffer)
	ruleTable := terminal.NewTable(buf, []string{T("ID"), T("Remote IP"), T("Remote Group ID"), T("Direction"), T("Ether Type"), T("Port Range Min"), T("Port Range Max"), T("Protocol")})
	for _, rule := range group.Rules {
		ruleTable.Add(utils.FormatIntPointer(rule.Id),
			utils.FormatStringPointer(rule.RemoteIp),
			utils.FormatIntPointer(rule.RemoteGroupId),
			utils.FormatStringPointer(rule.Direction),
			utils.FormatStringPointer(rule.Ethertype),
			utils.FormatIntPointer(rule.PortRangeMin),
			utils.FormatIntPointer(rule.PortRangeMax),
			utils.FormatStringPointer(rule.Protocol),
		)
	}
	ruleTable.Print()
	table.Add(T("Rules"), buf.String())

	buf = new(bytes.Buffer)
	serverTable := terminal.NewTable(buf, []string{T("ID"), T("Hostname"), T("Interface"), T("IP address")})
	for _, component := range group.NetworkComponentBindings {
		if component.NetworkComponent != nil && component.NetworkComponent.Guest != nil {
			var networkInterface, ipaddress string
			vsi := *component.NetworkComponent.Guest

			if (component.NetworkComponent.Port != nil && *component.NetworkComponent.Port == 0) || component.NetworkComponent.Port == nil {
				networkInterface = T("private")
				if vsi.PrimaryBackendIpAddress != nil {
					ipaddress = *vsi.PrimaryBackendIpAddress
				}
			} else {
				networkInterface = T("public")
				if vsi.PrimaryIpAddress != nil {
					ipaddress = *vsi.PrimaryIpAddress
				}
			}
			serverTable.Add(utils.FormatIntPointer(vsi.Id),
				utils.FormatStringPointer(vsi.Hostname),
				networkInterface,
				ipaddress,
			)
		}
	}
	serverTable.Print()
	table.Add(T("Servers"), buf.String())
	table.Print()
	return nil
}
