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

type SubnetsAssignCommand struct {
	UI             terminal.UI
	StorageManager managers.StorageManager
}

func NewSubnetsAssignCommand(ui terminal.UI, storageManager managers.StorageManager) (cmd *SubnetsAssignCommand) {
	return &SubnetsAssignCommand{
		UI:             ui,
		StorageManager: storageManager,
	}
}

func BlockSubnetsAssignMetaData() cli.Command {
	return cli.Command{
		Category:    "block",
		Name:        "subnets-assign",
		Description: T("Assign block storage subnets to the given host id."),
		Usage: T(`${COMMAND_NAME} sl block subnets-assign ACCESS_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl block subnets-assign 111111 --subnet-id 222222
   ${COMMAND_NAME} sl block subnets-assign 111111 --subnet-id 222222 --subnet-id 333333
   ACCESS_ID is the host_id obtained by: ibmcloud sl block access-list <volume_id>`),
		Flags: []cli.Flag{
			cli.IntSliceFlag{
				Name:     "subnet-id",
				Usage:    T("IDs of the subnets to assign"),
				Required: true,
			},
		},
	}
}

func (cmd *SubnetsAssignCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	accessID, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Access ID")
	}

	subnetsToAssign := c.IntSlice("subnet-id")

	subnetsResponse, err := cmd.StorageManager.AssignSubnetsToAcl(accessID, subnetsToAssign)
	if err != nil {
		return cli.NewExitError(T("Failed to assign subnets.")+"\n"+err.Error(), 2)
	}

	for _, subnet := range subnetsToAssign {
		values := map[string]interface{}{"subnetID": subnet, "accessID": accessID}
		if utils.IntInSlice(subnet, subnetsResponse) != -1 {
			cmd.UI.Print(T("Successfully assigned subnet id: {{.subnetID}} to allowed host id: {{.accessID}}", values))
		} else {
			cmd.UI.Print(T("Failed to assign subnet id: {{.subnetID}} to allowed host id: {{.accessID}}", values))
		}
	}

	return nil
}
