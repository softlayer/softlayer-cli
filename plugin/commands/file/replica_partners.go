package file

import (
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

type ReplicaPartnersCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewReplicaPartnersCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *ReplicaPartnersCommand) {
	return &ReplicaPartnersCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func FileReplicaPartnersMetaData() cli.Command {
	return cli.Command{
		Category:    "file",
		Name:        "replica-partners",
		Description: T("List existing replicant volumes for a file volume"),
		Usage: T(`${COMMAND_NAME} sl file replica-partners VOLUME_ID [OPTIONS]
		
EXAMPLE:
   ${COMMAND_NAME} sl file replica-partners 12345678
   This command lists existing replicant volumes for file volume with ID 12345678.`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}


func (cmd *ReplicaPartnersCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	volumeID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Volume ID")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	partners, err := cmd.StorageManager.GetReplicationPartners(volumeID)
	if err != nil {
		return cli.NewExitError(T("Failed to get replication partners for volume {{.VolumeID}}.\n", map[string]interface{}{"VolumeID": volumeID})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, partners)
	}

	if len(partners) == 0 {
		cmd.UI.Print(T("There are no replication partners for volume {{.VolumeID}}.\n", map[string]interface{}{"VolumeID": volumeID}))
	} else {
		table := cmd.UI.Table([]string{T("ID"), T("User name"), T("Account ID"), T("Capacity (GB)"), T("Hardware ID"), T("Guest ID"), T("Host ID")})
		for _, p := range partners {
			table.Add(
				utils.FormatIntPointer(p.Id),
				utils.FormatStringPointer(p.Username),
				utils.FormatIntPointer(p.AccountId),
				utils.FormatIntPointer(p.CapacityGb),
				utils.FormatIntPointer(p.HardwareId),
				utils.FormatIntPointer(p.GuestId),
				utils.FormatIntPointer(p.HostId),
			)
		}
		table.Print()
	}
	return nil
}
