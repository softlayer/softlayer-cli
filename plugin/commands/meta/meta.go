package meta

import (
	"fmt"
	"strings"

	"github.com/softlayer/softlayer-go/sl"
	"github.com/spf13/cobra"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)



type MetaCommand struct {
	*metadata.SoftlayerCommand
	Manager managers.MetadataManager
	Command              *cobra.Command
}

func NewMetaCommand(sl *metadata.SoftlayerCommand) *MetaCommand {
	thisCmd := &MetaCommand{
		SoftlayerCommand:     sl,
		Manager: managers.NewMetadataManager(sl.Session),
	}
	validOptions := availableMetadataOptions()
	cobraCmd := &cobra.Command{
		Use:   "metadata " + T("ARGUMENT"),
		Short: T("Find details about the machine making these API calls."),
		Long: T("ARGUMENT Choices: " + strings.Join(validOptions, ", ")),
		Args:  cobra.MatchAll(metadata.OneArgs, cobra.OnlyValidArgs),
		ValidArgs: validOptions,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}



func (cmd *MetaCommand) Run(args []string) error {

	option := args[0]
	if utils.StringInSlice(option, availableMetadataOptions()) == -1 {
		return slErrors.NewInvalidUsageError(T("This option is not available."))
	}

	var options sl.Options
	if option == "network" {
		arrayNetwork := []string{}
		var macAddress string
		var parameter string

		for _, network := range NetworkRequest() {
			requestMac := strings.Contains(network, "MacAddresses")
			if !requestMac {
				parameter = fmt.Sprintf(`["%s"]`, macAddress)
			} else {
				parameter = ""
			}

			response, err := cmd.Manager.CallAPIService("SoftLayer_Resource_Metadata", network, options, parameter)
			if err != nil {
				return err
			}

			response = cleanResponse(response)
			if requestMac {
				macAddress = obtainMac(response)
			}

			arrayNetwork = append(arrayNetwork, response)
		}

		printNetwork(cmd.UI, arrayNetwork)
		return nil
	}

	request := availableMetadataRequests()[option]
	response, err := cmd.Manager.CallAPIService("SoftLayer_Resource_Metadata", request, options, "")
	if err != nil {
		return err
	}
	cmd.UI.Print(cleanResponse(response))
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

func obtainMac(arrayMacs string) string {
	firstMac := strings.Split(arrayMacs, ",")
	return firstMac[0][1 : len(firstMac[0])-1]
}

func cleanResponse(response string) string {
	if strings.Contains(response, "[") {
		return response[1 : len(response)-1]
	}
	return response
}
