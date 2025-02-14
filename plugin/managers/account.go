package managers

import (
	"fmt"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

//counterfeiter:generate -o ../testhelpers/ . AccountManager
type AccountManager interface {
	GetBandwidthPools() ([]datatypes.Network_Bandwidth_Version1_Allotment, error)
	GetBandwidthPoolServers(identifier int) (uint, error)
	GetBillingItems(objectMask string, objectFilter string) ([]datatypes.Billing_Item, error)
	GetEvents(typeEvent string, mask string, dateFilter string) ([]datatypes.Notification_Occurrence_Event, error)
	GetEventDetail(identifier int, mask string) (datatypes.Notification_Occurrence_Event, error)
	AckEvent(identifier int) (bool, error)
	GetInvoiceDetail(identifier int, mask string) ([]datatypes.Billing_Invoice_Item, error)
	GetInvoices(limit int, closed bool, getAll bool) ([]datatypes.Billing_Invoice, error)
	CancelItem(identifier int) error
	GetItemDetail(identifier int, mask string) (datatypes.Billing_Item, error)
	GetActiveVirtualLicenses(mask string) ([]datatypes.Software_VirtualLicense, error)
	GetActiveAccountLicenses(mask string) ([]datatypes.Software_AccountLicense, error)
	GetAccountAllBillingOrders(mask string, limit int) ([]datatypes.Billing_Order, error)
	GetSummary(mask string) (datatypes.Account, error)
	GetBandwidthPoolDetail(bandwidthPoolId int, mask string) (datatypes.Network_Bandwidth_Version1_Allotment, error)
	GetPostProvisioningHooks(mask string) ([]datatypes.Provisioning_Hook, error)
	CreateProvisioningScript(template datatypes.Provisioning_Hook) (datatypes.Provisioning_Hook, error)
	DeleteProvisioningScript(idProvisioningScript int) (resp bool, err error)
	GetUpgradeRequests(mask string, limit int) ([]datatypes.Product_Upgrade_Request, error)
}

type accountManager struct {
	AccountService services.Account
	Session        *session.Session
}

func NewAccountManager(session *session.Session) *accountManager {
	return &accountManager{
		AccountService: services.GetAccountService(session),
		Session:        session,
	}
}

// https://sldn.softlayer.com/reference/services/SoftLayer_Account/getBandwidthAllotments/
func (a accountManager) GetBandwidthPools() ([]datatypes.Network_Bandwidth_Version1_Allotment, error) {
	mask := "mask[totalBandwidthAllocated,locationGroup, id, name, projectedPublicBandwidthUsage,billingCyclePublicBandwidthUsage[amountOut,amountIn],billingItem[id,nextInvoiceTotalRecurringAmount],outboundPublicBandwidthUsage,serviceProviderId,bandwidthAllotmentTypeId,activeDetailCount]"
	pools, err := a.AccountService.Mask(mask).GetBandwidthAllotments()
	return pools, err
}

/*
Gets a count of all servers in a bandwidth pool
Getting the server counts individually is significantly faster than pulling them in
with the GetBandwidthPools api call.
*/
func (a accountManager) GetBandwidthPoolServers(identifier int) (uint, error) {
	mask := "mask[id, bareMetalInstanceCount, hardwareCount, virtualGuestCount]"
	allotmentService := services.GetNetworkBandwidthVersion1AllotmentService(a.Session)
	counts, err := allotmentService.Mask(mask).Id(identifier).GetObject()
	var total uint
	total = 0
	if counts.BareMetalInstanceCount != nil {
		total += *counts.BareMetalInstanceCount
	}
	if counts.HardwareCount != nil {
		total += *counts.HardwareCount
	}
	if counts.VirtualGuestCount != nil {
		total += *counts.VirtualGuestCount
	}
	return total, err
}

/*
Gets All billing items of an account.
https://sldn.softlayer.com/reference/services/SoftLayer_Account/getAllTopLevelBillingItems/
*/
func (a accountManager) GetBillingItems(objectMask string, objectFilter string) ([]datatypes.Billing_Item, error) {

	i := 0
	resourceList := []datatypes.Billing_Item{}
	for {
		resp, err := a.AccountService.Mask(objectMask).Filter(objectFilter).Limit(metadata.LIMIT).Offset(i * metadata.LIMIT).GetAllTopLevelBillingItems()
		i++
		if err != nil {
			return []datatypes.Billing_Item{}, err
		}
		resourceList = append(resourceList, resp...)
		if len(resp) < metadata.LIMIT {
			break
		}
	}
	return resourceList, nil
}

/*
Gets a list of top-level invoice items that are on the currently pending invoice.
https://sldn.softlayer.com/reference/services/SoftLayer_Billing_Invoice/getInvoiceTopLevelItems/
*/
func (a accountManager) GetInvoiceDetail(identifier int, mask string) ([]datatypes.Billing_Invoice_Item, error) {
	BillingInoviceService := services.GetBillingInvoiceService(a.Session)

	filters := filter.New()
	filters = append(filters, filter.Path("invoiceTopLevelItems.id").OrderBy("DESC"))

	i := 0
	resourceList := []datatypes.Billing_Invoice_Item{}
	for {
		resp, err := BillingInoviceService.Mask(mask).Filter(filters.Build()).Limit(metadata.LIMIT).Offset(i * metadata.LIMIT).Id(identifier).GetInvoiceTopLevelItems()
		i++
		if err != nil {
			return []datatypes.Billing_Invoice_Item{}, err
		}
		resourceList = append(resourceList, resp...)
		if len(resp) < metadata.LIMIT {
			break
		}
	}
	return resourceList, nil
}

/*
Gets all events with the potential to cause a service interruption with a specific keyName.
https://sldn.softlayer.com/reference/services/SoftLayer_Notification_Occurrence_Event/getAllObjects/
*/
func (a accountManager) GetEvents(typeEvent string, mask string, dateFilter string) ([]datatypes.Notification_Occurrence_Event, error) {
	NotificationOccurrenceEventService := services.GetNotificationOccurrenceEventService(a.Session)
	filters := filter.New()
	filters = append(filters, filter.Path("id").OrderBy("ASC"))
	filters = append(filters, filter.Path("notificationOccurrenceEventType.keyName").Eq(typeEvent))
	if dateFilter != "" {
		if typeEvent == "PLANNED" {
			filters = append(filters, filter.Path("endDate").DateAfter(dateFilter))
		}
		if typeEvent == "UNPLANNED_INCIDENT" {
			filters = append(filters, filter.Path("modifyDate").DateAfter(dateFilter))
		}
	}
	if typeEvent == "ANNOUNCEMENT" {
		filters = append(filters, filter.Path("statusCode.keyName").Eq("PUBLISHED"))
	}

	resourceList, err := NotificationOccurrenceEventService.Mask(mask).Filter(filters.Build()).GetAllObjects()
	if err != nil {
		return []datatypes.Notification_Occurrence_Event{}, err
	}
	return resourceList, err
}

/*
Gets a event with the potential to cause a service interruption.
https://sldn.softlayer.com/reference/services/SoftLayer_Notification_Occurrence_Event/getObject/
*/
func (a accountManager) GetEventDetail(identifier int, mask string) (datatypes.Notification_Occurrence_Event, error) {
	NotificationOccurrenceEventService := services.GetNotificationOccurrenceEventService(a.Session)

	resourceList, err := NotificationOccurrenceEventService.Mask(mask).Id(identifier).GetObject()
	if err != nil {
		return datatypes.Notification_Occurrence_Event{}, err
	}
	return resourceList, err
}

/*
Acknowledge Event. Doing so will turn off the popup in the control portal
https://sldn.softlayer.com/reference/services/SoftLayer_Notification_Occurrence_Event/acknowledgeNotification/
*/
func (a accountManager) AckEvent(identifier int) (bool, error) {
	NotificationOccurrenceEventService := services.GetNotificationOccurrenceEventService(a.Session)

	return NotificationOccurrenceEventService.Id(identifier).AcknowledgeNotification()
}

/*
Gets all invoices from the account
https://sldn.softlayer.com/reference/services/SoftLayer_Account/getInvoices/
*/
func (a accountManager) GetInvoices(limit int, closed bool, getAll bool) ([]datatypes.Billing_Invoice, error) {
	mask := "mask[invoiceTotalAmount, itemCount]"
	filters := filter.New()
	filters = append(filters, filter.Path("invoices.id").OrderBy("DESC"))
	if !closed || !getAll {
		filters = append(filters, filter.Path("invoices.statusCode").Eq("OPEN"))
	}
	resourceList := []datatypes.Billing_Invoice{}
	if getAll {
		i := 0
		for {
			resp, err := a.AccountService.Mask(mask).Filter(filters.Build()).Limit(metadata.LIMIT).Offset(i * metadata.LIMIT).GetInvoices()
			i++
			if err != nil {
				return []datatypes.Billing_Invoice{}, err
			}
			resourceList = append(resourceList, resp...)
			if len(resp) < metadata.LIMIT {
				break
			}
		}
	} else {
		resp, err := a.AccountService.Mask(mask).Filter(filters.Build()).Limit(limit).GetInvoices()
		if err != nil {
			return []datatypes.Billing_Invoice{}, err
		}
		resourceList = append(resourceList, resp...)
	}
	return resourceList, nil
}

/*
Cancels the resource or service for a billing Item
https://sldn.softlayer.com/reference/services/SoftLayer_Billing_Item/cancelItem/
*/
func (a accountManager) CancelItem(identifier int) error {
	BillingItemService := services.GetBillingItemService(a.Session)

	CancelImmediately := false
	cancelAssociatedBillingItems := true
	Reason := "No longer needed"

	mask := "mask[id,displayName,email,username]"
	user, _ := a.AccountService.Mask(mask).GetCurrentUser()
	Note := fmt.Sprintf("Cancelled by %s with the ibmcloud sl", utils.FormatStringPointerName(user.Username))

	_, err := BillingItemService.Mask(mask).Id(identifier).CancelItem(&CancelImmediately, &cancelAssociatedBillingItems, &Reason, &Note)
	return err
}

/*
Gets the detail of a item
https://sldn.softlayer.com/reference/services/SoftLayer_Billing_Item/getObject/
*/
func (a accountManager) GetItemDetail(identifier int, mask string) (datatypes.Billing_Item, error) {
	BillingItemService := services.GetBillingItemService(a.Session)
	return BillingItemService.Mask(mask).Id(identifier).GetObject()
}

/*
Gets virtual software licenses controlled by an account
https://sldn.softlayer.com/reference/services/SoftLayer_Account/getActiveVirtualLicenses/
*/
func (a accountManager) GetActiveVirtualLicenses(mask string) ([]datatypes.Software_VirtualLicense, error) {
	filters := filter.New()
	filters = append(filters, filter.Path("activeVirtualLicenses.id").OrderBy("ASC"))
	i := 0
	resourceList := []datatypes.Software_VirtualLicense{}
	for {
		resp, err := a.AccountService.Mask(mask).Filter(filters.Build()).Limit(metadata.LIMIT).Offset(i * metadata.LIMIT).GetActiveVirtualLicenses()
		i++
		if err != nil {
			return []datatypes.Software_VirtualLicense{}, err
		}
		resourceList = append(resourceList, resp...)
		if len(resp) < metadata.LIMIT {
			break
		}
	}
	return resourceList, nil
}

/*
Gets active account software licenses owned by an account
https://sldn.softlayer.com/reference/services/SoftLayer_Account/getActiveVirtualLicenses/
*/
func (a accountManager) GetActiveAccountLicenses(mask string) ([]datatypes.Software_AccountLicense, error) {
	i := 0
	resourceList := []datatypes.Software_AccountLicense{}
	for {
		resp, err := a.AccountService.Mask(mask).Limit(metadata.LIMIT).Offset(i * metadata.LIMIT).GetActiveAccountLicenses()
		i++
		if err != nil {
			return []datatypes.Software_AccountLicense{}, err
		}
		resourceList = append(resourceList, resp...)
		if len(resp) < metadata.LIMIT {
			break
		}
	}
	return resourceList, nil
}

/*
Gets all billing orders for your account
https://sldn.softlayer.com/reference/services/SoftLayer_Billing_Order/getAllObjects/
*/
func (a accountManager) GetAccountAllBillingOrders(mask string, limit int) ([]datatypes.Billing_Order, error) {
	BillingOrderService := services.GetBillingOrderService(a.Session)
	filters := filter.New()
	filters = append(filters, filter.Path("id").OrderBy("DESC"))

	return BillingOrderService.Mask(mask).Filter(filters.Build()).Limit(limit).GetAllObjects()
}

/*
Gets a SoftLayer_Account record.
https://sldn.softlayer.com/reference/services/SoftLayer_Account/getObject/
*/
func (a accountManager) GetSummary(mask string) (datatypes.Account, error) {
	return a.AccountService.Mask(mask).GetObject()
}

/*
Gets a SoftLayer_Network_Bandwidth_Version1_Allotment object record.
https://sldn.softlayer.com/reference/services/SoftLayer_Network_Bandwidth_Version1_Allotment/getObject/
*/
func (a accountManager) GetBandwidthPoolDetail(bandwidthPoolId int, mask string) (datatypes.Network_Bandwidth_Version1_Allotment, error) {
	networkBandwidthService := services.GetNetworkBandwidthVersion1AllotmentService(a.Session)
	if mask == "" {
		mask = `mask[activeDetails[allocation],projectedPublicBandwidthUsage, billingCyclePublicBandwidthUsage,
        hardware[outboundBandwidthUsage,bandwidthAllotmentDetail[allocation]],inboundPublicBandwidthUsage,
        virtualGuests[outboundPublicBandwidthUsage,bandwidthAllotmentDetail[allocation]],
        bareMetalInstances[outboundBandwidthUsage,bandwidthAllotmentDetail[allocation]]]`
	}
	return networkBandwidthService.Id(bandwidthPoolId).Mask(mask).GetObject()
}

/*
Customer specified URIs that are downloaded onto a newly provisioned or reloaded server.
If the URI is sent over https it will be executed directly on the server.
https://sldn.softlayer.com/reference/services/SoftLayer_Account/getPostProvisioningHooks/
*/
func (a accountManager) GetPostProvisioningHooks(mask string) ([]datatypes.Provisioning_Hook, error) {
	if mask == "" {
		mask = "mask[id,name,uri]"
	}
	return a.AccountService.Mask(mask).GetPostProvisioningHooks()
}

/*
Create a provisioning script.
https://sldn.softlayer.com/reference/services/SoftLayer_Provisioning_Hook/createObject/
*/
func (a accountManager) CreateProvisioningScript(template datatypes.Provisioning_Hook) (datatypes.Provisioning_Hook, error) {
	provisioningHook := services.GetProvisioningHookService(a.Session)
	return provisioningHook.CreateObject(&template)
}

/*
Delete a provisioning script.
https://sldn.softlayer.com/reference/services/SoftLayer_Provisioning_Hook/deleteObject/
*/
func (a accountManager) DeleteProvisioningScript(idProvisioningScript int) (resp bool, err error) {
	provisioningHook := services.GetProvisioningHookService(a.Session)
	return provisioningHook.Id(idProvisioningScript).DeleteObject()
}

/*
Gets account's associated upgrade requests.
https://sldn.softlayer.com/reference/services/SoftLayer_Account/getUpgradeRequests/
*/
func (a accountManager) GetUpgradeRequests(mask string, limit int) ([]datatypes.Product_Upgrade_Request, error) {
	return a.AccountService.Limit(limit).Mask(mask).GetUpgradeRequests()
}
