package virtual

import (
	"fmt"
	"sort"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
)

type CreateOptionsCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
}

func NewCreateOptionsCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager) (cmd *CreateOptionsCommand) {
	return &CreateOptionsCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
	}
}

func (cmd *CreateOptionsCommand) Run(c *cli.Context) error {
	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	createOptions, err := cmd.VirtualServerManager.GetCreateOptions()
	if err != nil {
		return cli.NewExitError(T("Failed to get virtual server creation options.\n")+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, createOptions)
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	//datacenter
	var datacenters []string
	for _, datacenter := range createOptions.Datacenters {
		if datacenter.Template != nil && datacenter.Template.Datacenter != nil {
			datacenters = append(datacenters, utils.FormatStringPointer(datacenter.Template.Datacenter.Name))
		}
	}
	table.Add(T("datacenter"), utils.StringSliceToString(datacenters))

	//flavor
	for _, flavorKey := range managers.FlavorKeys {
		var flavors []string
		for _, f := range createOptions.Flavors {
			if f.Flavor != nil && f.Flavor.KeyName != nil {
				keyName := *f.Flavor.KeyName
				if strings.Index(keyName, flavorKey) == -1 {
					continue
				}
				flavors = append(flavors, keyName)
			}
		}
		table.Add(T("flavors ({{.Flavor}})", map[string]interface{}{"Flavor": managers.FlavorLabels[flavorKey]}), utils.StringSliceToString(flavors))
	}

	//cpus
	var standardCpus, dedicatedCpus, dedicatedHostCpus []int
	for _, processor := range createOptions.Processors {
		cpu := processor.Template.StartCpus
		if processor.Template.DedicatedAccountHostOnlyFlag == nil && processor.Template.DedicatedHost == nil {
			standardCpus = append(standardCpus, *cpu)
		} else if processor.Template != nil && processor.Template.DedicatedAccountHostOnlyFlag != nil && *processor.Template.DedicatedAccountHostOnlyFlag == true {
			dedicatedCpus = append(dedicatedCpus, *cpu)
		} else if processor.Template.DedicatedHost != nil {
			dedicatedHostCpus = append(dedicatedHostCpus, *cpu)
		}
	}
	sort.Ints(standardCpus)
	sort.Ints(dedicatedCpus)
	sort.Ints(dedicatedHostCpus)
	table.Add(T("cpu (standard)"), utils.IntSliceToString(standardCpus))
	table.Add(T("cpu (dedicated)"), utils.IntSliceToString(dedicatedCpus))
	table.Add(T("cpu (dedicated host)"), utils.IntSliceToString(dedicatedHostCpus))

	//memory
	var mems, dedicatedHostMems []int
	for _, mem := range createOptions.Memory {
		if mem.ItemPrice != nil && mem.ItemPrice.DedicatedHostInstanceFlag != nil && *mem.ItemPrice.DedicatedHostInstanceFlag == false {
			dedicatedHostMems = append(dedicatedHostMems, *mem.Template.MaxMemory)
		} else if mem.Template != nil && mem.Template.MaxMemory != nil {
			mems = append(mems, *mem.Template.MaxMemory)
		}
	}
	sort.Ints(mems)
	sort.Ints(dedicatedHostMems)
	table.Add(T("memory"), utils.IntSliceToString(mems))
	table.Add(T("memory (dedicated host)"), utils.IntSliceToString(dedicatedHostMems))

	//operating system
	osCodes := make(map[string][]string)
	for _, ops := range createOptions.OperatingSystems {
		if ops.Template != nil && ops.Template.OperatingSystemReferenceCode != nil {
			osCode := *ops.Template.OperatingSystemReferenceCode
			category := strings.Split(osCode, "_")[0]
			osCodes[category] = append(osCodes[category], osCode)
		}
	}
	var osKeys []string
	for key := range osCodes {
		osKeys = append(osKeys, key)
	}
	sort.Strings(osKeys)
	for _, key := range osKeys {
		var oss []string
		for _, value := range osCodes[key] {
			oss = append(oss, value)
		}
		sort.Strings(oss)
		table.Add(T("os ({{.OS}})", map[string]interface{}{"OS": key}), utils.StringSliceToString(oss))
	}

	//local disk/sandisk
	localDisks := make(map[string][]int)
	dedicatedDisks := make(map[string][]int)
	sanDisks := make(map[string][]int)
	for _, disk := range createOptions.BlockDevices {
		deviceNumber := disk.Template.BlockDevices[0].Device
		capacity := disk.Template.BlockDevices[0].DiskImage.Capacity
		if disk.Template != nil && disk.Template.LocalDiskFlag != nil && disk.ItemPrice != nil && disk.ItemPrice.DedicatedHostInstanceFlag != nil {
			if *disk.Template.LocalDiskFlag == false {
				key := fmt.Sprintf("san disk (%s)", *deviceNumber)
				sanDisks[key] = append(sanDisks[key], *capacity)
			} else {
				if *disk.ItemPrice.DedicatedHostInstanceFlag == true {
					key := fmt.Sprintf("local (dedicated host) disk (%s)", *deviceNumber)
					dedicatedDisks[key] = append(dedicatedDisks[key], *capacity)
				} else {
					key := fmt.Sprintf("local disk (%s)", *deviceNumber)
					localDisks[key] = append(localDisks[key], *capacity)
				}
			}
		}
	}
	keys := make([]string, 0, len(sanDisks))
	for key := range sanDisks {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		table.Add(key, utils.IntSliceToString(sanDisks[key]))
	}
	keys = make([]string, 0, len(localDisks))
	for key := range localDisks {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		table.Add(key, utils.IntSliceToString(localDisks[key]))
	}
	keys = make([]string, 0, len(dedicatedDisks))
	for key := range dedicatedDisks {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		table.Add(key, utils.IntSliceToString(dedicatedDisks[key]))
	}

	//nic
	var nics, dedicatedNics []int
	for _, nic := range createOptions.NetworkComponents {
		if nic.Template != nil && nic.ItemPrice != nil && nic.Template.NetworkComponents != nil && len(nic.Template.NetworkComponents) > 0 {
			if nic.Template.NetworkComponents[0].MaxSpeed != nil {
				maxSpeed := *nic.Template.NetworkComponents[0].MaxSpeed
				if nic.ItemPrice.DedicatedHostInstanceFlag != nil && *nic.ItemPrice.DedicatedHostInstanceFlag {
					if utils.IntInSlice(maxSpeed, dedicatedNics) == -1 {
						dedicatedNics = append(dedicatedNics, maxSpeed)
					}
				} else {
					if utils.IntInSlice(maxSpeed, nics) == -1 {
						nics = append(nics, maxSpeed)
					}
				}
			}
		}
	}
	sort.Ints(nics)
	sort.Ints(dedicatedNics)
	table.Add(T("nic"), utils.IntSliceToString(nics))
	table.Add(T("nic (dedicated host)"), utils.IntSliceToString(dedicatedNics))
	table.Print()
	return nil
}
