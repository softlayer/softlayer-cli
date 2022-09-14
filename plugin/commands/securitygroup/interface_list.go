package securitygroup

import (
	"sort"
	"strconv"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type InterfaceListCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
	Sortby         string
}

func NewInterfaceListCommand(sl *metadata.SoftlayerCommand) (cmd *InterfaceListCommand) {
	thisCmd := &InterfaceListCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "interface-list " + T("SECURITYGROUP_ID"),
		Short: T("List interfaces associated with security group"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVar(&thisCmd.Sortby, "sortby", "", T("Column to sort by. Options are: id,virtualServerId,hostname"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *InterfaceListCommand) Run(args []string) error {
	groupID, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Security group ID")
	}

	outputFormat := cmd.GetOutputFlag()

	sortColumns := []string{"id", "virtualServerId", "hostname"}
	sortby := cmd.Sortby
	if sortby != "" && utils.StringInSlice(sortby, sortColumns) == -1 {
		return errors.NewInvalidUsageError(T("--sortby {{.Column}} is not supported.", map[string]interface{}{"Column": sortby}))
	}
	securityGroup, err := cmd.NetworkManager.GetSecurityGroup(groupID, "")
	if err != nil {
		return errors.NewAPIError(T("Failed to get security group {{.GroupID}}.\n", map[string]interface{}{"GroupID": groupID}), err.Error(), 2)
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
