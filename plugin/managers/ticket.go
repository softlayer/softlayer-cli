package managers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

const mask = "mask[id, title, assignedUser[firstName, lastName], priority, createDate, lastEditDate, accountId, status, updateCount]"

type TicketManager interface {
	CreateStandardTicket(ticketArgs *TicketArguments) (*int, error)
	AttachDeviceToTicket(ticketId int, deviceId int, ishardware bool) error
	RemoveDeviceFromTicket(ticketId int, deviceId int, ishardware bool) error
	GetAllUpdates(ticketId int) (updates []datatypes.Ticket_Update, err error)
	AddUpdate(ticketId int, content string) error
	GetSubjects() (*[]datatypes.Ticket_Subject, error)
	ListTickets() ([]datatypes.Ticket, error)
	ListOpenTickets() ([]datatypes.Ticket, error)
	ListCloseTickets() ([]datatypes.Ticket, error)
	AttachFileToTicket(ticketId int, name string, path string) error
	Summary() (*TicketSummary, error)
	GetText() (string, error)
	GetTicket(ticketid int) (datatypes.Ticket, error)
}

type ticketManager struct {
	AccountService services.Account
	TicketService  services.Ticket
	TicketSubject  services.Ticket_Subject
}

func NewTicketManager(session *session.Session) *ticketManager {
	return &ticketManager{
		services.GetAccountService(session),
		services.GetTicketService(session),
		services.GetTicketSubjectService(session),
	}
}

type TicketArguments struct {
	AttachmentId   *int
	RootPassword   *string
	Content        *string
	Title          *string
	SubjectId      *int
	Priority       *int
	AttachmentType *string
}

type TicketSummary struct {
	Accounting uint
	Billing    uint
	Sales      uint
	Support    uint
	Other      uint
	Closed     uint
	Open       uint
}

var ( // Addressable constants.
	TRUE        = true // Might turn into a --dry flag depending on the response from bizdev.
	APIQuestion = 1522
)

const five_MB = 1024 * 1000 * 5

func (ticket ticketManager) CreateStandardTicket(ticketArgs *TicketArguments) (*int, error) {
	if ticketArgs.SubjectId == nil {
		ticketArgs.SubjectId = &APIQuestion
	}

	templateObject := datatypes.Ticket{Title: ticketArgs.Title, SubjectId: ticketArgs.SubjectId, Priority: ticketArgs.Priority}
	files := make([]datatypes.Container_Utility_File_Attachment, 0, 1) // Attaching files is not implemented for now.

	resp, err := ticket.TicketService.CreateStandardTicket(&templateObject, ticketArgs.Content, ticketArgs.AttachmentId, ticketArgs.RootPassword, nil, nil, files, ticketArgs.AttachmentType)
	if err != nil {
		return nil, err
	}
	return resp.Id, nil

}

func (ticket ticketManager) AttachDeviceToTicket(ticketId int, deviceId int, ishardware bool) error {
	t := ticket.TicketService.Id(ticketId)
	if !ishardware { // Virtual device
		vsi, err := t.AddAttachedVirtualGuest(&deviceId, &TRUE)
		if err != nil {
			return err
		}

		if vsi.AttachmentId != nil {
			return nil
		} else {
			return errors.New(T("Attachment failed. Confirm that {{.DeviceID}} is a virtual guest or not already attached.", map[string]interface{}{"DeviceID": deviceId}))
		}
	} else { // Hardware device
		hwd, err := t.AddAttachedHardware(&deviceId)
		if err != nil {
			return err
		}

		if hwd.AttachmentId != nil {
			return nil
		} else {
			return errors.New(T("Attachment failed. Confirm that {{.DeviceID}} is a hardware device or not already attached.", map[string]interface{}{"DeviceID": deviceId}))
		}
	}
}

func (ticket ticketManager) GetTicket(ticketid int) (datatypes.Ticket, error) {
	mask := "mask[id, title, assignedUser[firstName, lastName],status,createDate,lastEditDate,updates[entry,editor],updateCount, priority]"
	return ticket.TicketService.Id(ticketid).Mask(mask).GetObject()
}

func (ticket ticketManager) RemoveDeviceFromTicket(ticketId int, deviceId int, ishardware bool) error {
	t := ticket.TicketService.Id(ticketId)

	if !ishardware { // Virtual device
		vsi, err := t.RemoveAttachedVirtualGuest(&deviceId)
		if err != nil {
			return err
		}

		if !vsi {
			return errors.New(T("Could not remove device {{.DeviceID}} from ticket. Is it attached, or a virtual guest?", map[string]interface{}{"DeviceID": deviceId}))
		}

	} else { // Hardware device
		hwd, err := t.RemoveAttachedHardware(&deviceId)
		if err != nil {
			return err
		}

		if !hwd {
			return errors.New(T("Could not remove device {{.DeviceId}} from ticket. Is it attached, or a hardware device?", map[string]interface{}{"DeviceID": deviceId}))
		}

	}
	return nil

}

func (ticket ticketManager) GetAllUpdates(ticketId int) (updates []datatypes.Ticket_Update, err error) {
	t := ticket.TicketService.Id(ticketId)
	_updates, err := t.GetUpdates()
	return _updates, err
}

