package vlan

import (
	"bytes"

	"sync"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type OptionsCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
}

func NewOptionsCommand(ui terminal.UI, networkManager managers.NetworkManager) (cmd *OptionsCommand) {
	return &OptionsCommand{
		UI:             ui,
		NetworkManager: networkManager,
	}
}

func (cmd *OptionsCommand) Run(c *cli.Context) error {
	table := cmd.UI.Table([]string{T("name"), T("value")})
	datacenters, err := cmd.NetworkManager.ListDatacenters()
	if err != nil {
		return cli.NewExitError(T("Failed to list datacenters.\n")+err.Error(), 2)
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

func VlanOptionsMetaData() cli.Command {
	return cli.Command{
		Category:    "vlan",
		Name:        "options",
		Description: T("List all the options for creating VLAN"),
		Usage: T(`${COMMAND_NAME} sl vlan options
	
EXAMPLE:
   ${COMMAND_NAME} sl vlan options
   This command lists all options for creating a vlan, eg. vlan type, datacenters, subnet size, routers, etc.`),
	}
}
