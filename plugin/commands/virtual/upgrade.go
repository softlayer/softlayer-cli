package virtual

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type UpgradeCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
	Cpu                  int
	Private              bool
	Memory               int
	Network              int
	Flavor               string
	Force                bool
	AddDisk              int
	ResizeDisk           string
}

func NewUpgradeCommand(sl *metadata.SoftlayerCommand) (cmd *UpgradeCommand) {
	thisCmd := &UpgradeCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "upgrade " + T("IDENTIFIER"),
		Short: T("Upgrade a virtual server instance"),
		Long: T(`${COMMAND_NAME} sl vs upgrade IDENTIFIER [OPTIONS]
	Note: Classic infrastructure service automatically reboots the instance once upgrade request is
  	placed. The instance is halted until the upgrade transaction is completed.
  	However for Network, no reboot is required.

EXAMPLE:
   ${COMMAND_NAME} sl vs upgrade 12345678 -c 8 -m 8192 --network 1000
   This commands upgrades virtual server instance with ID 12345678 and set number of CPU cores to 8, memory to 8192M, network port speed to 1000 Mbps.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	cobraCmd.Flags().IntVarP(&thisCmd.Cpu, "cpu", "c", 0, T("Number of CPU cores"))
	cobraCmd.Flags().BoolVar(&thisCmd.Private, "private", false, T("CPU core will be on a dedicated host server"))
	cobraCmd.Flags().IntVarP(&thisCmd.Memory, "memory", "m", 0, T("Memory in megabytes"))
	// -1 as default since 0 is a valid value here
	cobraCmd.Flags().IntVar(&thisCmd.Network, "network", -1, T("Network port speed in Mbps"))
	cobraCmd.Flags().StringVar(&thisCmd.Flavor, "flavor", "", T("Flavor key name"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))
	cobraCmd.Flags().IntVar(&thisCmd.AddDisk, "add-disk", -1, T("Add Hard disk in GB"))
	cobraCmd.Flags().StringVar(&thisCmd.ResizeDisk, "resize-disk", "", T("Update disk number to size in GB [capacity,diskNumber]. --resize-disk 250,2"))
	return thisCmd
}

func (cmd *UpgradeCommand) Run(args []string) error {

	vsID, err := utils.ResolveVirtualGuestId(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	if cmd.Private && cmd.Cpu == 0 {
		return slErrors.NewInvalidUsageError(T("Must specify [--cpu] when using [--private]."))
	}

	if cmd.Cpu == 0 && cmd.Memory == 0 && cmd.Network == -1 && cmd.AddDisk == -1 && cmd.ResizeDisk == "" && cmd.Flavor == "" {
		return slErrors.NewInvalidUsageError(T("Must provide [--cpu], [--memory], [--network], [--add-disk], [--resize-disk] or [--flavor] to upgrade."))
	}

	if cmd.Flavor != "" && (cmd.Cpu != 0 || cmd.Memory != 0 || cmd.Private) {
		return slErrors.NewInvalidUsageError(T("Option [--flavor] is exclusive with [--cpu], [--memory] and [--private]."))
	}

	outputFormat := cmd.GetOutputFlag()

	if !cmd.Force && outputFormat != "JSON" {
		confirm, err := cmd.UI.Confirm(T("This action will incur charges on your account. Continue?"))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	subs := map[string]interface{}{"VsID": vsID, "VsId": vsID, "OrderId": 0}
	resizeDiskValues := []int{}
	if cmd.ResizeDisk != "" {
		resizeDiskStringValues := strings.Split(cmd.ResizeDisk, ",")
		if len(resizeDiskStringValues) != 2 {
			return slErrors.NewInvalidUsageError(T("--resize-disk requires capacity and disk number values separated by one comma."))
		}
		capacity, err := strconv.Atoi(resizeDiskStringValues[0])
		if err != nil {
			return slErrors.NewInvalidSoftlayerIdInputError("--resize-disk capacity")
		}
		diskNumber, err := strconv.Atoi(resizeDiskStringValues[1])
		if err != nil {
			return slErrors.NewInvalidSoftlayerIdInputError("--resize-disk disk number")
		}
		resizeDiskValues = []int{capacity, diskNumber}
	}
	orderReceipt, err := cmd.VirtualServerManager.UpgradeInstance(vsID, cmd.Cpu, cmd.Memory, cmd.Network, cmd.AddDisk, resizeDiskValues, cmd.Private, cmd.Flavor)
	if err != nil {
		return slErrors.NewAPIError(T("Failed to upgrade virtual server instance: {{.VsID}}.\n", subs), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, orderReceipt)
	}

	cmd.UI.Ok()
	subs["OrderId"] = *orderReceipt.OrderId
	cmd.UI.Print(T("Order {{.OrderId}} to upgrade virtual server instance: {{.VsId}} was placed.", subs))

	return nil
}
