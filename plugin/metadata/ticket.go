package metadata

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
)

var (
	NS_TICKET_NAME  = "ticket"
	CMD_TICKET_NAME = "ticket"

	CMD_TICKET_CREATE_NAME   = "create"
	CMD_TICKET_ATTACH_NAME   = "attach"
	CMD_TICKET_DETACH_NAME   = "detach"
	CMD_TICKET_DETAIL_NAME   = "detail"
	CMD_TICKET_UPDATE_NAME   = "update"
	CMD_TICKET_SUBJECTS_NAME = "subjects"
	CMD_TICKET_UPLOAD_NAME   = "upload"
	CMD_TICKET_LIST_NAME     = "list"
	CMD_TICKET_SUMMARY_NAME  = "summary"
)

func TicketNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  NS_SL_NAME,
		Name:        NS_TICKET_NAME,
		Description: T("Classic infrastructure Manage Tickets"),
	}
}

func TicketMetaData() cli.Command {
	return cli.Command{
		Category:    NS_SL_NAME,
		Name:        CMD_TICKET_NAME,
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

func TicketCreateMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_TICKET_NAME,
		Name:        CMD_TICKET_CREATE_NAME,
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

func TicketDetailMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_TICKET_NAME,
		Name:        CMD_TICKET_DETAIL_NAME,
		Description: T("Get details for a ticket"),
		Usage: T(`${COMMAND_NAME} sl ticket detail TICKETID [OPTIONS]
  
EXAMPLE:
  ${COMMAND_NAME} sl ticket detail 767676
  ${COMMAND_NAME} sl ticket detail 767676 --count 10`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "count",
				Usage: T("Number of updates"),
			},
		},
	}
}

func TicketAttachMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_TICKET_NAME,
		Name:        CMD_TICKET_ATTACH_NAME,
		Description: T("Attach devices to ticket"),
		Usage: T(`${COMMAND_NAME} sl ticket attach TICKETID [OPTIONS]
  
EXAMPLE:
  ${COMMAND_NAME} sl ticket attach 7676767 --hardware 8675654 
  ${COMMAND_NAME} sl ticket attach 7676767 --virtual 1234567 `),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "hardware",
				Usage: T("The identifier for hardware to attach"),
			},
			cli.IntFlag{
				Name:  "virtual",
				Usage: T("The identifier for a virtual server to attach"),
			},
		},
	}
}

func TicketDetachMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_TICKET_NAME,
		Name:        CMD_TICKET_DETACH_NAME,
		Description: T("Detach devices from a ticket"),
		Usage: T(`${COMMAND_NAME} sl ticket detach TICKETID [OPTIONS]
  
EXAMPLE:
  ${COMMAND_NAME} sl ticket detach 767676 --hardware 8675654
  ${COMMAND_NAME} sl ticket detach 767676 --virtual 1234567`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "hardware",
				Usage: T("The identifier for hardware to detach"),
			},
			cli.IntFlag{
				Name:  "virtual",
				Usage: T("The identifier for a virtual server to detach"),
			},
		},
	}
}

func TicketSubjectsMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_TICKET_NAME,
		Name:        CMD_TICKET_SUBJECTS_NAME,
		Description: T("List Subject IDs for ticket creation"),
		Usage: T(`${COMMAND_NAME} sl ticket subjects
  
EXAMPLE:
  ${COMMAND_NAME} sl ticket subjects
  `),
	}
}

func TicketUpdataMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_TICKET_NAME,
		Name:        CMD_TICKET_UPDATE_NAME,
		Description: T("Adds an update to an existing ticket"),
		Usage: T(`${COMMAND_NAME} sl ticket update TICKETID ["CONTENTS"] 
  
    If the second argument is not specified on a non-Windows machine, it will attempt to use either the value stored in the EDITOR environmental variable, or find either nano, vim, or emacs in that order.
  
EXAMPLE:
  ${COMMAND_NAME} sl ticket update 767676 "A problem has been detected."
  ${COMMAND_NAME} sl ticket update 767667`),
	}
}

func TicketListMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_TICKET_NAME,
		Name:        CMD_TICKET_LIST_NAME,
		Description: T("List tickets"),
		Usage:       T("${COMMAND_NAME} sl ticket list [OPTIONS]"),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "open",
				Usage: T("Display only open tickets"),
			},
			cli.BoolFlag{
				Name:  "closed",
				Usage: T("Display only closed tickets"),
			},
		},
	}
}

func TicketUploadMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_TICKET_NAME,
		Name:        CMD_TICKET_UPLOAD_NAME,
		Description: T("Adds an attachment to an existing ticket"),
		Usage: T(`${COMMAND_NAME} sl ticket upload TICKETID FILEPATH
  
EXAMPLE:
	${COMMAND_NAME} sl ticket upload 767676 "/home/user/screenshot.png"`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "name",
				Usage: T("The name of the attachment shown in the ticket"),
			},
		},
	}
}

func TicketSummaryMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_TICKET_NAME,
		Name:        CMD_TICKET_SUMMARY_NAME,
		Description: T("Summary info about tickets"),
		Usage:       "${COMMAND_NAME} sl ticket summary",
	}
}
