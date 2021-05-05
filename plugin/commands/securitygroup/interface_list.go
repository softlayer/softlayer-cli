package securitygroup

import (
	"sort"
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

type InterfaceListCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
}

func NewInterfaceListCommand(ui terminal.UI, networkManager managers.NetworkManager) (cmd *InterfaceListCommand) {
	return &InterfaceListCommand{
		UI:             ui,
		NetworkManager: networkManager,
	}
}

func (cmd *InterfaceListCommand) Run(c *cli.Context) error {
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

	sortColumns := []string{"id", "virtualServerId", "hostname"}
	sortby := c.String("sortby")
	if sortby != "" && utils.StringInSlice(sortby, sortColumns) == -1 {
		return errors.NewInvalidUsageError(T("--sortby {{.Column}} is not supported.", map[string]interface{}{"Column": sortby}))
	}
	securityGroup, err := cmd.NetworkManager.GetSecurityGroup(groupID, "")
	if err != nil {
		return cli.NewExitError(T("Failed to get security group {{.GroupID}}.\n", map[string]interface{}{"GroupID": groupID})+err.Error(), 2)
	}

	bindings := securityGroup.NetworkComponentBindings
	if sortby == "" || sortby == "id" {
		sort.Sort(utils.InterfaceByInterfaceId(bindings))
	} else if sortby == "virtualServerId" {
		sort.Sort(utils.InterfaceByVSId(bindings))
	} else if sortby == "hostname" {
		sort.Sort(utils.InterfaceByVSHost(bindings))
	} else {
		return errors.NewInvalidUsageError(T("--sortby {{.Column}} is not supported.", map[string]interface{}{"Column": sortby}))
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, bindings)
	}

	if len(bindings) == 0 {
		cmd.UI.Print(T("No interfaces are binded to security group {{.GroupID}}.", map[string]interface{}{"GroupID": groupID}))
		return nil
	}

	serverTable := cmd.UI.Table([]string{T("ID"), T("Server ID"), T("Hostname"), T("Interface"), T("IP address")})
	for _, component := range bindings {
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
			serverTable.Add(utils.FormatIntPointer(component.NetworkComponent.Id),
				utils.FormatIntPointer(vsi.Id),
				utils.FormatStringPointer(vsi.Hostname),
				networkInterface,
				ipaddress,
			)
		}
	}
	serverTable.Print()
	return nil
}
