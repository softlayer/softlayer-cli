package securitygroup

import (
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	*metadata.SoftlayerCommand
	NetworkManager managers.NetworkManager
	Command        *cobra.Command
	Sortby         string
}

func NewListCommand(sl *metadata.SoftlayerCommand) (cmd *ListCommand) {
	thisCmd := &ListCommand{
		SoftlayerCommand: sl,
		NetworkManager:   managers.NewNetworkManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "list",
		Short: T("List security groups"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVar(&thisCmd.Sortby, "sortby", "", T("Column to sort by. Options are: id,name,description,created"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ListCommand) Run(args []string) error {
	var sortby string
	if cmd.Sortby != "" {
		sortby = strings.ToLower(cmd.Sortby)
		if sortby != "id" && sortby != "name" && sortby != "description" && sortby != "created" {
			return errors.NewInvalidUsageError(T("Options for --sortby are: id,name,description,created"))
		}
	}

	outputFormat := cmd.GetOutputFlag()

	groups, err := cmd.NetworkManager.ListSecurityGroups()
	if err != nil {
		return errors.NewAPIError(T("Failed to get security groups.\n"), err.Error(), 2)
	}

	if sortby == "" || sortby == "id" {
		sort.Sort(utils.GroupById(groups))
	} else if sortby == "name" {
		sort.Sort(utils.GroupByName(groups))
	} else if sortby == "description" {
		sort.Sort(utils.GroupByDescription(groups))
	} else if sortby == "created" {
		sort.Sort(utils.GroupByCreated(groups))
	} else {
		return errors.NewInvalidUsageError(T("Options for --sortby are: id,name,description,created"))
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, groups)
	}

	if len(groups) == 0 {
		cmd.UI.Print(T("No security groups are found."))
		return nil
	}

	table := cmd.UI.Table([]string{T("ID"), T("Name"), T("Description"), T("Created")})
	for _, g := range groups {
		table.Add(utils.FormatIntPointer(g.Id),
			utils.FormatStringPointer(g.Name),
			utils.FormatStringPointer(g.Description),
			utils.FormatSLTimePointer(g.CreateDate))
	}
	table.Print()
	return nil
}
