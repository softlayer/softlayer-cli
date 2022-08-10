package block

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type SubnetsRemoveCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewSubnetsRemoveCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *SubnetsRemoveCommand) {
	return &SubnetsRemoveCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func BlockSubnetsRemoveMetaData() cli.Command {
	return cli.Command{
		Category:    "block",
		Name:        "subnets-remove",
		Description: T("Remove block storage subnets to the given host id."),
		Usage: T(`${COMMAND_NAME} sl block subnets-remove ACCESS_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl block subnets-remove 111111 --subnet-id 222222
   ${COMMAND_NAME} sl block subnets-remove 111111 --subnet-id 222222 --subnet-id 333333
   ACCESS_ID is the host_id obtained by: ibmcloud sl block access-list <volume_id>`),
		Flags: []cli.Flag{
			cli.IntSliceFlag{
				Name:     "subnet-id",
				Usage:    T("IDs of the subnets to remove"),
				Required: true,
			},
		},
	}
}

func (cmd *SubnetsRemoveCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	accessID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Access ID")
	}

	subnetsToRemove := c.IntSlice("subnet-id")

	subnetsResponse, err := cmd.StorageManager.RemoveSubnetsFromAcl(accessID, subnetsToRemove)
	if err != nil {
		return cli.NewExitError(T("Failed to remove subnets.")+"\n"+err.Error(), 2)
	}

	for _, subnet := range subnetsToRemove {
		values := map[string]interface{}{"subnetID": subnet, "accessID": accessID}
		if utils.IntInSlice(subnet, subnetsResponse) != -1 {
			cmd.UI.Print(T("Successfully removed subnet id: {{.subnetID}} to allowed host id: {{.accessID}}", values))
		} else {
			cmd.UI.Print(T("Failed to remove subnet id: {{.subnetID}} to allowed host id: {{.accessID}}", values))
		}
	}

	return nil
}
