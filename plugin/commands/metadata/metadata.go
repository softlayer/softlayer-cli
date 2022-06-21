package metadata

import (
	"fmt"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

func GetCommandActionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	metadataManager := managers.NewMetadataManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		"sl-metadata": func(c *cli.Context) error {
			return NewMetadataCommand(ui, metadataManager).Run(c)
		},
	}
	return CommandActionBindings
}

type MetadataCommand struct {
	UI              terminal.UI
	MetadataManager managers.MetadataManager
}

func NewMetadataCommand(ui terminal.UI, metadataManager managers.MetadataManager) (cmd *MetadataCommand) {
	return &MetadataCommand{
		UI:              ui,
		MetadataManager: metadataManager,
	}
}

func MetadataMetadata() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "metadata",
		Description: T("Find details about the machine making these API calls."),
		Usage: T(`${COMMAND_NAME} sl metadata {backend_ip|backend_mac|datacenter|datacenter_id|fqdn|frontend_mac|id|ip|network|provision_state|tags|user_data} [OPTIONS]
		
		.. csv-table:: Choices: 
	backend_ip     backend_mac     datacenter     datacenter_id     fqdn     frontend_mac     id     ip     network     provision_state     tags     user_data
	
	These commands only work on devices on the backend SoftLayer network. This allows for self-discovery for newly provisioned resources.`),
	}
}

func (cmd *MetadataCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	option := c.Args()[0]
	if utils.StringInSlice(option, availableMetadataOptions()) == -1 {
		return errors.NewInvalidUsageError(T("This option is not available."))
	}

	var options sl.Options
	if option == "network" {
		arrayNetwork := []string{}
		var macAddress string
		var parameter string

		for _, network := range NetworkRequest() {
			requestMac := strings.Contains(network, "MacAddresses")
			if !requestMac {
				parameter = fmt.Sprintf(`[{"macAddress":"%s"}]`, macAddress)
			} else {
				parameter = ""
			}

			response, err := cmd.MetadataManager.CallAPIService("SoftLayer_Resource_Metadata", network, options, parameter)
			if err != nil {
				return err
			}

			if requestMac {
				macAddress = response
			}

			arrayNetwork = append(arrayNetwork, response)
		}

		printNetwork(cmd.UI, arrayNetwork)
		return nil
	}

	request := availableMetadataRequests()[option]
	response, err := cmd.MetadataManager.CallAPIService("SoftLayer_Resource_Metadata", request, options, "")
	if err != nil {
		return err
	}
	cmd.UI.Print(response)
	return nil
}

func printNetwork(ui terminal.UI, arrayNetwork []string) {
	tableFront := ui.Table([]string{
		T("Name"),
		T("Value"),
	})
	tableFront.Add(T("Mac addresses"), arrayNetwork[0])
	tableFront.Add(T("Router"), arrayNetwork[1])
	tableFront.Add(T("Vlans"), arrayNetwork[2])
	tableFront.Add(T("Vlan ids"), arrayNetwork[3])
	tableFront.Print()

	tableBack := ui.Table([]string{
		T("Name"),
		T("Value"),
	})
	tableBack.Add(T("Mac addresses"), arrayNetwork[4])
	tableBack.Add(T("Router"), arrayNetwork[5])
	tableBack.Add(T("Vlans"), arrayNetwork[6])
	tableBack.Add(T("Vlan ids"), arrayNetwork[7])
	tableBack.Print()
}

func NetworkRequest() []string {
	NetworkRequest := []string{
		"getFrontendMacAddresses",
		"getRouter",
		"getVlans",
		"getVlanIds",
		"getBackendMacAddresses",
		"getRouter",
		"getVlans",
		"getVlanIds",
	}
	return NetworkRequest
}

func availableMetadataOptions() []string {
	availableMetadataOptions := []string{
		"backend_ip",
		"backend_mac",
		"datacenter",
		"datacenter_id",
		"fqdn",
		"frontend_mac",
		"id",
		"ip",
		"network",
		"provision_state",
		"tags",
		"user_data",
	}
	return availableMetadataOptions
}

func availableMetadataRequests() map[string]string {
	availableMetadataRequests := map[string]string{
		"backend_ip":          "getPrimaryBackendIpAddress", //primary_backend_ip
		"backend_mac":         "getBackendMacAddresses",
		"datacenter":          "getDatacenter",
		"datacenter_id":       "getDatacenterId",
		"domain":              "getDomain",
		"frontend_mac":        "getFrontendMacAddresses",
		"fqdn":                "getFullyQualifiedDomainName",
		"hostname":            "getHostname",
		"id":                  "getId",
		"ip":                  "getPrimaryIpAddress", // primary_ip
		"primary_frontend_ip": "getPrimaryIpAddress",
		"provision_state":     "getProvisionState",
		"router":              "getRouter",
		"tags":                "getTags",
		"user_data":           "getUserMetadata",
		"user_metadata":       "getUserMetadata",
		"vlan_ids":            "getVlanIds",
		"vlans":               "getVlans",
	}
	return availableMetadataRequests
}
