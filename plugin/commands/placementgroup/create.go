package placementgroup

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type PlacementGroupCreateCommand struct {
	UI                terminal.UI
	PlaceGroupManager managers.PlaceGroupManager
}

func NewPlacementGroupCreateCommand(ui terminal.UI, placeGroupManager managers.PlaceGroupManager) (cmd *PlacementGroupCreateCommand) {
	return &PlacementGroupCreateCommand{
		UI:                ui,
		PlaceGroupManager: placeGroupManager,
	}
}

func (cmd *PlacementGroupCreateCommand) Run(c *cli.Context) error {
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
	result, err := cmd.PlaceGroupManager.Create(&placementObject)
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