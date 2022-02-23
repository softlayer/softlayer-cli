package block

import (
	"fmt"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ReplicaOrderCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
	Context        plugin.PluginContext
}

func NewReplicaOrderCommand(ui terminal.UI, storageManager managers.StorageManager, context plugin.PluginContext) (cmd *ReplicaOrderCommand) {
	return &ReplicaOrderCommand{
		UI:             ui,
		StorageManager: storageManager,
		Context:        context,
	}
}

func BlockReplicaOrderMetaData() cli.Command {
	return cli.Command{
		Category:    "block",
		Name:        "replica-order",
		Description: T("Order a block storage replica volume"),
		Usage: T(`${COMMAND_NAME} sl block replica-order VOLUME_ID [OPTIONS]
		
EXAMPLE:
   ${COMMAND_NAME} sl block replica-order 12345678 -s DAILY -d dal09 --tier 4 --os-type LINUX
   This command orders a replica for volume with ID 12345678, which performs DAILY replication, is located at dal09, tier level is 4, OS type is Linux.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "s,snapshot-schedule",
				Usage: T("Snapshot schedule to use for replication. Options are: HOURLY,DAILY,WEEKLY [required]"),
			},
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Short name of the datacenter for the replica. For example, dal09 [required]"),
			},
			cli.Float64Flag{
				Name:  "t,tier",
				Usage: T("Endurance Storage Tier (IOPS per GB) of the primary volume for which a replica is ordered [optional], options are: 0.25,2,4,10,if no tier is specified, the tier of the original volume will be used"),
			},
			cli.IntFlag{
				Name:  "i,iops",
				Usage: T("Performance Storage IOPs, between 100 and 6000 in multiples of 100,if no IOPS value is specified, the IOPS value of the original volume will be used"),
			},
			cli.StringFlag{
				Name:  "o,os-type",
				Usage: T("Operating System Type (eg. LINUX) of the primary volume for which a replica is ordered [optional], options are: HYPER_V,LINUX,VMWARE,WINDOWS_2008,WINDOWS_GPT,WINDOWS,XEN"),
			},
			metadata.ForceFlag(),
			metadata.OutputFlag(),
		},
	}
}

func (cmd *ReplicaOrderCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	volumeID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}

	snapshotSchedule := c.String("s")
	if snapshotSchedule == "" || (snapshotSchedule != "HOURLY" && snapshotSchedule != "DAILY" && snapshotSchedule != "WEEKLY") {
		return errors.NewInvalidUsageError(T("[-s|--snapshot-schedule] is required, options are: HOURLY, DAILY, WEEKLY."))
	}

	datacenter := c.String("d")
	if datacenter == "" {
		return errors.NewInvalidUsageError(T("[-d|--datacenter] is required.\n Run '{{.CommandName}} sl block volume-options' to get available options.",
			map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
	}

	tier := c.Float64("t")
	if tier != 0 {
		if tier != 0.25 && tier != 2 && tier != 4 && tier != 10 {
			return errors.NewInvalidUsageError(T("[-t|--tier] is optional, options are: 0.25,2,4,10."))
		}
	}
	iops := c.Int("i")
	if iops != 0 {
		if iops < 100 || iops > 6000 {
			return errors.NewInvalidUsageError(T("-i|--iops must be between 100 and 6000, inclusive.\nRun '{{.CommandName}} sl block volume-options' to check available options.",
				map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
		}
		if iops%100 != 0 {
			return errors.NewInvalidUsageError(T("-i|--iops must be a multiple of 100.\nRun '{{.CommandName}} sl block volume-options' to check available options.",
				map[string]interface{}{"CommandName": cmd.Context.CLIName()}))
		}
	}

	osType := c.String("o")
	if osType != "" {
		if osType != "HYPER_V" && osType != "LINUX" && osType != "VMWARE" && osType != "WINDOWS_2008" && osType != "WINDOWS_GPT" && osType != "WINDOWS" && osType != "XEN" {
			return errors.NewInvalidUsageError(T("-o|--os-type is optional, options are: HYPER_V,LINUX,VMWARE,WINDOWS_2008,WINDOWS_GPT,WINDOWS,XEN."))
		}
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	if !c.IsSet("f") && outputFormat != "JSON" {
		confirm, err := cmd.UI.Confirm(T("This action will incur charges on your account. Continue?"))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}
	orderReceipt, err := cmd.StorageManager.OrderReplicantVolume("block", volumeID, snapshotSchedule, datacenter, tier, iops, osType)
	if err != nil {
		return cli.NewExitError(T("Failed to order replicant for volume {{.VolumeID}}.Please verify your options and try again.\n", map[string]interface{}{"VolumeID": volumeID})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, orderReceipt)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Order {{.OrderID}} was placed.", map[string]interface{}{"OrderID": *orderReceipt.OrderId}))
	for _, item := range orderReceipt.PlacedOrder.Items {
		cmd.UI.Print(fmt.Sprintf(" > %s", *item.Description))
		cmd.UI.Print("")
	}
	return nil
}
