package file

import (
	"sort"
	"strconv"

	"github.com/spf13/cobra"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type VolumeCountCommand struct {
	*metadata.SoftlayerStorageCommand
	Command        *cobra.Command
	StorageManager managers.StorageManager
	Datacenter     string
}

func NewVolumeCountCommand(sl *metadata.SoftlayerStorageCommand) (cmd *VolumeCountCommand) {
	thisCmd := &VolumeCountCommand{
		SoftlayerStorageCommand: sl,
		StorageManager:          managers.NewStorageManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "volume-count",
		Short: T("List number of file storage volumes per datacenter"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVarP(&thisCmd.Datacenter, "datacenter", "d", "", T("Filter by datacenter shortname"))
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *VolumeCountCommand) Run(args []string) error {
	mask := "mask[id,serviceResource.datacenter.name]"
	volumes, err := cmd.StorageManager.ListVolumes(managers.VOLUME_TYPE_FILE, cmd.Datacenter, "", "", "", 0, mask)
	if err != nil {
		return slErr.NewAPIError(T("Failed to list volumes on your account.\n"), err.Error(), 2)
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
