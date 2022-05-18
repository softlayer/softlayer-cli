package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type EmailManager interface {
	GetNetworkMessageDeliveryAccounts(mask string) ([]datatypes.Network_Message_Delivery, error)
	GetAccountOverview(emailId int) (datatypes.Container_Network_Message_Delivery_Email_Sendgrid_Account_Overview, error)
	GetStatistics(emailId int) ([]datatypes.Container_Network_Message_Delivery_Email_Sendgrid_Statistics, error)
	GetInstance(emailId int, mask string) (datatypes.Network_Message_Delivery_Email_Sendgrid, error)
	UpdateEmail(emailId int, emailAddress string) (error)
	EditObject(emailId int, username string, password string) (error)
}

type emailManager struct {
	EmailService services.Network_Message_Delivery_Email_Sendgrid
	Session      *session.Session
}

func NewEmailManager(session *session.Session) *emailManager {
	return &emailManager{
		EmailService: services.GetNetworkMessageDeliveryEmailSendgridService(session),
		Session:      session,
	}
}

/*
Gets all emails by account.
https://sldn.softlayer.com/reference/services/SoftLayer_Account/getNetworkMessageDeliveryAccounts/
*/
func (a emailManager) GetNetworkMessageDeliveryAccounts(mask string) ([]datatypes.Network_Message_Delivery, error) {
	NetworkMessageDeliveryService := services.GetAccountService(a.Session)

	filters := filter.New()
	filters = append(filters, filter.Path("id").OrderBy("DESC"))

	i := 0
	resourceList := []datatypes.Network_Message_Delivery{}
	for {
		resp, err := NetworkMessageDeliveryService.Mask(mask).Filter(filters.Build()).Limit(metadata.LIMIT).Offset(i * metadata.LIMIT).GetNetworkMessageDeliveryAccounts()
		i++
		if err != nil {
			return []datatypes.Network_Message_Delivery{}, err
		}
		resourceList = append(resourceList, resp...)
		if len(resp) < metadata.LIMIT {
			break
		}
	}
	return resourceList, nil
}

/*
Gets account overview by email.
https://sldn.softlayer.com/reference/services/SoftLayer_Network_Message_Delivery_Email_Sendgrid/getAccountOverview/
*/
func (a emailManager) GetAccountOverview(emailId int) (datatypes.Container_Network_Message_Delivery_Email_Sendgrid_Account_Overview, error) {
	return a.EmailService.Id(emailId).GetAccountOverview()
}

/*
Gets all statistics by email.
https://sldn.softlayer.com/reference/services/SoftLayer_Network_Message_Delivery_Email_Sendgrid/getStatistics/
*/
func (a emailManager) GetStatistics(emailId int) ([]datatypes.Container_Network_Message_Delivery_Email_Sendgrid_Statistics, error) {
	options := datatypes.Container_Network_Message_Delivery_Email_Sendgrid_Statistics_Options{}
	return a.EmailService.Id(emailId).GetStatistics(&options)
}

/*
Gets a SoftLayer_Network_Message_Delivery_Email_Sendgrid record.
https://sldn.softlayer.com/reference/services/SoftLayer_Network_Message_Delivery_Email_Sendgrid/getObject/
*/
func (a emailManager) GetInstance(emailId int, mask string) (datatypes.Network_Message_Delivery_Email_Sendgrid, error) {
	return a.EmailService.Mask(mask).Id(emailId).GetObject()
}

/*
Edits a email adrress from a user.
https://sldn.softlayer.com/reference/services/SoftLayer_Network_Message_Delivery_Email_Sendgrid/updateEmailAddress/
*/
func (a emailManager) UpdateEmail(emailId int, emailAddress string) (error) {
	_, err :=a.EmailService.Id(emailId).UpdateEmailAddress(&emailAddress)
	return err
}

/*
Edits the email object from a user.
https://sldn.softlayer.com/reference/services/SoftLayer_Network_Message_Delivery_Email_Sendgrid/editObject/
*/
func (a emailManager) EditObject(emailId int, username string, password string) (error) {

	templateObject := datatypes.Network_Message_Delivery{}
	if username != "" {
		templateObject.Username = sl.String(username)
	}
	if password != "" {
		templateObject.Password = sl.String(password)
	}
	_, err :=a.EmailService.Id(emailId).EditObject(&templateObject)
	return err
}
