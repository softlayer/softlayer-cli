package virtual

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CapacityCreateCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
	Name                 string
	BackendRouterId      int
	Instances            int
	Flavor               string
	Test                 bool
	Force                bool
}

func NewCapacityCreateCommand(sl *metadata.SoftlayerCommand) (cmd *CapacityCreateCommand) {
	thisCmd := &CapacityCreateCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "capacity-create",
		Short: T("Create a Reserved Capacity instance."),
		Long: T(`${COMMAND_NAME} sl vs capacity-create [OPTIONS]
EXAMPLE:
${COMMAND_NAME} sl vs capacity-create -n myvsi -b 1234567 -fl C1_2X2_1_YEAR_TERM -i 2
This command orders a Reserved Capacity instance with name is myvsi, backendRouterId 1234567, flavor C1_2X2_1_YEAR_TERM and 2 instances,
${COMMAND_NAME} sl vs capacity-create --name myvsi --backendRouterId 1234567 --flavor C1_2X2_1_YEAR_TERM --instances 2 --test
This command tests whether the order is valid with above options before the order is actually placed.

WARNING: Reserved Capacity is on a yearly contract and not cancelable until the contract is expired.`),
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	cobraCmd.Flags().StringVarP(&thisCmd.Name, "name", "n", "", T("Name for your new reserved capacity  [required]"))
	cobraCmd.Flags().IntVarP(&thisCmd.BackendRouterId, "backendRouterId", "b", 0, T("BackendRouterId, create-options has a list of valid ids to use. [required]"))
	cobraCmd.Flags().IntVarP(&thisCmd.Instances, "instances", "i", 0, T("Number of VSI instances this capacity reservation can support. [required]"))
	cobraCmd.Flags().StringVarP(&thisCmd.Flavor, "flavor", "l", "", T(" Capacity keyname (C1_2X2_1_YEAR_TERM for example). [required]"))
	cobraCmd.Flags().BoolVar(&thisCmd.Test, "test", false, T(" Do not actually create the reserved capacity"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	return thisCmd
}
func (cmd *CapacityCreateCommand) Run(args []string) error {
	var params map[string]interface{}
	var err error
	capacity_create := datatypes.Container_Product_Order_Virtual_ReservedCapacity{}
	if cmd.Name == "" || cmd.BackendRouterId == 0 || cmd.Instances == 0 || cmd.Flavor == "" {
		confirm, err := cmd.UI.Confirm(T("Please make sure you know all the creation options by running command: '{{.CommandName}} sl vs options'. Continue?",
			map[string]interface{}{"CommandName": "ibmcloud"}))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
		params = make(map[string]interface{})
		params["backendRouterId"], _ = cmd.UI.Ask(T("backendRouterId: "))
		params["flavor"], _ = cmd.UI.Ask(T("flavor: "))
		params["quantity"], _ = cmd.UI.Ask(T("instances: "))
		params["test"] = cmd.Test
		//params["complexType"] = T("SoftLayer_Container_Product_Order_Virtual_ReservedCapacity")
		params["hourly"] = true
	} else {
		//create virtual reservedCapacity server with customized parameters
		params, err = cmd.verifyCapacityParams()
		if err != nil {
			return err
		}
	}

	orderReceipt, err := cmd.VirtualServerManager.GenerateInstanceCapacityCreationTemplate(&capacity_create, params)
	if err != nil {
		return err
	}
	createTable(cmd, orderReceipt, cmd.Test)

	return nil
}

func createTable(cmd *CapacityCreateCommand, receipt interface{}, set bool) {
	if set {
		val := receipt.(datatypes.Container_Product_Order)
		table := cmd.UI.Table([]string{T("name"), T("value")})
		table.Add("name", utils.FormatStringPointer(val.QuoteName))
		table.Add("location", utils.FormatStringPointer(val.LocationObject.LongName))
		if val.Prices != nil {
			for _, price := range val.Prices {
				table.Add("Contract", utils.FormatStringPointer(price.Item.Description))
			}
		}
		table.Add("Hourly Total", utils.FormatSLFloatPointerToInt(val.PostTaxRecurringHourly))
		table.Print()
	} else {
		val := receipt.(datatypes.Container_Product_Order_Receipt)
		table := cmd.UI.Table([]string{T("name"), T("value")})
		table.Add("Order Date", utils.FormatSLTimePointer(val.OrderDate))
		table.Add("Order Id", utils.FormatIntPointer(val.OrderId))
		table.Add("Status", *val.PlacedOrder.Status)
		table.Add("Hourly Total", utils.FormatSLFloatPointerToInt(val.OrderDetails.PostTaxRecurringHourly))
		table.Print()
	}
}

func (cmd *CapacityCreateCommand) verifyCapacityParams() (map[string]interface{}, error) {
	params := make(map[string]interface{})
	if cmd.Flavor != "" {
		params["flavor"] = cmd.Flavor
	}
	if cmd.BackendRouterId != 0 {
		params["backendRouterId"] = cmd.BackendRouterId
	}
	if cmd.Instances != 0 {
		params["quantity"] = cmd.Instances
	}
	if cmd.Name != "" {
		params["name"] = cmd.Name
	}
	if cmd.Test {
		params["test"] = cmd.Test
	}

	return params, nil
}
