package account

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type BandwidthPoolsDetailCommand struct {
	*metadata.SoftlayerCommand
	AccountManager managers.AccountManager
	Command        *cobra.Command
}

func NewBandwidthPoolsDetailCommand(sl *metadata.SoftlayerCommand) *BandwidthPoolsDetailCommand {

	thisCmd := &BandwidthPoolsDetailCommand{
		SoftlayerCommand: sl,
		AccountManager:   managers.NewAccountManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "bandwidth-pools-detail " + T("IDENTIFIER"),
		Short: T("Get bandwidth pool details."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *BandwidthPoolsDetailCommand) Run(args []string) error {

	outputFormat := cmd.GetOutputFlag()

	bandwidthPoolId, err := strconv.Atoi(args[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Bandwidth Pool ID")
	}

	bandwidthPool, err := cmd.AccountManager.GetBandwidthPoolDetail(bandwidthPoolId, "")
	if err != nil {
		return errors.NewAPIError(T("Failed to get Bandwidth Pool."), err.Error(), 2)
	}

	currentUsage := "-"
	if bandwidthPool.BillingCyclePublicBandwidthUsage != nil {
		if bandwidthPool.BillingCyclePublicBandwidthUsage.AmountOut != nil {
			currentUsage = fmt.Sprintf("%.2f GB", *bandwidthPool.BillingCyclePublicBandwidthUsage.AmountOut/1000000000.0)
		}
	}
	projectedUsage := "-"
	if bandwidthPool.ProjectedPublicBandwidthUsage != nil {
		projectedUsage = fmt.Sprintf("%.2f GB", *bandwidthPool.ProjectedPublicBandwidthUsage/1000000000.0)
	}
	inboundUsage := "-"
	if bandwidthPool.InboundPublicBandwidthUsage != nil {
		inboundUsage = fmt.Sprintf("%.2f GB", *bandwidthPool.InboundPublicBandwidthUsage/1000000000.0)
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("ID"), utils.FormatIntPointer(bandwidthPool.Id))
	table.Add(T("Name"), utils.FormatStringPointer(bandwidthPool.Name))
	table.Add(T("Create Date"), utils.FormatSLTimePointer(bandwidthPool.CreateDate))
	table.Add(T("Current Usage"), currentUsage)
	table.Add(T("Projected Usage"), projectedUsage)
	table.Add(T("Inbound Usage"), inboundUsage)

	if bandwidthPool.Hardware != nil && len(bandwidthPool.Hardware) != 0 {
		hardwareTableBuffer := getHardwareTable(bandwidthPool.Hardware)
		table.Add(T("Hardware"), hardwareTableBuffer.String())
	} else {
		table.Add(T("Hardware"), T("Not Found"))
	}

	if bandwidthPool.VirtualGuests != nil && len(bandwidthPool.VirtualGuests) != 0 {
		virtualTableBuffer := getVirtualTable(bandwidthPool.VirtualGuests)
		table.Add(T("Virtual"), virtualTableBuffer.String())
	} else {
		table.Add(T("Virtual"), T("Not Found"))
	}

	if bandwidthPool.BareMetalInstances != nil && len(bandwidthPool.BareMetalInstances) != 0 {
		bareMetalInstancesTableBuffer := getHardwareTable(bandwidthPool.BareMetalInstances)
		table.Add(T("Netscaler"), bareMetalInstancesTableBuffer.String())
	} else {
		table.Add(T("Netscaler"), T("Not Found"))
	}

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

func getHardwareTable(hardwares []datatypes.Hardware) *bytes.Buffer {
	buf := new(bytes.Buffer)
	hardwareTable := terminal.NewTable(buf, []string{T("Id"), T("Hostname"), T("IP Address"), T("Amount"), T("Current Usage")})
	for _, hardware := range hardwares {
		ipAddress := "-"
		if hardware.PrimaryIpAddress != nil {
			ipAddress = *hardware.PrimaryIpAddress
		}
		amount := "-"
		if hardware.BandwidthAllotmentDetail != nil && hardware.BandwidthAllotmentDetail.Allocation != nil && hardware.BandwidthAllotmentDetail.Allocation.Amount != nil {
			amount = fmt.Sprintf("%.2f GB", float64(*hardware.BandwidthAllotmentDetail.Allocation.Amount)/1000000000)
		}
		current := "-"
		if hardware.OutboundBandwidthUsage != nil {
			current = fmt.Sprintf("%.2f GB", float64(*hardware.OutboundBandwidthUsage)/1000000000)
		}
		hardwareTable.Add(
			utils.FormatIntPointer(hardware.Id),
			utils.FormatStringPointer(hardware.FullyQualifiedDomainName),
			ipAddress,
			amount,
			current,
		)
	}
	hardwareTable.Print()
	return buf
}

func getVirtualTable(virtuals []datatypes.Virtual_Guest) *bytes.Buffer {
	buf := new(bytes.Buffer)
	virtualTable := terminal.NewTable(buf, []string{T("Id"), T("Hostname"), T("IP Address"), T("Amount"), T("Current Usage")})
	for _, virtual := range virtuals {
		ipAddress := "-"
		if virtual.PrimaryIpAddress != nil {
			ipAddress = *virtual.PrimaryIpAddress
		}
		amount := "-"
		if virtual.BandwidthAllotmentDetail != nil && virtual.BandwidthAllotmentDetail.Allocation != nil && virtual.BandwidthAllotmentDetail.Allocation.Amount != nil {
			amount = fmt.Sprintf("%.2f GB", float64(*virtual.BandwidthAllotmentDetail.Allocation.Amount)/1000000000)
		}
		current := "-"
		if virtual.OutboundPublicBandwidthUsage != nil {
			current = fmt.Sprintf("%.2f GB", float64(*virtual.OutboundPublicBandwidthUsage)/1000000000)
		}
		virtualTable.Add(
			utils.FormatIntPointer(virtual.Id),
			utils.FormatStringPointer(virtual.FullyQualifiedDomainName),
			ipAddress,
			amount,
			current,
		)
	}
	virtualTable.Print()
	return buf
}
