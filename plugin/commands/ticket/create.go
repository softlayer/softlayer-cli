package ticket

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CreateStandardTicketCommand struct {
	*metadata.SoftlayerCommand
	TicketManager  managers.TicketManager
	Command        *cobra.Command
	Attachment     int
	RootPwd        string
	SubjectId      int
	Title          string
	Body           string
	Priority       int
	AttachmentType string
}

func NewCreateStandardTicketCommand(sl *metadata.SoftlayerCommand) *CreateStandardTicketCommand {
	thisCmd := &CreateStandardTicketCommand{
		SoftlayerCommand: sl,
		TicketManager:    managers.NewTicketManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "create",
		Short: T("Create a support ticket"),
		Long: T(`${COMMAND_NAME} sl ticket create [OPTIONS]

EXAMPLE: 	
    ${COMMAND_NAME} sl ticket create --title "Example title" --subject-id 1522 --body "This is an example ticket. Please disregard."
    ${COMMAND_NAME} sl ticket create --title "Example title" --subject-id 1522 --body "This is an example ticket. Please disregard." --attachment 8675654 --attachment-type hardware --rootpwd passw0rd
    ${COMMAND_NAME} sl ticket create --title "Example title" --subject-id 1522 --body "This is an example ticket. Please disregard." --attachment 1234567 --attachment-type virtual --rootpwd passw0rd
    ${COMMAND_NAME} sl ticket create --title "Example title" --subject-id 1522 --attachment 8675654 --rootpwd passw0rd
    ${COMMAND_NAME} sl ticket create --title "Example title" --subject-id 1522`),
		Args: metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntVar(&thisCmd.Attachment, "attachment", 0, T("Initial object ID number to attach to ticket"))
	cobraCmd.Flags().StringVar(&thisCmd.RootPwd, "rootpwd", "", T("Root password associated with attached device id"))
	cobraCmd.Flags().IntVar(&thisCmd.SubjectId, "subject-id", 0, T("The subject id to use for the ticket, issue '${COMMAND_NAME} sl ticket subjects' to get the list. [required]"))
	cobraCmd.Flags().StringVar(&thisCmd.Title, "title", "", T("The title of the ticket. [required]"))
	cobraCmd.Flags().StringVar(&thisCmd.Body, "body", "", T("The ticket body"))
	cobraCmd.Flags().IntVar(&thisCmd.Priority, "priority", 0, T("Ticket priority [1|2|3|4], from 1 (Critical) to 4 (Minimal Impact). Only settable with Advanced and Premium support. See https://www.ibm.com/cloud/support"))
	cobraCmd.Flags().StringVar(&thisCmd.AttachmentType, "attachment-type", "", T("Specify the type of attachment, hardware or virtual. default is hardware"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CreateStandardTicketCommand) Run(args []string) error {
	if cmd.SubjectId == 0 {
		return errors.NewInvalidUsageError(T("This command requires the --subject-id option."))
	}

	if cmd.Title == "" {
		return errors.NewInvalidUsageError(T("This command requires the --title option."))
	}

	var content string
	var err error

	if cmd.Body == "" {
		content, err = cmd.TicketManager.GetText()
		if err != nil {
			return err
		}

	} else {
		content = cmd.Body
	}

	ticketArgs := managers.TicketArguments{}

	title := cmd.Title
	subjectId := cmd.SubjectId
	priority := cmd.Priority

	ticketArgs.Title = &title
	ticketArgs.Content = &content
	ticketArgs.SubjectId = &subjectId
	ticketArgs.Priority = &priority

	if cmd.Attachment != 0 {

		attachmentType := cmd.AttachmentType
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

		if cmd.RootPwd == "" {
			return cli.NewExitError(T("Root password must be provided with rootpwd flag if attachment is set."), 1)
		}
		id := cmd.Attachment
		ticketArgs.AttachmentId = &id
		pwd := cmd.RootPwd
		ticketArgs.RootPassword = &pwd
	}

	ticket_id, err := cmd.TicketManager.CreateStandardTicket(&ticketArgs)
	if err != nil {
		return cli.NewExitError(T("Ticket could not be created: {{.Error}}.", map[string]interface{}{"Error": err.Error()}), 2)
	}

	cmd.UI.Print(T("Ticket ID: {{.TicketID}}.", map[string]interface{}{"TicketID": *ticket_id}))

	return err
}