func (ticket ticketManager) AddUpdate(ticketId int, content string) error {
	t := ticket.TicketService.Id(ticketId)

	var update datatypes.Ticket_Update

	update.TicketId = &ticketId
	update.Entry = &content

	files := make([]datatypes.Container_Utility_File_Attachment, 0, 1)

	_, err := t.AddUpdate(&update, files)
	return err

}

func (ticket ticketManager) GetSubjects() (*[]datatypes.Ticket_Subject, error) {
	subs, err := ticket.TicketSubject.GetAllObjects()
	if err != nil {
		return nil, err
	} else {
		return &subs, err
	}
}

func (ticket ticketManager) ListTickets() ([]datatypes.Ticket, error) {
	tickets, err := ticket.AccountService.Mask(mask).GetTickets()
	return tickets, err
}

func (ticket ticketManager) ListOpenTickets() ([]datatypes.Ticket, error) {
	tickets, err := ticket.AccountService.Mask(mask).GetOpenTickets()
	return tickets, err
}
func (ticket ticketManager) ListCloseTickets() ([]datatypes.Ticket, error) {
	tickets, err := ticket.AccountService.Mask(mask).GetClosedTickets()
	return tickets, err
}

func (ticket ticketManager) AttachFileToTicket(ticketId int, name string, path string) error {
	file_stat, err := os.Stat(path)
	if err != nil {
		return err
	}

	file, err := os.Open(path) // #nosec
	if err != nil {
		return err
	}

	buffer := make([]byte, file_stat.Size())
	_, err = file.Read(buffer)

	if name == "" {
		name = file_stat.Name()
	}

	attachment := datatypes.Container_Utility_File_Attachment{Data: &buffer, Filename: &name}

	t := ticket.TicketService.Id(ticketId)
	_, err = t.AddAttachedFile(&attachment)

	return err
}

func (ticket ticketManager) GetText() (string, error) {
	_editor := os.Getenv("EDITOR")
	if os.Getenv("OS") != "Windows_NT" {
		if _editor == "" {
			if _editor = getpath("nano"); _editor != "" {

			} else if _editor = getpath("vim"); _editor != "" {

			} else if _editor = getpath("emacs"); _editor != "" {

			} else {
				return "", errors.New(T("Editor could not be found. Please set EDITOR environmental variable or specifiy a message argument."))
			}
		} else if isfullpath := strings.Split(_editor, "/"); len(isfullpath) == 1 {
			_editor = getpath(_editor)
			if _editor == "" {
				return "", errors.New(T("Editor could not be found. Please set EDITOR environmental variable or specifiy a message argument."))
			}
		}
	} else {
		_editor = "notepad"
	}
	tmp_file, err := ioutil.TempFile("", "ibmcloud-sl-updatebuf")
	defer os.Remove(tmp_file.Name())
	if err != nil {
		return "", errors.New(T("Buffer file could not be created: {{.Error}}.\n", map[string]interface{}{"Error": err.Error()}))
	} else {
		fileCloseErr := tmp_file.Close()
		if fileCloseErr != nil {
			return "", errors.New(T("Buffer file could not be closed: {{.Error}}.\n", map[string]interface{}{"Error": fileCloseErr.Error()}))
		}
	}

	editor := exec.Command(_editor, tmp_file.Name()) // #nosec
	editor.Stdin = os.Stdin
	editor.Stdout = os.Stdout

	if err := editor.Run(); err != nil {
		return "", errors.New(T("Editor could not be ran: {{.Error}}.\n", map[string]interface{}{"Error": err.Error()}))
	}

	output := make([]byte, 1024, 4000)
	_file, _ := os.Open(tmp_file.Name())
	count, err := _file.Read(output)
	if err != nil {
		return "", errors.New(T("Read Failure: {{.Error}}.\n", map[string]interface{}{"Error": err.Error()}))
	}

	content := fmt.Sprintf("%s\n\n***POSTED FROM IBMCLOUD SL***", output[:count]) // #nosec
	return content, nil
}

func (ticket ticketManager) Summary() (*TicketSummary, error) {
	mask := "mask[openTicketCount,closedTicketCount,openBillingTicketCount,openOtherTicketCount,openSalesTicketCount,openSupportTicketCount,openAccountingTicketCount]"
	resp, err := ticket.AccountService.Mask(mask).GetObject()
	if err != nil {
		return nil, err
	}

	summary := TicketSummary{utils.UIntPointertoUInt(resp.OpenAccountingTicketCount), utils.UIntPointertoUInt(resp.OpenBillingTicketCount), utils.UIntPointertoUInt(resp.OpenSalesTicketCount), utils.UIntPointertoUInt(resp.OpenSupportTicketCount), utils.UIntPointertoUInt(resp.OpenOtherTicketCount), utils.UIntPointertoUInt(resp.ClosedTicketCount), utils.UIntPointertoUInt(resp.OpenTicketCount)}
	return &summary, nil
}

func getpath(file string) string {
	path := strings.Split(os.Getenv("PATH"), ":")
	for _, p := range path {
		full_path := filepath.Join(p, file)
		_, err := os.Stat(full_path)
		if err == nil {
			return full_path
		}
	}

	return ""
}
