package virtual

import (
	"errors"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"

	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
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

func VSDNSSyncMetaData() cli.Command {
	return cli.Command{
		Category:    "vs",
		Name:        "dns-sync",
		Description: T("Synchronize DNS records for a virtual server instance"),
		Usage: T(`${COMMAND_NAME} sl vs dns-sync IDENTIFIER [OPTIONS]
   Note: If you don't specify any arguments, it will attempt to update both the A
   and PTR records. If you don't want to update both records, you may use the
   -a or --ptr arguments to limit the records updated.
 
EXAMPLE:
   ${COMMAND_NAME} sl vs dns-sync 12345678 --a-record --ttl 3600
   This command synchronizes A record(IP V4 address) of virtual server instance with ID 12345678 to DNS server and sets ttl of this A record to 3600.
   ${COMMAND_NAME} sl vs dns-sync 12345678 --aaaa-record --ptr
   This command synchronizes both AAAA record(IP V6 address) and PTR record of virtual server instance with ID 12345678 to DNS server.`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "a,a-record",
				Usage: T("Sync the A record for the host"),
			},
			cli.BoolFlag{
				Name:  "aaaa-record",
				Usage: T("Sync the AAAA record for the host"),
			},
			cli.BoolFlag{
				Name:  "ptr",
				Usage: T("Sync the PTR record for the host"),
			},
			cli.IntFlag{
				Name:  "ttl",
				Usage: T("Sets the TTL for the A and/or PTR records, default is: 7200"),
			},
			metadata.ForceFlag(),
		},
	}
}