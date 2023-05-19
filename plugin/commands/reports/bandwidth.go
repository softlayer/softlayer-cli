package reports

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"

	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type BandwidthCommand struct {
	*metadata.SoftlayerCommand
	SearchManager managers.SearchManager
	Command       *cobra.Command
}

func NewBandwidthCommand(sl *metadata.SoftlayerCommand) *BandwidthCommand {
	thisCmd := &BandwidthCommand{
		SoftlayerCommand: sl,
		SearchManager:    managers.NewSearchManager(sl.Session),
	}
	cobraCmd := &cobra.Command{

		Use:   "bandwidth",
		Short: T("Bandwidth report for every pool/server."),
		Long: `EXAMPLE:
${COMMAND_NAME} sl report bandwidth`,
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd

}

func (cmd *BandwidthCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	mask := "mask[resource(SoftLayer_Hardware)[id,bandwidthAllocation,bandwidthAllotmentDetail[id,bandwidthAllotment[id,bandwidthAllotmentTypeId,name]],billingItem[id,createDate,lastBillDate],datacenter[id,name],fullyQualifiedDomainName,inboundPublicBandwidthUsage,outboundPublicBandwidthUsage,primaryIpAddress,tagReferences[id,tag[id,name]]],resource(SoftLayer_Network_Application_Delivery_Controller)[id,billingItem[id,bandwidthAllocation[id,amount],bandwidthAllotmentDetail[id,bandwidthAllotment[id,bandwidthAllotmentTypeId,name]],createDate,lastBillDate],datacenter[id,name],name,outboundPublicBandwidthUsage,primaryIpAddress,tagReferences[id,tag[id,name]]],resource(SoftLayer_Virtual_Guest)[id,bandwidthAllocation,bandwidthAllotmentDetail[id,bandwidthAllotment[id,bandwidthAllotmentTypeId,name]],billingItem[id,createDate,lastBillDate],datacenter[id,name],fullyQualifiedDomainName,inboundPublicBandwidthUsage,outboundPublicBandwidthUsage,primaryIpAddress,tagReferences[id,tag[id,name]]]]"

	searchString := "_objectType:SoftLayer_Hardware,SoftLayer_Virtual_Guest,SoftLayer_Network_Application_Delivery_Controller _sort:[fullyQualifiedDomainName:asc]"

	bandwidths, err := cmd.SearchManager.AdvancedSearch(mask, searchString)
	if err != nil {
		return err
	}
	table := cmd.UI.Table([]string{
		T("Id"),
		T("Device name"),
		T("Location"),
		T("Allocation"),
		T("Data in"),
		T("Data out"),
		T("Total usage"),
		T("Pool"),
		T("Tags"),
	})
	for _, bandwidth := range bandwidths {

		resourceJSON, err := json.Marshal(bandwidth.Resource)
		if err != nil {
			fmt.Printf("Error marshalling resource: %v\n", err)
			continue
		}

		typeDataType := reflect.TypeOf(bandwidth.Resource).String()

		if typeDataType == "*datatypes.Hardware" {
			var hardware datatypes.Hardware
			if err = json.Unmarshal(resourceJSON, &hardware); err == nil {
				table.Add(cleanDatas(
					utils.FormatIntPointer(hardware.Id),
					utils.FormatStringPointer(hardware.FullyQualifiedDomainName),
					utils.FormatStringPointer(hardware.Datacenter.Name),
					utils.FormatSLFloatPointerToFloat(hardware.BandwidthAllocation),
					utils.FormatSLFloatPointerToFloat(hardware.InboundPublicBandwidthUsage),
					utils.FormatSLFloatPointerToFloat(hardware.OutboundPublicBandwidthUsage),
					utils.FormatStringPointer(hardware.BandwidthAllotmentDetail.BandwidthAllotment.Name),
					utils.TagRefsToString(hardware.TagReferences),
				)...)
				continue
			}
		}
		if typeDataType == "*datatypes.Virtual_Guest" {
			var virtual datatypes.Virtual_Guest
			if err = json.Unmarshal(resourceJSON, &virtual); err == nil {

				pool := "-"
				if virtual.BandwidthAllotmentDetail != nil {
					pool = utils.FormatStringPointer(virtual.BandwidthAllotmentDetail.BandwidthAllotment.Name)
				}

				table.Add(cleanDatas(
					utils.FormatIntPointer(virtual.Id),
					utils.FormatStringPointer(virtual.FullyQualifiedDomainName),
					utils.FormatStringPointer(virtual.Datacenter.Name),
					utils.FormatSLFloatPointerToFloat(virtual.BandwidthAllocation),
					utils.FormatSLFloatPointerToFloat(virtual.InboundPublicBandwidthUsage),
					utils.FormatSLFloatPointerToFloat(virtual.OutboundPublicBandwidthUsage),
					pool,
					utils.TagRefsToString(virtual.TagReferences),
				)...)
				continue
			}
		}
		if typeDataType == "*datatypes.Network_Application_Delivery_Controller" {
			var delivery datatypes.Network_Application_Delivery_Controller
			if err = json.Unmarshal(resourceJSON, &delivery); err == nil {
				table.Add(cleanDatas(
					utils.FormatIntPointer(delivery.Id),
					utils.FormatStringPointer(delivery.Name),
					utils.FormatStringPointer(delivery.Datacenter.Name),
					"",
					"",
					utils.FormatSLFloatPointerToFloat(delivery.OutboundPublicBandwidthUsage),
					"",
					utils.TagRefsToString(delivery.TagReferences),
				)...)
				continue
			}
		}
	}

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

func cleanDatas(id string, deviceName string, location string, allocation string, dataIn string, dataOut string, pool string, tags string) []string {
	convertedDataIn := utils.ConvertSizes(dataIn, "GB", false)
	convertedDataOut := utils.ConvertSizes(dataOut, "GB", false)
	if allocation == "0.000000" {
		allocation = "Pay-As-You-Go"
	} else {
		if allocation == "-" || allocation == "" {
			allocation = "Unlimited"
		} else {
			allocation = utils.ConvertSizes(allocation, "GB", true)
		}
	}
	if allocation == "Pay-As-You-Go" || allocation == "Unlimited" {
		pool = "Not Applicable"
	}
	var rowBandiwdth []string
	rowBandiwdth = append(rowBandiwdth, id)
	rowBandiwdth = append(rowBandiwdth, deviceName)
	rowBandiwdth = append(rowBandiwdth, location)
	rowBandiwdth = append(rowBandiwdth, allocation)
	rowBandiwdth = append(rowBandiwdth, convertedDataIn)
	rowBandiwdth = append(rowBandiwdth, convertedDataOut)
	rowBandiwdth = append(rowBandiwdth, utils.SumSizes(convertedDataIn, convertedDataOut))
	rowBandiwdth = append(rowBandiwdth, pool)
	rowBandiwdth = append(rowBandiwdth, tags)
	return rowBandiwdth
}
