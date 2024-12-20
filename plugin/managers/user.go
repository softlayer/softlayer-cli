package managers

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"time"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

const (
	LIST_USER_MASK = "mask[id, username, displayName, userStatus[name], hardwareCount, virtualGuestCount, email, roles, externalBindingCount, apiAuthenticationKeyCount]"
	GET_USER_MASK  = "mask[id, firstName, lastName, email, companyName, address1, city, country, postalCode, state, userStatusId, timezoneId]"
)

// Manages SoftLayer Block and File Storage volumes.
// See product information here: https://www.ibm.com/cloud-computing/bluemix/block-storage, https://www.ibm.com/cloud-computing/bluemix/file-storage
//counterfeiter:generate -o ../testhelpers/ . UserManager
type UserManager interface {
	ListUsers(mask string) ([]datatypes.User_Customer, error)
	GetUser(userId int, mask string) (datatypes.User_Customer, error)
	GetCurrentUser() (datatypes.User_Customer, error)
	GetAllPermission() ([]datatypes.User_Customer_CustomerPermission_Permission, error)
	GetAllPermissionDepartments() ([]datatypes.User_Permission_Department, error)
	AddPermission(userId int, permissions []datatypes.User_Customer_CustomerPermission_Permission) (bool, error)
	RemovePermission(userId int, permissions []datatypes.User_Customer_CustomerPermission_Permission) (bool, error)
	PermissionFromUser(userId, fromUserId int) error
	GetUserPermissions(userId int) ([]datatypes.User_Customer_CustomerPermission_Permission, error)
	GetLogins(userId int, startDate time.Time) ([]datatypes.User_Customer_Access_Authentication, error)
	GetEvents(userId int, startDate time.Time) ([]datatypes.Event_Log, error)
	GetIdFromUsername(userName string) (int, error)
	FormatPermissionObject(permissionsKeyNames []string) ([]datatypes.User_Customer_CustomerPermission_Permission, error)
	CreateUser(templateObject datatypes.User_Customer, password string, vpnPassword string) (datatypes.User_Customer, error)
	EditUser(templateObject datatypes.User_Customer, UserId int) (bool, error)
	AddApiAuthenticationKey(UserId int) (string, error)
	GetApiAuthenticationKeys(userId int) ([]datatypes.User_Customer_ApiAuthentication, error)
	RemoveApiAuthenticationKey(userId int) (bool, error)
	GetAllNotifications(mask string) ([]datatypes.Email_Subscription, error)
	EnableEmailSubscriptionNotification(notificationId int) (bool, error)
	DisableEmailSubscriptionNotification(notificationId int) (bool, error)
	AddHardwareAccess(userId int, hardwareId int) (bool, error)
	AddDedicatedHostAccess(userId int, dedicatedHostId int) (bool, error)
	AddVirtualGuestAccess(userId int, virtualGuestId int) (bool, error)
	RemoveHardwareAccess(userId int, hardwareId int) (bool, error)
	RemoveDedicatedHostAccess(userId int, dedicatedHostId int) (bool, error)
	RemoveVirtualGuestAccess(userId int, virtualGuestId int) (bool, error)
	GetUserAllowDevicesPermissions(userId int) ([]datatypes.User_Customer_CustomerPermission_Permission, error)
	GetDedicatedHosts(userId int, mask string) ([]datatypes.Virtual_DedicatedHost, error)
	GetHardware(userId int, mask string) ([]datatypes.Hardware, error)
	GetVirtualGuests(userId int, mask string) ([]datatypes.Virtual_Guest, error)
	CreateUserVpnOverride(userId int, subnetId int) (bool, error)
	UpdateVpnUser(userId int) (bool, error)
	GetOverrides(userId int) ([]datatypes.Network_Service_Vpn_Overrides, error)
	DeleteUserVpnOverride(overrideId int) (bool, error)
	UpdateVpnPassword(userID int, password string) (bool, error)
}

type userManager struct {
	AccountService        services.Account
	UserCustomerService   services.User_Customer
	UserPermissionService services.User_Customer_CustomerPermission_Permission
	EventLogService       services.Event_Log
	Email_Subscription    services.Email_Subscription
	UserPermissionAction  services.User_Permission_Action
	Session               *session.Session
}

