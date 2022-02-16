package block

import (
	"sort"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type VolumeCountCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewVolumeCountCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *VolumeCountCommand) {
	return &VolumeCountCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func BlockVolumeCountMetaData() cli.Command {
	return cli.Command{
		Category:    "block",
		Name:        "volume-count",
		Description: T("List number of block storage volumes per datacenter"),
		Usage:       "${COMMAND_NAME} sl block volume-count [OPTIONS]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Filter by datacenter shortname"),
			},
		},
	}
}

func (cmd *VolumeCountCommand) Run(c *cli.Context) error {
	mask := "mask[id,serviceResource.datacenter.name]"
	volumes, err := cmd.StorageManager.ListVolumes(managers.VOLUME_TYPE_BLOCK, c.String("d"), "", "", 0, mask)
	if err != nil {
		return cli.NewExitError(T("Failed to list volumes on your account.\n")+err.Error(), 2)
	}
	result := make(map[string]int)
	for _, v := range volumes {
		if v.ServiceResource != nil && v.ServiceResource.Datacenter != nil && v.ServiceResource.Datacenter.Name != nil {
			datacenterName := *v.ServiceResource.Datacenter.Name
			if count, ok := result[datacenterName]; ok {
				result[datacenterName] = count + 1
			} else {
				result[datacenterName] = 1
			}
		}
	}
	var datacenters []string
	for key, _ := range result {
		datacenters = append(datacenters, key)
	}
	sort.Strings(datacenters)
	table := cmd.UI.Table([]string{T("Data center"), T("Count")})
	for _, dc := range datacenters {
		table.Add(dc, strconv.Itoa(result[dc]))
	}
	table.Print()
	return nil
}
