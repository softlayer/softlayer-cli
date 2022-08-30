package vlan

import (
	"bytes"

	"sync"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type OptionsCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
}

func NewOptionsCommand(sl *metadata.SoftlayerCommand) *OptionsCommand {
	thisCmd := &OptionsCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "options",
		Short: T("List all the options for creating VLAN."),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *OptionsCommand) Run(args []string) error {
	table := cmd.UI.Table([]string{T("name"), T("value")})
	datacenters, err := cmd.NetworkManager.ListDatacenters()
	if err != nil {
		return errors.NewAPIError(T("Failed to list datacenters.\n"), err.Error(), 2)
	}
	var datacenternames []string
	buf := new(bytes.Buffer)
	dTable := terminal.NewTable(buf, []string{T("datacenter"), T("hostname")})

	var wg sync.WaitGroup
	var mutex sync.Mutex
	for id, name := range datacenters {
		wg.Add(1)
		go func(id int, name string) {
			names, err := cmd.NetworkManager.ListRouters(id, "mask[hostname]")
			mutex.Lock()
			if err != nil {
				dTable.Add(name, err.Error())
			} else {
				dTable.Add(name, utils.StringSliceToString(names))
			}
			mutex.Unlock()
			wg.Done()
		}(id, name)
		datacenternames = append(datacenternames, name)
	}
	wg.Wait()
	dTable.Print()
	table.Add(T("VLAN type"), "public,private")
	table.Add(T("datacenters"), utils.StringSliceToString(datacenternames))
	table.Add(T("routers"), buf.String())
	table.Print()
	return nil
}