func NewUserManager(session *session.Session) *userManager {
	return &userManager{
		services.GetAccountService(session),
		services.GetUserCustomerService(session),
		services.GetUserCustomerCustomerPermissionPermissionService(session),
		services.GetEventLogService(session),
		services.GetEmailSubscriptionService(session),
		services.GetUserPermissionActionService(session),
		session,
	}
}

func (u userManager) ListUsers(mask string) ([]datatypes.User_Customer, error) {
	if mask == "" {
		mask = LIST_USER_MASK
	}
	return u.AccountService.Mask(mask).GetUsers()
}

func (u userManager) GetUser(userId int, mask string) (datatypes.User_Customer, error) {
	if mask == "" {
		return u.UserCustomerService.Id(userId).Mask(GET_USER_MASK).GetObject()
	}
	return u.UserCustomerService.Id(userId).Mask(mask).GetObject()
}

func (u userManager) GetCurrentUser() (datatypes.User_Customer, error) {
	return u.AccountService.Mask(GET_USER_MASK).GetCurrentUser()
}

func (u userManager) GetAllPermission() ([]datatypes.User_Customer_CustomerPermission_Permission, error) {
	permissions, err := u.UserPermissionAction.GetAllObjects()
	if err != nil {
		return nil, err
	}

	// parsing permission from datatypes.User_Permission_Action to datatypes.User_Customer_CustomerPermission_Permission
	parsedPermission := []datatypes.User_Customer_CustomerPermission_Permission{}
	for _, permission := range permissions {
		parsedPermission = append(parsedPermission, datatypes.User_Customer_CustomerPermission_Permission{
			Key:     sl.String(utils.FormatStringPointer(permission.Key)),
			KeyName: sl.String(utils.FormatStringPointer(permission.KeyName)),
			Name:    sl.String(utils.FormatStringPointer(permission.Name)),
		})
	}

	sort.Sort(utils.PermissionsBykeyName(parsedPermission))
	return parsedPermission, nil
}

func (u userManager) GetAllPermissionDepartments() ([]datatypes.User_Permission_Department, error) {
	permissionService := services.GetUserPermissionDepartmentService(u.Session)
	mask := "mask[permissions[id,description,name,keyName]]"
	permissions, err := permissionService.Mask(mask).GetAllObjects()
	if err != nil {
		return nil, err
	}

	return permissions, nil
}


func (u userManager) AddPermission(userId int, permissions []datatypes.User_Customer_CustomerPermission_Permission) (bool, error) {
	return u.UserCustomerService.Id(userId).AddBulkPortalPermission(permissions)
}

func (u userManager) RemovePermission(userId int, permissions []datatypes.User_Customer_CustomerPermission_Permission) (bool, error) {
	return u.UserCustomerService.Id(userId).RemoveBulkPortalPermission(permissions, nil)
}

func (u userManager) PermissionFromUser(userId, fromUserId int) error {
	fromPermission, err := u.GetUserPermissions(fromUserId)

	if err != nil {
		return err
	}

	_, err = u.AddPermission(userId, fromPermission)
	if err != nil {
		return err
	}

	allPermissions, err := u.GetAllPermission()
	if err != nil {
		return err
	}

	removePermission := []datatypes.User_Customer_CustomerPermission_Permission{}

	for _, permission := range allPermissions {
		if keyNameSearch(fromPermission, *permission.KeyName) {
			continue
		} else {
			removePermission = append(removePermission, permission)
		}
	}
	_, err = u.RemovePermission(userId, removePermission)

	if err != nil {
		return err
	}
	return nil
}

func (u userManager) GetUserPermissions(userId int) ([]datatypes.User_Customer_CustomerPermission_Permission, error) {
	var permissions []datatypes.User_Customer_CustomerPermission_Permission
	var err error
	isMasterUser, err := u.UserCustomerService.Id(userId).IsMasterUser()
	if isMasterUser {
		permissions, err = u.GetAllPermission()
	} else {
		permissions, err = u.UserCustomerService.Id(userId).GetPermissions()
	}

	if err != nil {
		return nil, err
	}
	sort.Sort(utils.PermissionsBykeyName(permissions))

	return permissions, err
}

func (u userManager) GetLogins(userId int, startDate time.Time) ([]datatypes.User_Customer_Access_Authentication, error) {
	if startDate.IsZero() {
		month, _ := time.ParseDuration("-24h")
		startDate = time.Now().Add(month)
	}
	filters := filter.New(filter.Path("loginAttempts.createDate").DateAfter(startDate.Format("01/02/2006 15:04:05")))
	return u.UserCustomerService.Filter(filters.Build()).Id(userId).GetLoginAttempts()
}

