package block

import (
	"bytes"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

var (
	volumeSizes = []string{"20", "40", "80", "100", "250", "500", "1000", "2000-3000", "4000-7000", "8000-9000", "10000-12000"}
)

type VolumeOptionsCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewVolumeOptionsCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *VolumeOptionsCommand) {
	return &VolumeOptionsCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func BlockVolumeOptionsMetaData() cli.Command {
	return cli.Command{
		Category:    "block",
		Name:        "volume-options",
		Description: T("List all options for ordering a block storage"),
		Usage: T(`${COMMAND_NAME} sl block volume-options
	
EXAMPLE:
   ${COMMAND_NAME} sl block volume-options
   This command lists all options for creating a block storage volume, including storage type, volume size, OS type, IOPS, tier level, datacenter, and snapshot size.`),
	}
}

//refer to here about volume size and iops ranges: http://knowledgelayer.softlayer.com/learning/block-storage
func (cmd *VolumeOptionsCommand) Run(c *cli.Context) error {
	table := cmd.UI.Table([]string{"name", "value"})
	locations, err := cmd.StorageManager.GetAllDatacenters()
	if err != nil {
		return cli.NewExitError(T("Failed to get all datacenters.\n")+err.Error(), 2)
	}
	table.Add(T("Storage Type"), "performance,endurance")
	table.Add(T("Size (GB)"), utils.StringSliceToString(volumeSizes))
	table.Add(T("OS Type"), "HYPER_V,LINUX,VMWARE,WINDOWS_2008,WINDOWS_GPT,WINDOWS,XEN")
	buf := new(bytes.Buffer)
	iopsTable := terminal.NewTable(buf, append([]string{T("Size (GB)")}, volumeSizes...))
	iopsTable.Add(T("Min IOPS"), "100", "100", "100", "100", "100", "100", "100", "200", "300", "500", "1000")
	iopsTable.Add(T("Max IOPS"), "1000", "2000", "4000", "6000", "6000", "6000 or 10000", "6000 or 20000", "6000 or 40000", "6000 or 48000", "6000 or 48000", "6000 or 48000")
	iopsTable.Print()
	table.Add(T("IOPS"), buf.String())
	table.Add(T("Tier"), "0.25,2,4,10")
	table.Add(T("Location"), utils.StringSliceToString(locations))
	buf = new(bytes.Buffer)
	snapshotTable := terminal.NewTable(buf, []string{T("Storage Size (GB)"), T("Available Snapshot Size (GB)")})
	snapshotTable.Add(volumeSizes[0], "0,5,10,20")
	snapshotTable.Add(volumeSizes[1], "0,5,10,20,40")
	snapshotTable.Add(volumeSizes[2], "0,5,10,20,40,60,80")
	snapshotTable.Add(volumeSizes[3], "0,5,10,20,40,60,80,100")
	snapshotTable.Add(volumeSizes[4], "0,5,10,20,40,60,80,100,150,200,250")
	snapshotTable.Add(volumeSizes[5], "0,5,10,20,40,60,80,100,150,200,250,300,350,400,450,500")
	snapshotTable.Add(volumeSizes[6], "0,5,10,20,40,60,80,100,150,200,250,300,350,400,450,500,600,700,1000")
	snapshotTable.Add(volumeSizes[7], "0,5,10,20,40,60,80,100,150,200,250,300,350,400,450,500,600,700,1000,2000")
	snapshotTable.Add(volumeSizes[8], "0,5,10,20,40,60,80,100,150,200,250,300,350,400,450,500,600,700,1000,2000,4000")
	snapshotTable.Add(volumeSizes[9], "0,5,10,20,40,60,80,100,150,200,250,300,350,400,450,500,600,700,1000,2000,4000")
	snapshotTable.Add(volumeSizes[10], "0,5,10,20,40,60,80,100,150,200,250,300,350,400,450,500,600,700,1000,2000,4000")
	snapshotTable.Print()
	table.Add(T("Snapshot Size (GB)"), buf.String())
	table.Add(T("Note:"), T("IOPs limit above 6000 available in select data centers, refer to:http://knowledgelayer.softlayer.com/articles/new-ibm-block-and-file-storage-location-and-features"))
	table.Print()
	return nil
}
