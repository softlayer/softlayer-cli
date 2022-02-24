package virtual

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type PlacementgroupCreateCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
	Context              plugin.PluginContext
}

func NewVSPlacementGroupCreateCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager, context plugin.PluginContext) (cmd *PlacementgroupCreateCommand) {
	return &PlacementgroupCreateCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
		Context:              context,
	}
}

func (cmd *PlacementgroupCreateCommand) Run(c *cli.Context) error {
	name := c.String("n")
	if name == "" {
		return errors.NewMissingInputError("-n, --name")

	}
	backendRouter := c.Int("b")
	if backendRouter == 0 {
		return errors.NewMissingInputError("-b, --backend-router-id")
	}
	rule := c.Int("r")
	if rule == 0 {
		return errors.NewMissingInputError("-r, --rule-id")
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	placementObject := datatypes.Virtual_PlacementGroup{
		Name:            &name,
		BackendRouterId: &backendRouter,
		RuleId:          &rule,
	}
	result, err := cmd.VirtualServerManager.PlacementCreate(&placementObject)
	if err != nil {
		return cli.NewExitError(T("Failed to create placement group\n")+err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, result)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Successfully created placement group: ID: {{.ID}}, Name: {{.Name}}.", map[string]interface{}{"ID": result.Id, "Name": result.Name}))
	return nil
}

func VSPlacementGroupCreateMetaData() cli.Command {
	return cli.Command{
		Category:    "vs",
		Name:        "placementgroup-create",
		Description: T("Create a placement group."),
		Usage: T(`${COMMAND_NAME} sl vs placementgroup-create [OPTIONS]
EXAMPLE:
${COMMAND_NAME} sl vs placementgroup-create -n myvsi -b 1234567 -r 258369 
This command orders a Placement group instance with name is myvsi, backendRouterId 1234567, and rule 258369`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "n,name",
				Usage: T("Name for this new placement group. [required]"),
			},
			cli.IntFlag{
				Name:  "b,backend-router-id",
				Usage: T("Backend router ID. [required]"),
			},
			cli.IntFlag{
				Name:  "r,rule-id",
				Usage: T("The ID of the rule to govern this placement group. [required]"),
			},
			metadata.OutputFlag(),
		},
	}
}
