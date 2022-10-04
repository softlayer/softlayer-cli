package virtual

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"

	slErrors "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type DnsSyncCommand struct {
	*metadata.SoftlayerCommand
	VirtualServerManager managers.VirtualServerManager
	Command              *cobra.Command
	DNSManager           managers.DNSManager
	ARecord              bool
	AAAARecord           bool
	PTR                  bool
	TTL                  int
	Force                bool
}

func NewDnsSyncCommand(sl *metadata.SoftlayerCommand) (cmd *DnsSyncCommand) {
	thisCmd := &DnsSyncCommand{
		SoftlayerCommand:     sl,
		VirtualServerManager: managers.NewVirtualServerManager(sl.Session),
		DNSManager:           managers.NewDNSManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "dns-sync " + T("IDENTIFIER"),
		Short: T("Synchronize DNS records for a virtual server instance"),
		Long: T(`${COMMAND_NAME} sl vs dns-sync IDENTIFIER [OPTIONS]
   Note: If you don't specify any arguments, it will attempt to update both the A
   and PTR records. If you don't want to update both records, you may use the
   -a or --ptr arguments to limit the records updated.
 
EXAMPLE:
   ${COMMAND_NAME} sl vs dns-sync 12345678 --a-record --ttl 3600
   This command synchronizes A record(IP V4 address) of virtual server instance with ID 12345678 to DNS server and sets ttl of this A record to 3600.
   ${COMMAND_NAME} sl vs dns-sync 12345678 --aaaa-record --ptr
   This command synchronizes both AAAA record(IP V6 address) and PTR record of virtual server instance with ID 12345678 to DNS server.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	cobraCmd.Flags().BoolVarP(&thisCmd.ARecord, "a-record", "a", false, T("Sync the A record for the host"))
	cobraCmd.Flags().BoolVar(&thisCmd.AAAARecord, "aaaa-record", false, T("Sync the AAAA record for the host"))
	cobraCmd.Flags().BoolVar(&thisCmd.PTR, "ptr", false, T("Sync the PTR record for the host"))
	cobraCmd.Flags().IntVar(&thisCmd.TTL, "ttl", 7200, T("Sets the TTL for the A and/or PTR records, default is: 7200"))
	cobraCmd.Flags().BoolVarP(&thisCmd.Force, "force", "f", false, T("Force operation without confirmation"))

	return thisCmd
}

func (cmd *DnsSyncCommand) Run(args []string) error {

	vsID, err := utils.ResolveVirtualGuestId(args[0])
	if err != nil {
		return slErrors.NewInvalidSoftlayerIdInputError("Virtual server ID")
	}

	subs := map[string]interface{}{
		"VsID": vsID,
		"VsId": vsID,
		"Zone": "",
	}

	if !cmd.Force {
		confirm, err := cmd.UI.Confirm(T("Attempt to update DNS records for virtual server instance: {{.VsID}}. Continue?", subs))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	vs, err := cmd.VirtualServerManager.GetInstance(vsID, "id,globalIdentifier,fullyQualifiedDomainName,hostname,domain,primaryBackendIpAddress,primaryIpAddress,primaryNetworkComponent[id,primaryIpAddress,primaryVersion6IpAddressRecord[ipAddress]]")
	if err != nil {
		return slErrors.NewAPIError(T("Failed to get virtual server instance: {{.VsID}}.\n", subs), err.Error(), 2)
	}

	both := false
	syncA := cmd.ARecord
	syncAAAA := cmd.AAAARecord
	syncPtr := cmd.PTR
	ttl := cmd.TTL

	if !syncPtr && !syncA && !syncAAAA {
		both = true
	}
	zoneID, err := cmd.DNSManager.GetZoneIdFromName(utils.StringPointertoString(vs.Domain))
	if err != nil {
		subs["Zone"] = utils.StringPointertoString(vs.Domain)
		return slErrors.NewAPIError(T("Failed to get zone ID from zone name: {{.Zone}}.\n", subs), err.Error(), 2)
	}

	var multiErrors []error

	if both || syncA {
		err := cmd.DNSManager.SyncARecord(vs, zoneID, ttl)
		if err != nil {
			newError := errors.New(T("Failed to synchronize A record for virtual server instance: {{.VsId}}.\n", subs) + err.Error())
			multiErrors = append(multiErrors, newError)
		} else {
			cmd.UI.Ok()
			cmd.UI.Print(T("Synchronized A record for virtual server instance: {{.VsId}}.", subs))
		}
	}

	if both || syncPtr {
		err := cmd.DNSManager.SyncPTRRecord(vs, ttl)
		if err != nil {
			newError := errors.New(T("Failed to synchronize PTR record for virtual server instance: {{.VsId}}.\n", subs) + err.Error())
			multiErrors = append(multiErrors, newError)
		} else {
			cmd.UI.Ok()
			cmd.UI.Print(T("Synchronized PTR record for virtual server instance: {{.VsId}}.", subs))
		}
	}

	if syncAAAA {
		err := cmd.DNSManager.SyncAAAARecord(vs, zoneID, ttl)
		if err != nil {
			newError := errors.New(T("Failed to synchronize AAAA record for virtual server instance: {{.VsId}}.\n", subs) + err.Error())
			multiErrors = append(multiErrors, newError)
		} else {
			cmd.UI.Ok()
			cmd.UI.Print(T("Synchronized AAAA record for virtual server instance: {{.VsId}}.", subs))
		}
	}

	if len(multiErrors) > 0 {
		errorString := ""
		for _, theError := range multiErrors {
			errorString = fmt.Sprintf("%v\n%v", errorString, theError.Error())
		}
		return errors.New(errorString)
	}
	return nil
}
