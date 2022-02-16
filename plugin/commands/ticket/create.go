package ticket

import (
	"fmt"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CreateStandardTicketCommand struct {
	UI            terminal.UI
	TicketManager managers.TicketManager
}

func NewCreateStandardTicketCommand(ui terminal.UI, ticketManager managers.TicketManager) (cmd *CreateStandardTicketCommand) {
	return &CreateStandardTicketCommand{
		UI:            ui,
		TicketManager: ticketManager,
	}
}

func (cmd *CreateStandardTicketCommand) Run(c *cli.Context) error {
	if !c.IsSet("subject-id") {
		return errors.NewInvalidUsageError(T("This command requires the --subject-id option."))
	}

	if !c.IsSet("title") {
		return errors.NewInvalidUsageError(T("This command requires the --title option."))
	}

	var content string
	var err error

	if !c.IsSet("body") {
		content, err = cmd.TicketManager.GetText()
		if err != nil {
			return err
		}

	} else {
		content = c.String("body")
	}

	ticketArgs := managers.TicketArguments{}

	title := c.String("title")
	subjectId := c.Int("subject-id")
	priority := c.Int("priority")

	ticketArgs.Title = &title
	ticketArgs.Content = &content
	ticketArgs.SubjectId = &subjectId
	ticketArgs.Priority = &priority

	if c.IsSet("attachment") {

		attachmentType := c.String("attachment-type")
		attachmentType = strings.ToLower(attachmentType)
		switch attachmentType {
		case "hardware", "":
			HardwareAttachmentType := "HARDWARE"
			ticketArgs.AttachmentType = &HardwareAttachmentType
		case "virtual":
			VirtualAttachmentType := "VIRTUAL_GUEST"
			ticketArgs.AttachmentType = &VirtualAttachmentType
		default:
			return utils.FailWithError(fmt.Sprintf(T("options for %s are hardware or virtual"), "attachment-type"), cmd.UI)
		}

		if !c.IsSet("rootpwd") {
			return cli.NewExitError(T("Root password must be provided with rootpwd flag if attachment is set."), 1)
		}
		id := c.Int("attachment")
		ticketArgs.AttachmentId = &id
		pwd := c.String("rootpwd")
		ticketArgs.RootPassword = &pwd
	}

	ticket_id, err := cmd.TicketManager.CreateStandardTicket(&ticketArgs)
	if err != nil {
		return cli.NewExitError(T("Ticket could not be created: {{.Error}}.", map[string]interface{}{"Error": err.Error()}), 2)
	}

	cmd.UI.Print(T("Ticket ID: {{.TicketID}}.", map[string]interface{}{"TicketID": *ticket_id}))

	return err
}

func TicketCreateMetaData() cli.Command {
	return cli.Command{
		Category:    "ticket",
		Name:        "create",
		Description: T("Create a support ticket"),
		Usage: T(`${COMMAND_NAME} sl ticket create [OPTIONS]

EXAMPLE: 	
    ${COMMAND_NAME} sl ticket create --title "Example title" --subject-id 1522 --body "This is an example ticket. Please disregard."
    ${COMMAND_NAME} sl ticket create --title "Example title" --subject-id 1522 --body "This is an example ticket. Please disregard." --attachment 8675654 --attachment-type hardware --rootpwd passw0rd
    ${COMMAND_NAME} sl ticket create --title "Example title" --subject-id 1522 --body "This is an example ticket. Please disregard." --attachment 1234567 --attachment-type virtual --rootpwd passw0rd
    ${COMMAND_NAME} sl ticket create --title "Example title" --subject-id 1522 --attachment 8675654 --rootpwd passw0rd
    ${COMMAND_NAME} sl ticket create --title "Example title" --subject-id 1522`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "attachment",
				Usage: T("Initial object ID number to attach to ticket"),
			},
			cli.StringFlag{
				Name:  "rootpwd",
				Usage: T("Root password associated with attached device id"),
			},
			cli.IntFlag{
				Name:  "subject-id",
				Usage: T("The subject id to use for the ticket, issue '${COMMAND_NAME} sl ticket subjects' to get the list. [required]"),
			},
			cli.StringFlag{
				Name:  "title",
				Usage: T("The title of the ticket. [required]"),
			},
			cli.StringFlag{
				Name:  "body",
				Usage: T("The ticket body"),
			},
			cli.StringFlag{
				Name:  "priority",
				Usage: T("Ticket priority [1|2|3|4], from 1 (Critical) to 4 (Minimal Impact). Only settable with Advanced and Premium support. See https://www.ibm.com/cloud/support"),
			},
			cli.StringFlag{
				Name:  "attachment-type",
				Usage: T("Specify the type of attachment, hardware or virtual. default is hardware"),
			},
		},
	}
}
