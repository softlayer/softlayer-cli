package search

import (
	"fmt"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/spf13/cobra"
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type SearchCommand struct {
	*metadata.SoftlayerCommand
	SearchManager managers.SearchManager
	Command       *cobra.Command
	Query         string
}

func SetupCobraCommands(sl *metadata.SoftlayerCommand) *cobra.Command {
	thisCmd := &SearchCommand{
		SoftlayerCommand: sl,
		SearchManager:    managers.NewSearchManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "search",
		Short: T("Perform a query against the SoftLayer search database."),
		Long: T(`Read More: https://sldn.softlayer.com/reference/services/SoftLayer_Search/search/
Examples::

    sl search --query 'test.com'
    sl search --query '_objectType:SoftLayer_Virtual_Guest test.com'
`),
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().StringVarP(&thisCmd.Query, "query", "q", "", T("The search query you want to use."))
	cobraCmd.AddCommand(NewSearchTypesCommand(sl).Command)
	return cobraCmd
}

func SearchNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "search",
		Description: T("Perform a query against the SoftLayer search database."),
	}
}

func (cmd *SearchCommand) Run(args []string) error {

	results, err := cmd.SearchManager.AdvancedSearch("", cmd.Query)
	if err != nil {
		return err
	}
	if cmd.GetOutputFlag() == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, results)
	}
	table := cmd.UI.Table([]string{"Type", "Matched Terms", "Resource"})
	for _, result := range results {
		resource_string := ""
		switch *result.ResourceType {
		case "SoftLayer_Virtual_Guest":
			resource_object := result.Resource.(*datatypes.Virtual_Guest)
			resource_string = parseVirtual_Guest(*resource_object)
		case "SoftLayer_Event_Log":
			resource_object := result.Resource.(*datatypes.Event_Log)
			resource_string = parseEvent_Log(*resource_object)
		case "SoftLayer_Virtual_DedicatedHost":
			resource_object := result.Resource.(*datatypes.Virtual_DedicatedHost)
			resource_string = parseVirtual_DedicatedHost(*resource_object)
		case "SoftLayer_Hardware":
			resource_object := result.Resource.(*datatypes.Hardware)
			resource_string = parseHardware(*resource_object)
		case "SoftLayer_Network_Application_Delivery_Controller":
			resource_object := result.Resource.(*datatypes.Network_Application_Delivery_Controller)
			resource_string = parseNetwork_Application_Delivery_Controller(*resource_object)
		case "SoftLayer_Network_Subnet_IpAddress":
			resource_object := result.Resource.(*datatypes.Network_Subnet_IpAddress)
			resource_string = parseNetwork_Subnet_IpAddress(*resource_object)
		case "SoftLayer_Network_Vlan":
			resource_object := result.Resource.(*datatypes.Network_Vlan)
			resource_string = parseNetwork_Vlan(*resource_object)
		case "SoftLayer_Network_Vlan_Firewall":
			resource_object := result.Resource.(*datatypes.Network_Vlan_Firewall)
			resource_string = parseNetwork_Vlan_Firewall(*resource_object)
		case "SoftLayer_Ticket":
			resource_object := result.Resource.(*datatypes.Ticket)
			resource_string = parseTicket(*resource_object)
		}
		table.Add(*result.ResourceType, parseMatchedTerms(result), resource_string)
	}
	table.Print()
	return nil
}

func parseMatchedTerms(searchResult datatypes.Container_Search_Result) string {
	return strings.Join(searchResult.MatchedTerms, "\n")
}
func parseVirtual_Guest(resource datatypes.Virtual_Guest) string {
	return fmt.Sprintf(T("ID")+": %d\n"+T("FQDN")+": %s\n", *resource.Id, *resource.FullyQualifiedDomainName)
}
func parseEvent_Log(resource datatypes.Event_Log) string {
	return fmt.Sprintf(T("ID")+": %s\n"+T("Event")+": %s\n", *resource.TraceId, *resource.EventName)
}
func parseVirtual_DedicatedHost(resource datatypes.Virtual_DedicatedHost) string {
	return fmt.Sprintf(T("ID")+": %d\n"+T("Name")+": %s\n", *resource.Id, *resource.Name)
}
func parseHardware(resource datatypes.Hardware) string {
	return fmt.Sprintf(T("ID")+": %d\n"+T("FQDN")+": %s\n", *resource.Id, *resource.FullyQualifiedDomainName)
}
func parseNetwork_Application_Delivery_Controller(resource datatypes.Network_Application_Delivery_Controller) string {
	return fmt.Sprintf(T("ID")+": %d\n"+T("Name")+": %s\n", *resource.Id, *resource.Name)
}
func parseNetwork_Subnet_IpAddress(resource datatypes.Network_Subnet_IpAddress) string {
	return fmt.Sprintf(T("ID")+": %d\n"+T("Ip Address")+": %s\n", *resource.Id, *resource.IpAddress)
}
func parseNetwork_Vlan(resource datatypes.Network_Vlan) string {
	return fmt.Sprintf(T("ID")+": %d\n"+T("VLAN")+": %d\n", *resource.Id, *resource.VlanNumber)
}
func parseNetwork_Vlan_Firewall(resource datatypes.Network_Vlan_Firewall) string {
	return fmt.Sprintf(T("ID")+": %d\n"+T("Ip Address")+": %s\n", *resource.Id, *resource.PrimaryIpAddress)
}
func parseTicket(resource datatypes.Ticket) string {
	return fmt.Sprintf(T("ID")+": %d\n"+T("Subject")+": %s\n", *resource.Id, *resource.Title)
}
