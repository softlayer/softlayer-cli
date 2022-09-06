package virtual

import (
	"bytes"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CapacityDetailCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
	Column               []string
	SortBy               string
}

func NewCapacityDetailCommand(sl *metadata.SoftlayerCommand) (cmd *CapacityDetailCommand) {
	thisCmd := &CapacityDetailCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "capacity-detail " + T("IDENTIFIER"),
		Short: T("Get Reserved Capacity Group details."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	cobraCmd.Flags().StringVar(&thisCmd.SortBy, "sortby", "id", T("Column to sort by. Options are: id, hostname, domain, primary_ip, backend_ip"))
	cobraCmd.Flags().StringSliceVar(&thisCmd.Column, "column", []string{}, T("Column to display. Options are: id, hostname, domain, primary_ip, backend_ip. This option can be specified multiple times"))
	return thisCmd
}

func (cmd *CapacityDetailCommand) Run(args []string) error {

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Reserved Capacity Group Virtual server ID")
	}
	capacity, err := cmd.VirtualServerManager.GetCapacityDetail(id)
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Reserved Capacity Gruop Virtual server ID")
	}

	outputFormat := cmd.GetOutputFlag()
	sortby := cmd.SortBy
	if sortby == "" {
		sortby = "hostname"
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, capacity)
	}

	defaultColumns := []string{"id", "hostname", "domain", "primary_id", "backend_id"}
	optionalColumns := []string{"id", "hostname", "domain", "primary_id", "backend_id"}
	sortColumns := []string{"id", "hostname", "domain", "primary_id", "backend_id"}

	showColumns, err := utils.ValidateColumns2(sortby, cmd.Column, defaultColumns, optionalColumns, sortColumns)
	if err != nil {
		return err
	}

	mainTable := cmd.UI.Table([]string{T("detail")})
	mainTable.Add(utils.FormatStringPointer(capacity.Name))
	buf := new(bytes.Buffer)
	table := terminal.NewTable(buf, utils.GetColumnHeader(showColumns))
	for _, instance := range capacity.Instances {
		values := make(map[string]string)
		if instance.Guest != nil {
			values["id"] = utils.FormatIntPointer(instance.Id)
			values["hostname"] = utils.FormatStringPointer(instance.Guest.Hostname)
			values["domain"] = utils.FormatStringPointer(instance.Guest.Domain)
			values["primary_id"] = utils.FormatStringPointer(instance.Guest.PrimaryIpAddress)
			values["backend_id"] = utils.FormatStringPointer(instance.Guest.PrimaryBackendIpAddress)
		} else {
			values["id"] = utils.EMPTY_VALUE
			values["hostname"] = utils.EMPTY_VALUE
			values["domain"] = utils.EMPTY_VALUE
			values["primary_id"] = utils.EMPTY_VALUE
			values["backend_id"] = utils.EMPTY_VALUE
		}
		row := make([]string, len(showColumns))
		for i, col := range showColumns {
			row[i] = values[col]
		}
		table.Add(row...)
	}
	table.Print()
	mainTable.Add(buf.String())
	mainTable.Print()
	return nil
}
