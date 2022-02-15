package user

import (
	"strconv"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type DetailsCommand struct {
	UI          terminal.UI
	UserManager managers.UserManager
}

func NewDetailsCommand(ui terminal.UI, userManager managers.UserManager) (cmd *DetailsCommand) {
	return &DetailsCommand{
		UI:          ui,
		UserManager: userManager,
	}
}

func (cmd *DetailsCommand) Run(c *cli.Context) error {

	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	userId := c.Args()[0]
	id, err := strconv.Atoi(userId)
	if err != nil {
		return errors.NewInvalidUsageError(T("User ID should be a number."))
	}

	keys := c.Bool("keys")
	permissions := c.Bool("permissions")
	hardware := c.Bool("hardware")
	virtual := c.Bool("virtual")
	logins := c.Bool("logins")
	events := c.Bool("events")

	object_mask := "userStatus[name],parent[id,username],apiAuthenticationKeys[authenticationKey]"
	user, err := cmd.UserManager.GetUser(id, object_mask)
	if err != nil {
		return cli.NewExitError(T("Failed to show user detail.\n")+err.Error(), 2)
	}

	baseUserPrint(user, keys, cmd.UI)

	if permissions {
		perms, err := cmd.UserManager.GetUserPermissions(id)
		if err != nil {
			return cli.NewExitError(T("Failed to show user permissions.\n")+err.Error(), 2)
		}
		table := cmd.UI.Table([]string{T("keyName"), T("name")})
		for _, perm := range perms {
			table.Add(utils.FormatStringPointer(perm.KeyName), utils.FormatStringPointer(perm.Name))
		}
		table.Add("", "")
		table.Print()
	}

	if hardware {
		mask := "id, hardware, dedicatedHosts"
		access, err := cmd.UserManager.GetUser(id, mask)
		if err != nil {
			return cli.NewExitError(T("Failed to show hardware.\n")+err.Error(), 2)
		}

		table := cmd.UI.Table([]string{T("ID"), T("Name"), T("Cpus"), T("Memory"), T("Disk"), T("Created"), T("Dedicated Access")})
		for _, host := range access.DedicatedHosts {
			hostId := utils.FormatIntPointer(host.Id)
			hostFqdn := utils.FormatStringPointer(host.Name)
			hostCpu := utils.FormatIntPointer(host.CpuCount)
			hostMem := utils.FormatIntPointer(host.MemoryCapacity)
			hostDisk := utils.FormatIntPointer(host.DiskCapacity)
			hostCreated := utils.FormatSLTimePointer(host.CreateDate)
			table.Add(hostId, hostFqdn, hostCpu, hostMem, hostDisk, hostCreated)
		}
		table.Add("", "")
		table.Print()

		tableAccess := cmd.UI.Table([]string{T("ID"), T("Host Name"), T("Primary Public IP"), T("Primary Private IP"), T("Created")})
		for _, host := range access.Hardware {
			hostId := utils.FormatIntPointer(host.Id)
			hostFqdn := utils.FormatStringPointer(host.FullyQualifiedDomainName)
			hostPrimary := utils.FormatStringPointer(host.PrimaryIpAddress)
			hostPrivate := utils.FormatStringPointer(host.PrimaryBackendIpAddress)
			hostCreated := utils.FormatSLTimePointer(host.ProvisionDate)
			tableAccess.Add(hostId, hostFqdn, hostPrimary, hostPrivate, hostCreated)
		}
		tableAccess.Add("", "")
		tableAccess.Print()
	}

	if virtual {
		mask := "id, virtualGuests"
		access, err := cmd.UserManager.GetUser(id, mask)
		if err != nil {
			return cli.NewExitError(T("Failed to show virual server.\n")+err.Error(), 2)
		}

		tableAccess := cmd.UI.Table([]string{T("ID"), T("Host Name"), T("Primary Public IP"), T("Primary Private IP"), T("Created")})
		for _, host := range access.VirtualGuests {
			hostId := utils.FormatIntPointer(host.Id)
			hostFqdn := utils.FormatStringPointer(host.FullyQualifiedDomainName)
			hostPrimary := utils.FormatStringPointer(host.PrimaryIpAddress)
			hostPrivate := utils.FormatStringPointer(host.PrimaryBackendIpAddress)
			hostCreated := utils.FormatSLTimePointer(host.ProvisionDate)
			tableAccess.Add(hostId, hostFqdn, hostPrimary, hostPrivate, hostCreated)
		}
		tableAccess.Add("", "")
		tableAccess.Print()
	}

	if logins {
		var t time.Time
		loginLog, err := cmd.UserManager.GetLogins(id, t)
		if err != nil {
			return cli.NewExitError(T("Failed to show login history.\n")+err.Error(), 2)
		}

		table := cmd.UI.Table([]string{T("Date"), T("IP Address"), T("Successful Login?")})
		for _, login := range loginLog {
			loginData := utils.FormatSLTimePointer(login.CreateDate)
			loginIp := utils.FormatStringPointer(login.IpAddress)
			loginSucc := utils.FormatBoolPointer(login.SuccessFlag)

			table.Add(loginData, loginIp, loginSucc)
		}
		table.Add("", "")
		table.Print()
	}

	if events {
		var t time.Time
		events, err := cmd.UserManager.GetEvents(id, t)
		if err != nil {
			return cli.NewExitError(T("Failed to show event log.\n")+err.Error(), 2)
		}

		table := cmd.UI.Table([]string{T("Date"), T("Type"), T("IP Address"), T("Label"), T("Username")})
		for _, event := range events {
			eventData := utils.FormatSLTimePointer(event.EventCreateDate)
			eventName := utils.FormatStringPointer(event.EventName)
			eventIp := utils.FormatStringPointer(event.IpAddress)
			eventLabel := utils.FormatStringPointer(event.Label)
			eventUsername := utils.FormatStringPointer(event.Username)
			table.Add(eventData, eventName, eventIp, eventLabel, eventUsername)
		}
		table.Add("", "")
		table.Print()
	}

	return nil

}

func baseUserPrint(user datatypes.User_Customer, keys bool, ui terminal.UI) {
	table := ui.Table([]string{T("name"), T("value")})
	table.Add(T("ID"), utils.FormatIntPointer(user.Id))
	table.Add(T("Username"), utils.FormatStringPointer(user.Username))

	if keys {
		for _, key := range user.ApiAuthenticationKeys {
			table.Add(T("APIKEY"), utils.FormatStringPointer(key.AuthenticationKey))
		}
	}

	table.Add(T("Name"), utils.FormatStringPointer(user.FirstName)+" "+utils.FormatStringPointer(user.LastName))
	table.Add(T("Email"), utils.FormatStringPointer(user.Email))
	table.Add(T("OpenID"), utils.FormatStringPointer(user.OpenIdConnectUserName))
	table.Add(T("Address"), utils.FormatStringPointer(user.Address1)+" "+utils.FormatStringPointer(user.Address2)+" "+utils.FormatStringPointer(user.City)+" "+utils.FormatStringPointer(user.State)+" "+utils.FormatStringPointer(user.Country)+" "+utils.FormatStringPointer(user.PostalCode))
	table.Add(T("Company"), utils.FormatStringPointer(user.CompanyName))
	table.Add(T("Created"), utils.FormatSLTimePointer(user.CreateDate))
	table.Add(T("Phone Number"), utils.FormatStringPointer(user.OfficePhone))

	if user.Parent != nil {
		table.Add(T("Parent User"), utils.FormatStringPointer(user.Parent.Username))
	}

	if user.UserStatus != nil {
		table.Add(T("Status"), utils.FormatStringPointer(user.UserStatus.Name))
	}

	table.Add(T("PPTP VPN"), utils.FormatBoolPointer(user.PptpVpnAllowedFlag))
	table.Add(T("SSL VPN"), utils.FormatBoolPointer(user.SslVpnAllowedFlag))

	if len(user.SuccessfulLogins) != 0 {
		loginString := user.SuccessfulLogins[0].CreateDate.String() + " From: " + utils.FormatStringPointer(user.SuccessfulLogins[0].IpAddress)
		table.Add(T("Last Login"), loginString)
	}

	if len(user.UnsuccessfulLogins) != 0 {
		unloginString := user.UnsuccessfulLogins[0].CreateDate.String() + " From: " + utils.FormatStringPointer(user.UnsuccessfulLogins[0].IpAddress)
		table.Add(T("Last Failed Login"), unloginString)
	}
	table.Add("", "")
	table.Print()
}

func UserDetailMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_USER_NAME,
		Name:        CMD_USER_DETAIL_NAME,
		Description: T("User details"),
		Usage:       "${COMMAND_NAME} sl user detail IDENTIFIER [OPTIONS]",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "keys",
				Usage: T("Show the users API key"),
			},
			cli.BoolFlag{
				Name:  "permissions",
				Usage: T("Display permissions assigned to this user. Master users do not show permissions"),
			},
			cli.BoolFlag{
				Name:  "hardware",
				Usage: T("Display hardware this user has access to"),
			},
			cli.BoolFlag{
				Name:  "virtual",
				Usage: T("Display virtual guests this user has access to"),
			},
			cli.BoolFlag{
				Name:  "logins",
				Usage: T("Show login history of this user for the last 24 hours"),
			},
			cli.BoolFlag{
				Name:  "events",
				Usage: T("Show audit log for this user"),
			},
		},
	}
}