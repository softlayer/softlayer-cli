package block

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type SubnetsListCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewSubnetsListCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *SubnetsListCommand) {
	return &SubnetsListCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func BlockSubnetsListMetaData() cli.Command {
	return cli.Command{
		Category:    "block",
		Name:        "subnets-list",
		Description: T("List block storage assigned subnets for the given host id."),
		Usage: T(`${COMMAND_NAME} sl block subnets-list ACCESS_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl block subnets-list 12345678 
   ACCESS_ID is the host_id obtained by: ibmcloud sl block access-list <volume_id>`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}

func (cmd *SubnetsListCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	accessID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Access ID")
	}

	subnets, err := cmd.StorageManager.GetSubnetsInAcl(accessID, "")
	if err != nil {
		return cli.NewExitError(T("Failed to get subnets.")+"\n"+err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("Id"), T("Network Identifier"), T("CIDR")})
	for _, subnet := range subnets {
		table.Add(
			utils.FormatIntPointer(subnet.Id),
			utils.FormatStringPointer(subnet.NetworkIdentifier),
			utils.FormatIntPointer(subnet.Cidr),
		)
	}

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}
