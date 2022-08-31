package vlan

import (
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CreateCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
	VlanType       string
	Router         string
	Datacenter     string
	Name           string
	Force          bool
}

func NewCreateCommand(sl *metadata.SoftlayerCommand) *CreateCommand {
	thisCmd := &CreateCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "create",
		Short: T("Create a new VLAN."),
		Long: T(`${COMMAND_NAME} sl vlan create [OPTIONS]
	
EXAMPLE:
	${COMMAND_NAME} sl vlan create -t public -d dal09 -n myvlan
	This command creates a public vlan located in datacenter dal09 named "myvlan".
	${COMMAND_NAME} sl vlan create -r bcr01a.dal09 -n myvlan
	This command creates a vlan on router bcr01a.dal09 named "myvlan".`),
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVarP(&thisCmd.VlanType, "vlan-type", "t", "", T("The type of the VLAN, either public or private"))
	cobraCmd.Flags().StringVarP(&thisCmd.Router, "router", "r", "", T("The hostname of the router"))
	cobraCmd.Flags().StringVarP(&thisCmd.Datacenter, "datacenter", "d", "", T("The short name of the datacenter"))
	cobraCmd.Flags().StringVarP(&thisCmd.Name, "name", "n", "", T("The name of the VLAN"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CreateCommand) Run(args []string) error {
	if cmd.Router != "" {
		//set routers, then no need to set vlan-type or datacenter
		if cmd.Datacenter != "" || cmd.VlanType != "" {
			return errors.NewInvalidUsageError(T("[-r|--router] is not allowed with [-d|--datacenter] or [-t|--vlan-type].\nRun '{{.CommandName}} sl vlan options' to check available options."))
		}
	} else {
		//not set router, then need to set both vlan-type and datacenter
		if cmd.Datacenter == "" || cmd.VlanType == "" {
			return errors.NewInvalidUsageError(T("[-d|--datacenter] and [-t|--vlan-type] are required.\nRun '{{.CommandName}} sl vlan options' to check available options."))
		}
		vlanType := cmd.VlanType
		if vlanType != "public" && vlanType != "private" {
			return errors.NewInvalidUsageError(T("[-t|--vlan-type] is required, must be either public or private.\nRun '{{.CommandName}} sl vlan options' to check available options."))
		}

	}

	outputFormat := cmd.GetOutputFlag()

	if !cmd.Force {
		confirm, err := cmd.UI.Confirm(T("This action will incur charges on your account. Continue?"))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	orderReceipt, err := cmd.NetworkManager.AddVlan(cmd.VlanType, cmd.Datacenter, cmd.Router, cmd.Name)
	if err != nil {
		return errors.NewAPIError(T("Failed to add VLAN.\n"), err.Error(), 2)
	}
	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, orderReceipt)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("The order {{.OrderID}} was placed.", map[string]interface{}{"OrderID": *orderReceipt.OrderId}))
	return nil
}
