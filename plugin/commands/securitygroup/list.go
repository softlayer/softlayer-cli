package securitygroup

import (
	"sort"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	UI             terminal.UI
	NetworkManager managers.NetworkManager
}

func NewListCommand(ui terminal.UI, networkManager managers.NetworkManager) (cmd *ListCommand) {
	return &ListCommand{
		UI:             ui,
		NetworkManager: networkManager,
	}
}

func (cmd *ListCommand) Run(c *cli.Context) error {
	var sortby string
	if c.IsSet("sortby") {
		sortby = strings.ToLower(c.String("sortby"))
		if sortby != "id" && sortby != "name" && sortby != "description" && sortby != "created" {
			return errors.NewInvalidUsageError(T("Options for --sortby are: id,name,description,created"))
		}
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	groups, err := cmd.NetworkManager.ListSecurityGroups()
	if err != nil {
		return cli.NewExitError(T("Failed to get security groups.\n")+err.Error(), 2)
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
