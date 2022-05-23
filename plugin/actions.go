package plugin

import (
	"fmt"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/licenses"

	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration/core_config"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/account"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/autoscale"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/callapi"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dedicatedhost"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dns"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/email"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/eventlog"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/file"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/globalip"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/image"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/ipsec"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/loadbal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/order"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/placementgroup"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/reports"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/security"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/securitygroup"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/subnet"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/tags"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/ticket"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/user"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/vlan"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

func GetCommandAcionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {

	CommandActionBindings := map[string]func(c *cli.Context) error{}

	// ibmcloud sl account
	accountCommands := account.GetCommandActionBindings(context, ui, session)
	for name, action := range accountCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl autoscale
	autoScaleCommands := autoscale.GetCommandActionBindings(context, ui, session)
	for name, action := range autoScaleCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl dedicatedhost
	dedicatedhostCommands := dedicatedhost.GetCommandActionBindings(context, ui, session)
	for name, action := range dedicatedhostCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl dns
	dnsCommands := dns.GetCommandActionBindings(context, ui, session)
	for name, action := range dnsCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl dns
	emailCommands := email.GetCommandActionBindings(context, ui, session)
	for name, action := range emailCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl vlan
	vlanCommands := vlan.GetCommandActionBindings(context, ui, session)
	for name, action := range vlanCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl block
	blockCommands := block.GetCommandAcionBindings(context, ui, session)
	for name, action := range blockCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl callapi
	callapiCommands := callapi.GetCommandActionBindings(context, ui, session)
	for name, action := range callapiCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl event-log
	eventLogCommands := eventlog.GetCommandActionBindings(context, ui, session)
	for name, action := range eventLogCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl file
	fileCommands := file.GetCommandAcionBindings(context, ui, session)
	for name, action := range fileCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl hardware
	hardwareCommands := hardware.GetCommandActionBindings(context, ui, session)
	for name, action := range hardwareCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl image
	imageCommands := image.GetCommandActionBindings(context, ui, session)
	for name, action := range imageCommands {
		CommandActionBindings[name] = action
	}
	// ibmcloud sl ipsec
	ipsecCommands := ipsec.GetCommandActionBindings(context, ui, session)
	for name, action := range ipsecCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl tags
	tagsCommands := tags.GetCommandActionBindings(context, ui, session)
	for name, action := range tagsCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl loadbal
	loadbalCommands := loadbal.GetCommandActionBindings(context, ui, session)
	for name, action := range loadbalCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl vs
	vsCommands := virtual.GetCommandActionBindings(context, ui, session)
	for name, action := range vsCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl order
	orderCommands := order.GetCommandActionBindings(context, ui, session)
	for name, action := range orderCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl placement-group
	placementgroupCommands := placementgroup.GetCommandActionBindings(context, ui, session)
	for name, action := range placementgroupCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl globalip
	globalipCommands := globalip.GetCommandActionBindings(context, ui, session)
	for name, action := range globalipCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl subnet
	subnetCommands := subnet.GetCommandActionBindings(context, ui, session)
	for name, action := range subnetCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl ticket
	ticketCommands := ticket.GetCommandActionBindings(context, ui, session)
	for name, action := range ticketCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl securitygroup
	securitygroupCommands := securitygroup.GetCommandActionBindings(context, ui, session)
	for name, action := range securitygroupCommands {
		CommandActionBindings[name] = action
	}

	userCommands := user.GetCommandActionBindings(context, ui, session)
	for name, action := range userCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl security
	securityCommands := security.GetCommandActionBindings(context, ui, session)
	for name, action := range securityCommands {
		CommandActionBindings[name] = action
	}

	//ibmcloud sl licenses
	licenseCommands := licenses.GetCommandActionBindings(context, ui, session)
	for name, action := range licenseCommands {
		CommandActionBindings[name] = action
	}

	//ibmcloud sl reports
	for name, action := range reports.GetCommandActionBindings(context, ui, session) {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl firewall
	// Deprecated for now.
	/*
		firewallCommands := firewall.GetCommandActionBindings(context, ui, session)
		for name, action := range firewallCommands {
			CommandActionBindings[name] = action
		}
	*/

	actionWithPreCheck := make(map[string]func(c *cli.Context) error)

	for name, action := range CommandActionBindings {
		actionCopy := action
		actionWithPreCheck[name] = func(c *cli.Context) (err error) {
			err = PreChecktRequirement(context, ui)
			if err != nil {
				return err
			}

			defer func() {
				// catch panic
				if recoverErr := recover(); recoverErr != nil {
					err = cli.NewExitError(fmt.Sprintf("%v", recoverErr), 1)
				}
				switch err.(type) {
				case *errors.InvalidUsageError:
					ui.Failed("%v", err)
					showCmdErr := cli.ShowCommandHelp(c, c.Command.Name)
					if showCmdErr != nil {
						fmt.Println(showCmdErr.Error())
					}
					err = cli.NewExitError("", 2)
				}
			}()
			err = actionCopy(c)
			return err
		}
	}
	return actionWithPreCheck
}

func PreChecktRequirement(context plugin.PluginContext, ui terminal.UI) error {
	var errorMessage error
	switch {
	case !context.IsLoggedIn():
		errorMessage = fmt.Errorf(T("Not logged in. Use '{{.Command}}' to log in.",
			map[string]interface{}{"Command": terminal.CommandColor(context.CLIName() + " login")}))
	case context.IAMToken() == "":
		errorMessage = fmt.Errorf(T("IAM token is required. Use '{{.Command}}' to log in.",
			map[string]interface{}{"Command": terminal.CommandColor(context.CLIName() + " login")}))
	case context.IMSAccountID() == "":
		errorMessage = fmt.Errorf(T("Current account is not linked to a Softlayer account. Use '{{.Command}}' to switch account.",
			map[string]interface{}{"Command": terminal.CommandColor(context.CLIName() + " target -c")}))
	case !core_config.NewIAMTokenInfo(context.IAMToken()).Accounts.Valid:
		errorMessage = fmt.Errorf(T("The linked Softlayer account is not validated. Use '{{.Command}}' to re-login.",
			map[string]interface{}{"Command": terminal.CommandColor(context.CLIName() + " login")}))
	}
	if errorMessage != nil {
		return utils.FailWithError(errorMessage.Error(), ui)
	}
	return nil
}