func (u userManager) GetEvents(userId int, startDate time.Time) ([]datatypes.Event_Log, error) {
	if startDate.IsZero() {
		month, _ := time.ParseDuration("-720h")
		startDate = time.Now().Add(month)
	}
	filters := filter.New(filter.Path("eventCreateDate").DateAfter(startDate.Format("2006-01-02T15:04:05")))

	filterUserId := filter.Path("userId")
	filterUserId.Val = strconv.Itoa(userId)
	filters = append(filters, filterUserId)

	return u.EventLogService.Filter(filters.Build()).GetAllObjects()
}

func (u userManager) GetIdFromUsername(userName string) (int, error) {

	mask := "mask[id, username]"
	filters := filter.New()
	filters = append(filters, utils.QueryFilter(userName, "users.username"))

	users, err := u.AccountService.Mask(mask).Filter(filters.Build()).GetUsers()
	if err != nil {
		return 0, err
	}

	if len(users) == 1 {
		return *users[0].Id, nil
	} else if len(users) > 1 {
		return 0, errors.New(T("Multiple users found with the name: %s", userName))
	} else {
		return 0, errors.New(T("Unable to find user id for %s", userName))

	}
}

func (u userManager) FormatPermissionObject(permissionsKeyNames []string) ([]datatypes.User_Customer_CustomerPermission_Permission, error) {
	var pretty_permissions []datatypes.User_Customer_CustomerPermission_Permission
	availablePermissions, err := u.GetAllPermission()
	if err != nil {
		return nil, err
	}

	for _, permissionsKeyName := range permissionsKeyNames {
		permissionsKeyName = strings.ToUpper(permissionsKeyName)
		if permissionsKeyName == "ALL" {
			return availablePermissions, nil
		}
		if keyNameSearch(availablePermissions, permissionsKeyName) {
			tmp := permissionsKeyName
			pretty_permissions = append(pretty_permissions, datatypes.User_Customer_CustomerPermission_Permission{KeyName: &tmp})
		} else {
			return nil, errors.New(fmt.Sprintf(T("%s is not a valid permission"), permissionsKeyName))
		}
	}
	return pretty_permissions, nil

}

func (u userManager) CreateUser(templateObject datatypes.User_Customer, password string, vpnPassword string) (datatypes.User_Customer, error) {
	return u.UserCustomerService.CreateObject(&templateObject, &password, &vpnPassword)
}

func (u userManager) EditUser(templateObject datatypes.User_Customer, UserId int) (bool, error) {
	return u.UserCustomerService.Id(UserId).EditObject(&templateObject)
}

func (u userManager) AddApiAuthenticationKey(UserId int) (string, error) {
	return u.UserCustomerService.Id(UserId).AddApiAuthenticationKey()
}

func keyNameSearch(permissions []datatypes.User_Customer_CustomerPermission_Permission, keyname string) bool {
	for _, permission := range permissions {
		if permission.KeyName != nil && *permission.KeyName == keyname {
			return true
		}
	}
	return false
}

func (u userManager) GetAllNotifications(mask string) ([]datatypes.Email_Subscription, error) {
	return u.Email_Subscription.Mask(mask).GetAllObjects()
}

func (u userManager) EnableEmailSubscriptionNotification(notificationId int) (bool, error) {
	return u.Email_Subscription.Id(notificationId).Enable()
}

func (u userManager) DisableEmailSubscriptionNotification(notificationId int) (bool, error) {
	return u.Email_Subscription.Id(notificationId).Disable()
}

func (u userManager) AddHardwareAccess(userId int, hardwareId int) (bool, error) {
	return u.UserCustomerService.Id(userId).AddHardwareAccess(&hardwareId)
}

func (u userManager) AddDedicatedHostAccess(userId int, dedicatedHostId int) (bool, error) {
	return u.UserCustomerService.Id(userId).AddDedicatedHostAccess(&dedicatedHostId)
}

func (u userManager) AddVirtualGuestAccess(userId int, virtualGuestId int) (bool, error) {
	return u.UserCustomerService.Id(userId).AddVirtualGuestAccess(&virtualGuestId)
}

func (u userManager) RemoveHardwareAccess(userId int, hardwareId int) (bool, error) {
	return u.UserCustomerService.Id(userId).RemoveHardwareAccess(&hardwareId)
}

