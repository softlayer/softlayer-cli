package virtual

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type PlacementGroupCreateCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
	Name                 string
	BackendRouterId      int
	RuleId               int
}

func NewPlacementGroupCreateCommand(sl *metadata.SoftlayerCommand) (cmd *PlacementGroupCreateCommand) {
	thisCmd := &PlacementGroupCreateCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "placementgroup-create",
		Short: T("Create a placement group"),
		Long: T(`${COMMAND_NAME} sl vs placementgroup-create [OPTIONS]
EXAMPLE:
${COMMAND_NAME} sl vs placementgroup-create -n myvsi -b 1234567 -r 258369 
This command orders a Placement group instance with name is myvsi, backendRouterId 1234567, and rule 258369`),
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	cobraCmd.Flags().StringVarP(&thisCmd.Name, "name", "n", "", T("Name for this new placement group. [required]"))
	cobraCmd.Flags().IntVarP(&thisCmd.BackendRouterId, "backend-router-id", "b", 0, T("Backend router ID. [required]"))
	cobraCmd.Flags().IntVarP(&thisCmd.RuleId, "rule-id", "r", 0, T("The ID of the rule to govern this placement group. [required]"))
	return thisCmd
}

func (cmd *PlacementGroupCreateCommand) Run(args []string) error {
	name := cmd.Name
	if name == "" {
		return slErrors.NewMissingInputError("-n, --name")

	}
	backendRouter := cmd.BackendRouterId
	if backendRouter == 0 {
		return slErrors.NewMissingInputError("-b, --backend-router-id")
	}
	rule := cmd.RuleId
	if rule == 0 {
		return slErrors.NewMissingInputError("-r, --rule-id")
	}

	outputFormat := cmd.GetOutputFlag()

	placementObject := datatypes.Virtual_PlacementGroup{
		Name:            &name,
		BackendRouterId: &backendRouter,
		RuleId:          &rule,
	}
	result, err := cmd.VirtualServerManager.PlacementCreate(&placementObject)
	if err != nil {
		return slErrors.NewAPIError(T("Failed to create placement group\n"), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, result)
	}

	cmd.UI.Ok()
	cmd.UI.Print(T("Successfully created placement group: ID: {{.ID}}, Name: {{.Name}}.", map[string]interface{}{"ID": result.Id, "Name": result.Name}))
	return nil
}
