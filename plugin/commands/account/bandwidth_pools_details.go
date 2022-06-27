package account

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type BandwidthPoolsDetailCommand struct {
	UI             terminal.UI
	AccountManager managers.AccountManager
}

func NewBandwidthPoolsDetailCommand(ui terminal.UI, accountManager managers.AccountManager) (cmd *BandwidthPoolsDetailCommand) {
	return &BandwidthPoolsDetailCommand{
		UI:             ui,
		AccountManager: accountManager,
	}
}

func BandwidthPoolsDetailMetaData() cli.Command {
	return cli.Command{
		Category:    "account",
		Name:        "bandwidth-pools-detail",
		Description: T("Get bandwidth pool details."),
		Usage: T(`${COMMAND_NAME} sl account bandwidth-pools-detail
EXAMPLE: 
	${COMMAND_NAME} sl account bandwidth-pools-detail 123456`),

		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}

func (cmd *BandwidthPoolsDetailCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	bandwidthPoolId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return errors.NewInvalidSoftlayerIdInputError("Bandwidth Pool ID")
	}

	bandwidthPool, err := cmd.AccountManager.GetBandwidthPoolDetail(bandwidthPoolId, "")
	if err != nil {
		return cli.NewExitError(T("Failed to get Bandwidth Pool.\n")+err.Error(), 2)
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
	hardwareTable := terminal.NewTable(buf, []string{T("Id"), T("HostName"), T("IP Address"), T("Amount"), T("Current Usage")})
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
	virtualTable := terminal.NewTable(buf, []string{T("Id"), T("HostName"), T("IP Address"), T("Amount"), T("Current Usage")})
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
