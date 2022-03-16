package loadbal

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

type MembersDelCommand struct {
	UI                  terminal.UI
	LoadBalancerManager managers.LoadBalancerManager
}

func NewMembersDelCommand(ui terminal.UI, lbManager managers.LoadBalancerManager) (cmd *MembersDelCommand) {
	return &MembersDelCommand{
		UI:                  ui,
		LoadBalancerManager: lbManager,
	}
}

func (cmd *MembersDelCommand) Run(c *cli.Context) error {
	loadbalID := c.Int("lb-id")
	if loadbalID == 0 {
		return errors.NewMissingInputError("--lb-id")
	}

	memberUUID := c.String("m")
	if memberUUID == "" {
		return errors.NewMissingInputError("-m, --member-uuid")
	}

	loadbalancerUUID, err := cmd.LoadBalancerManager.GetLoadBalancerUUID(loadbalID)
	if err != nil {
		return cli.NewExitError(T("Failed to get load balancer: {{.ERR}}.", map[string]interface{}{"ERR": err.Error()}), 2)
	}

	if !c.IsSet("f") {
		confirm, err := cmd.UI.Confirm(T("This will delete the load balancer member: {{.MemberID}} and cannot be undone. Continue?", map[string]interface{}{"MemberID": memberUUID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Say(T("Aborted."))
			return nil
		}
	}

	_, err = cmd.LoadBalancerManager.DeleteLoadBalancerMember(&loadbalancerUUID, []string{memberUUID})
	if err != nil {
		return cli.NewExitError(T("Failed to delete load balancer member {{.Member}}: {{.Error}}.\n",
			map[string]interface{}{"Member": memberUUID, "Error": err.Error()}), 2)
	}
	cmd.UI.Ok()
	cmd.UI.Say(T("Member {{.MemberID}} removed", map[string]interface{}{"MemberID": memberUUID}))
	return nil
}

func LoadbalMemberDelMetadata() cli.Command {
	return cli.Command{
		Category:    "loadbal",
		Name:        "member-delete",
		Description: T("Remove a load balancer member"),
		Usage:       "${COMMAND_NAME} sl loadbal member-del (--lb-id LOADBAL_ID) (-m, --member-uuid MEMBER_UUID)",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "lb-id",
				Usage: T("ID for the load balancer [required]"),
			},
			cli.StringFlag{
				Name:  "m,member-uuid",
				Usage: T("Member UUID [required]"),
			},
		},
	}
}