func (u userManager) RemoveDedicatedHostAccess(userId int, dedicatedHostId int) (bool, error) {
	return u.UserCustomerService.Id(userId).RemoveDedicatedHostAccess(&dedicatedHostId)
}

func (u userManager) RemoveVirtualGuestAccess(userId int, virtualGuestId int) (bool, error) {
	return u.UserCustomerService.Id(userId).RemoveVirtualGuestAccess(&virtualGuestId)
}

func (u userManager) GetUserAllowDevicesPermissions(userId int) ([]datatypes.User_Customer_CustomerPermission_Permission, error) {
	filters := filter.New(filter.Path("permissions.key").Contains("All_"))
	return u.UserCustomerService.Id(userId).Filter(filters.Build()).GetPermissions()
}

func (u userManager) GetDedicatedHosts(userId int, mask string) ([]datatypes.Virtual_DedicatedHost, error) {
	if mask == "" {
		mask = "mask[id,name,notes]"
	}
	return u.UserCustomerService.Id(userId).Mask(mask).GetDedicatedHosts()
}

func (u userManager) GetHardware(userId int, mask string) ([]datatypes.Hardware, error) {
	if mask == "" {
		mask = "mask[id,fullyQualifiedDomainName,primaryIpAddress,primaryBackendIpAddress,notes]"
	}
	return u.UserCustomerService.Id(userId).GetHardware()
}

func (u userManager) GetVirtualGuests(userId int, mask string) ([]datatypes.Virtual_Guest, error) {
	if mask == "" {
		mask = "mask[id,fullyQualifiedDomainName,primaryIpAddress,primaryBackendIpAddress,notes]"
	}
	return u.UserCustomerService.Id(userId).GetVirtualGuests()
}

// Create Softlayer portal user VPN overrides.
// int userId: The user customer identifier.
// int subnetId: The subnet identifier.
func (u userManager) CreateUserVpnOverride(userId int, subnetId int) (bool, error) {
	vpnOverrideTemplates := []datatypes.Network_Service_Vpn_Overrides{
		datatypes.Network_Service_Vpn_Overrides{
			UserId:   sl.Int(userId),
			SubnetId: sl.Int(subnetId),
		},
	}
	networkServiceVpnOverridesService := services.GetNetworkServiceVpnOverridesService(u.Session)
	return networkServiceVpnOverridesService.CreateObjects(vpnOverrideTemplates)
}

// Creates or updates a user’s VPN access privileges.
// int userId: The user customer identifier.
func (u userManager) UpdateVpnUser(userId int) (bool, error) {
	return u.UserCustomerService.Id(userId).UpdateVpnUser()
}

// Return user’s vpn accessible subnets.
// int userId: The user customer identifier.
func (u userManager) GetOverrides(userId int) ([]datatypes.Network_Service_Vpn_Overrides, error) {
	return u.UserCustomerService.Id(userId).GetOverrides()
}

// Delete overrides.
// int overrideId: Override to be deleted.
func (u userManager) DeleteUserVpnOverride(overrideId int) (bool, error) {
	networkServiceVpnOverridesService := services.GetNetworkServiceVpnOverridesService(u.Session)
	return networkServiceVpnOverridesService.Id(overrideId).DeleteObject()
}

// Update a user’s VPN password.
// int userID: The user customer identifier.
// string password: New password
func (u userManager) UpdateVpnPassword(userID int, password string) (bool, error) {
	return u.UserCustomerService.Id(userID).UpdateVpnPassword(&password)
}

// Returns user's API authentication keys.
// int keyId: The user customer identifier.
func (u userManager) GetApiAuthenticationKeys(userId int) ([]datatypes.User_Customer_ApiAuthentication, error) {
	return u.UserCustomerService.Id(userId).GetApiAuthenticationKeys()
}

// Remove user's API authentication key.
// int userId: The user customer identifier.
func (u userManager) RemoveApiAuthenticationKey(userId int) (bool, error) {
	apiAuthenticationKeys, err := u.GetApiAuthenticationKeys(userId)
	if err != nil {
		return false, slErrors.NewAPIError(T("Failed to get user's API authentication keys"), err.Error(), 2)
	}
	if len(apiAuthenticationKeys) == 0 {
		return true, nil
	}
	return u.UserCustomerService.RemoveApiAuthenticationKey(apiAuthenticationKeys[0].Id)
}
