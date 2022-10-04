package dedicatedhost

import (
	"sort"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CreateOptionsCommand struct {
	*metadata.SoftlayerCommand
	DedicatedHostManager managers.DedicatedHostManager
	Command              *cobra.Command
	Datacenter           string
	Flavor               string
}

func NewCreateOptionsCommand(sl *metadata.SoftlayerCommand) *CreateOptionsCommand {
	thisCmd := &CreateOptionsCommand{
		SoftlayerCommand:     sl,
		DedicatedHostManager: managers.NewDedicatedhostManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "create-options",
		Short: T("Host order options for a given dedicated host."),
		Long: T(`${COMMAND_NAME} sl dedicatedhost create-options [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl dedicatedhost create-options

   To get the list of available private vlans use this command: ${COMMAND_NAME} sl dedicatedhost create-options --datacenter dal05 --flavor 56_CORES_X_242_RAM_X_1_4_TB"`),
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVarP(&thisCmd.Datacenter, "datacenter", "d", "", T("Filter private vlans by Datacenter shortname e.g. ams01, (requires --flavor)"))
	cobraCmd.Flags().StringVarP(&thisCmd.Flavor, "flavor", "f", "", T("Dedicated Virtual Host flavor (requires --datacenter) e.g. 56_CORES_X_242_RAM_X_1_4_TB"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CreateOptionsCommand) Run(args []string) error {
	var productPackage datatypes.Product_Package
	productPackage, err := cmd.DedicatedHostManager.GetPackage()
	if err != nil {
		return errors.NewAPIError(T("Failed to get product package for hardware server."), err.Error(), 2)
	}
	if cmd.Datacenter == "" && cmd.Flavor == "" {
		options := cmd.DedicatedHostManager.GetCreateOptions(productPackage)

		//datacenters
		dcTable := cmd.UI.Table([]string{T("Data center"), T("Value")})
		locations := options[managers.LOCATIONS]
		var sortedLocations []string
		for key, _ := range locations {
			sortedLocations = append(sortedLocations, key)
		}
		sort.Strings(sortedLocations)
		for _, key := range sortedLocations {
			dcTable.Add(locations[key], key)
		}
		dcTable.Print()
		cmd.UI.Print("")

		//Dedicated Virtual Host Flavor(s)
		flavorTable := cmd.UI.Table([]string{T("Dedicated Virtual Host Flavor(s)"), T("Value")})
		flavors := options[managers.DEDICATED_HOST]
		var sortedFlavors []string
		for key, _ := range flavors {
			sortedFlavors = append(sortedFlavors, key)
		}
		sort.Strings(sortedFlavors)
		for _, key := range sortedFlavors {
			flavorTable.Add(flavors[key], key)
		}
		flavorTable.Print()
		cmd.UI.Print("")
	} else {
		if (cmd.Datacenter != "" && cmd.Flavor == "") || (cmd.Datacenter == "" && cmd.Flavor != "") {
			return errors.NewMissingInputError("Both -d|--datacenter and -f|--flavor need to be passed as arguments e.g. ibmcloud sl dedicatedhost create-options -d ams01 -f 56_CORES_X_242_RAM_X_1_4_TB")
		}
		privateVlans, err := cmd.DedicatedHostManager.GetVlansOptions(cmd.Datacenter, cmd.Flavor, productPackage)
		if err != nil {
			return errors.NewAPIError(T("Failed to get the vlans available for datacener: {{.DATACENTER}} and flavor: {{.FLAVOR}}.", map[string]interface{}{"DATACENTER": cmd.Datacenter, "FLAVOR": cmd.Flavor}), err.Error(), 2)
		}
		table := cmd.UI.Table([]string{T("Id"), T("Name"), T("PrimaryRouter Hostname")})
		for _, vlans := range privateVlans {
			table.Add(utils.FormatIntPointer(vlans.Id), utils.FormatStringPointer(vlans.Name), utils.FormatStringPointer(vlans.PrimaryRouter.Hostname))
		}
		table.Print()
	}

	return nil
}
