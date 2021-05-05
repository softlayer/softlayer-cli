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

type ReplicaLocationsCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewReplicaLocationsCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *ReplicaLocationsCommand) {
	return &ReplicaLocationsCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func (cmd *ReplicaLocationsCommand) Run(c *cli.Context) error {
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

	datacenters, err := cmd.StorageManager.GetReplicationLocations(volumeID)
	if err != nil {
		return cli.NewExitError(T("Failed to get datacenters for volume {{.VolumeID}}.\n", map[string]interface{}{"VolumeID": volumeID})+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, datacenters)
	}

	if len(datacenters) == 0 {
		cmd.UI.Print(T("No data centers compatible for replication."))
	} else {
		table := cmd.UI.Table([]string{T("ID"), T("Short Name"), T("Long Name")})
		for _, d := range datacenters {
			table.Add(utils.FormatIntPointer(d.Id), utils.FormatStringPointer(d.Name), utils.FormatStringPointer(d.LongName))
		}
		table.Print()
	}
	return nil
}
