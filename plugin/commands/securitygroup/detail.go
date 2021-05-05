package securitygroup

import (
	"bytes"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type DetailCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
}

func NewDetailCommand(ui terminal.UI, networkManager managers.NetworkManager) (cmd *DetailCommand) {
	return &DetailCommand{
		UI:             ui,
		NetworkManager: networkManager,
	}
}

func (cmd *DetailCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	groupID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Security group ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	group, err := cmd.NetworkManager.GetSecurityGroup(groupID, "")
	if err != nil {
		return cli.NewExitError(T("Failed to get security group {{.ID}}.\n", map[string]interface{}{"ID": groupID})+err.Error(), 2)
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
