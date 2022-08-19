package placementgroup

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type PlacementGroupCreateCommand struct {
	*metadata.SoftlayerCommand
	PlaceGroupManager managers.PlaceGroupManager
	Command           *cobra.Command
	Name              string
	BackendRouterId   int
	RuleId            int
}

func NewPlacementGroupCreateCommand(sl *metadata.SoftlayerCommand) (cmd *PlacementGroupCreateCommand) {
	thisCmd := &PlacementGroupCreateCommand{
		SoftlayerCommand:  sl,
		PlaceGroupManager: managers.NewPlaceGroupManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "create",
		Short: T("Create a placement group"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVarP(&thisCmd.Name, "name", "n", "", T("Name for this new placement group. [required]"))
	cobraCmd.Flags().IntVarP(&thisCmd.BackendRouterId, "backend-router-id", "b", 0, T("Backend router ID. [required]"))
	cobraCmd.Flags().IntVarP(&thisCmd.RuleId, "rule-id", "r", 0, T("The ID of the rule to govern this placement group. [required]"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *PlacementGroupCreateCommand) Run(args []string) error {
	name := cmd.Name
	if name == "" {
		return errors.NewMissingInputError("-n, --name")

	}
	backendRouter := cmd.BackendRouterId
	if backendRouter == 0 {
		return errors.NewMissingInputError("-b, --backend-router-id")
	}
	rule := cmd.RuleId
	if rule == 0 {
		return errors.NewMissingInputError("-r, --rule-id")
	}

	outputFormat := cmd.GetOutputFlag()

	placementObject := datatypes.Virtual_PlacementGroup{
		Name:            &name,
		BackendRouterId: &backendRouter,
		RuleId:          &rule,
	}
	result, err := cmd.PlaceGroupManager.Create(&placementObject)
	if err != nil {
		return errors.NewAPIError(T("Failed to create placement group"), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, result)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Successfully created placement group: ID: {{.ID}}, Name: {{.Name}}.", map[string]interface{}{"ID": result.Id, "Name": result.Name}))
	return nil
}
