package ticket

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
)

func GetCommandActionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	ticketManager := managers.NewTicketManager(session)
	userManager := managers.NewUserManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{
		"ticket-attach": func(c *cli.Context) error {
			return NewAttachDeviceTicketCommand(ui, ticketManager).Run(c)
		},
		"ticket-create": func(c *cli.Context) error {
			return NewCreateStandardTicketCommand(ui, ticketManager).Run(c)
		},
		"ticket-detach": func(c *cli.Context) error {
			return NewDetachDeviceTicketCommand(ui, ticketManager).Run(c)
		},
		"ticket-detail": func(c *cli.Context) error {
			return NewDetailTicketCommand(ui, ticketManager, userManager).Run(c)
		},
		"ticket-list": func(c *cli.Context) error {
			return NewListTicketCommand(ui, ticketManager).Run(c)
		},
		"ticket-subjects": func(c *cli.Context) error {
			return NewSubjectsTicketCommand(ui, ticketManager).Run(c)
		},
		"ticket-summary": func(c *cli.Context) error {
			return NewSummaryTicketCommand(ui, ticketManager).Run(c)
		},
		"ticket-update": func(c *cli.Context) error {
			return NewUpdateTicketCommand(ui, ticketManager).Run(c)
		},
		"ticket-upload": func(c *cli.Context) error {
			return NewUploadFileTicketCommand(ui, ticketManager).Run(c)
		},
	}

	return CommandActionBindings
}

func TicketNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  "sl",
		Name:        "ticket",
		Description: T("Classic infrastructure Manage Tickets"),
	}
}

func TicketMetaData() cli.Command {
	return cli.Command{
		Category:    "sl",
		Name:        "ticket",
		Usage:       "${COMMAND_NAME} sl ticket",
		Description: T("Classic infrastructure Manage Tickets"),
		Subcommands: []cli.Command{
			TicketCreateMetaData(),
			TicketDetailMetaData(),
			TicketAttachMetaData(),
			TicketDetachMetaData(),
			TicketSubjectsMetaData(),
			TicketUpdataMetaData(),
			TicketListMetaData(),
			TicketUploadMetaData(),
			TicketSummaryMetaData(),
		},
	}
}
