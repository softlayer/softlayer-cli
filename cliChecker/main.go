package main

import (
	"fmt"
	"sort"
	// "github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	sl_plugin "github.ibm.com/SoftLayer/softlayer-cli/plugin"
)


func main() {
	fmt.Printf("IBMCLOUD SL Command Directory\n")

	slPlugin := new(sl_plugin.SoftlayerPlugin)
	slMeta := slPlugin.GetMetadata()
	sort.Slice(slMeta.Commands, func(i, j int) bool {
		one := fmt.Sprintf("%s %s", slMeta.Commands[i].Namespace, slMeta.Commands[i].Name)
		two := fmt.Sprintf("%s %s", slMeta.Commands[j].Namespace, slMeta.Commands[j].Name)
		return one < two
	})
	fmt.Printf("==============================================================\n")
	for _, slCmd := range slMeta.Commands {
		fmt.Printf("%s %s\n", slCmd.Namespace, slCmd.Name)
		sort.Slice(slCmd.Flags, func(i, j int) bool {
			return slCmd.Flags[i].Name < slCmd.Flags[j].Name
		})
		for _, slCmdFlag := range slCmd.Flags {
			fmt.Printf("\tFlag: %s: %s\n", slCmdFlag.Name, slCmdFlag.Description)
		}
		fmt.Printf("\t--------------------------------\n")
		fmt.Printf("\tDescription: %s\n", slCmd.Description)
		fmt.Printf("\t--------------------------------\n")
		fmt.Printf("\tUsage: %s\n", slCmd.Usage)
		fmt.Printf("==============================================================\n")
	}
}