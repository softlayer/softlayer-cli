package virtual

import (
	"errors"

	bmxErr "github.ibm.com/cgallo/softlayer-cli/plugin/errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
	slErrors "github.ibm.com/cgallo/softlayer-cli/plugin/errors"
	"github.ibm.com/cgallo/softlayer-cli/plugin/managers"
	"github.ibm.com/cgallo/softlayer-cli/plugin/utils"
)

type DnsSyncCommand struct {
	UI                   terminal.UI
	VirtualServerManager managers.VirtualServerManager
	DNSManager           managers.DNSManager
}

func NewDnsSyncCommand(ui terminal.UI, virtualServerManager managers.VirtualServerManager, dnsManager managers.DNSManager) (cmd *DnsSyncCommand) {
	return &DnsSyncCommand{
		UI:                   ui,
		VirtualServerManager: virtualServerManager,
		DNSManager:           dnsManager,
	}
}

func (cmd *DnsSyncCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return bmxErr.NewInvalidUsageError("This command requires one argument.")
	}
	vsID, err := utils.ResolveVirtualGuestId(c.Args()[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	if !c.IsSet("f") && !c.IsSet("force") {
		confirm, err := cmd.UI.Confirm(T("Attempt to update DNS records for virtual server instance: {{.VsID}}. Continue?", map[string]interface{}{"VsID": vsID}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	vs, err := cmd.VirtualServerManager.GetInstance(vsID, "id,globalIdentifier,fullyQualifiedDomainName,hostname,domain,primaryBackendIpAddress,primaryIpAddress,primaryNetworkComponent[id,primaryIpAddress,primaryVersion6IpAddressRecord[ipAddress]]")
	if err != nil {
		return cli.NewExitError(T("Failed to get virtual server instance: {{.VsID}}.\n", map[string]interface{}{"VsID": vsID})+err.Error(), 2)
	}

	both := false
	syncA := c.Bool("a")
	syncAAAA := c.Bool("aaaa-record")
	syncPtr := c.Bool("ptr")
	ttl := 7200
	if c.IsSet("ttl") {
		ttl = c.Int("ttl")
	}

	if !syncPtr && !syncA && !syncAAAA {
		both = true
	}
	zoneID, err := cmd.DNSManager.GetZoneIdFromName(utils.StringPointertoString(vs.Domain))
	if err != nil {
		return cli.NewExitError(T("Failed to get zone ID from zone name: {{.Zone}}.\n", map[string]interface{}{"Zone": utils.StringPointertoString(vs.Domain)})+err.Error(), 2)
	}

	var multiErrors []error

	if both || syncA {
		err := cmd.DNSManager.SyncARecord(vs, zoneID, ttl)
		if err != nil {
			newError := errors.New(T("Failed to synchronize A record for virtual server instance: {{.VsId}}.\n",
				map[string]interface{}{"VsId": vsID}) + err.Error())
			multiErrors = append(multiErrors, newError)
		} else {
			cmd.UI.Ok()
			cmd.UI.Print(T("Synchronized A record for virtual server instance: {{.VsId}}.", map[string]interface{}{"VsId": vsID}))
		}
	}

	if both || syncPtr {
		err := cmd.DNSManager.SyncPTRRecord(vs, ttl)
		if err != nil {
			newError := errors.New(T("Failed to synchronize PTR record for virtual server instance: {{.VsId}}.\n",
				map[string]interface{}{"VsId": vsID}) + err.Error())
			multiErrors = append(multiErrors, newError)
		} else {
			cmd.UI.Ok()
			cmd.UI.Print(T("Synchronized PTR record for virtual server instance: {{.VsId}}.", map[string]interface{}{"VsId": vsID}))
		}
	}

	if syncAAAA {
		err := cmd.DNSManager.SyncAAAARecord(vs, zoneID, ttl)
		if err != nil {
			newError := errors.New(T("Failed to synchronize AAAA record for virtual server instance: {{.VsId}}.\n",
				map[string]interface{}{"VsId": vsID}) + err.Error())
			multiErrors = append(multiErrors, newError)
		} else {
			cmd.UI.Ok()
			cmd.UI.Print(T("Synchronized AAAA record for virtual server instance: {{.VsId}}.", map[string]interface{}{"VsId": vsID}))
		}
	}
	if len(multiErrors) > 0 {
		return cli.NewExitError(cli.NewMultiError(multiErrors...).Error(), 2)
	}
	return nil
}
