package autoscale

import (
	"bytes"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CreateCommand struct {
	UI               terminal.UI
	AutoScaleManager managers.AutoScaleManager
}

func NewCreateCommand(ui terminal.UI, autoScaleManager managers.AutoScaleManager) (cmd *CreateCommand) {
	return &CreateCommand{
		UI:               ui,
		AutoScaleManager: autoScaleManager,
	}
}

func (cmd *CreateCommand) Run(c *cli.Context) error {
	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	datacenterInput := c.String("datacenter")
	datacenter, err := cmd.AutoScaleManager.GetDatacenterByName(datacenterInput, "shortname")
	if err != nil {
		return cli.NewExitError(T("Failed to get Datacenter {{.datacenter}}.\n", map[string]interface{}{"datacenter": datacenterInput})+err.Error(), 2)
	}
	if len(datacenter) > 1 {
		return cli.NewExitError(T("Failed to get Datacenter {{.datacenter}}.\n", map[string]interface{}{"datacenter": datacenterInput}), 2)
	}
	if len(datacenter) == 0 {
		datacenter, err = cmd.AutoScaleManager.GetDatacenterByName(datacenterInput, "longname")
		if err != nil {
			return cli.NewExitError(T("Failed to get Datacenter {{.datacenter}}.\n", map[string]interface{}{"datacenter": datacenterInput})+err.Error(), 2)
		}
		if len(datacenter) != 1 {
			return cli.NewExitError(T("Failed to get Datacenter {{.datacenter}}.\n", map[string]interface{}{"datacenter": datacenterInput}), 2)
		}
	}

	disks := c.IntSlice("disk")
	numDisk := 0
	block := []datatypes.Virtual_Guest_Block_Device{}
	for _, disk := range disks {
		// disk 1 is reserved to boot
		if numDisk == 1 {
			numDisk++
		}
		numberDisk := strconv.Itoa(numDisk)
		block = append(block,
			datatypes.Virtual_Guest_Block_Device{
				DiskImage: &datatypes.Virtual_Disk_Image{
					Capacity: sl.Int(disk),
				},
				Device: sl.String(numberDisk),
			},
		)
		numDisk++
	}

	virtualGuestMemberTemplate := datatypes.Virtual_Guest{
		Domain:   sl.String(c.String("domain")),
		Hostname: sl.String(c.String("hostname")),
		Datacenter: &datatypes.Location{
			Id: sl.Int(*datacenter[0].Id),
		},
		OperatingSystemReferenceCode: sl.String(c.String("os")),
		StartCpus:                    sl.Int(c.Int("cpu")),
		MaxMemory:                    sl.Int(c.Int("memory")),
		BlockDevices:                 block,
		LocalDiskFlag:                sl.Bool(false),
		HourlyBillingFlag:            sl.Bool(true),
		PrivateNetworkOnlyFlag:       sl.Bool(false),
		NetworkComponents: []datatypes.Virtual_Guest_Network_Component{
			datatypes.Virtual_Guest_Network_Component{
				MaxSpeed: sl.Int(100),
			},
		},
		TypeId:       sl.Int(1),
		NetworkVlans: []datatypes.Network_Vlan{},
	}

	autoSacaleGroupTemplate := datatypes.Scale_Group{
		Name:                       sl.String(c.String("name")),
		Cooldown:                   sl.Int(c.Int("cooldown")),
		MinimumMemberCount:         sl.Int(c.Int("min")),
		MaximumMemberCount:         sl.Int(c.Int("max")),
		RegionalGroupId:            sl.Int(c.Int("regional")),
		TerminationPolicyId:        sl.Int(c.Int("termination-policy")),
		SuspendedFlag:              sl.Bool(false),
		BalancedTerminationFlag:    sl.Bool(false),
		VirtualGuestMemberTemplate: &virtualGuestMemberTemplate,
		VirtualGuestMemberCount:    sl.Uint(0),
	}

	if c.IsSet("postinstall") {
		autoSacaleGroupTemplate.VirtualGuestMemberTemplate.PostInstallScriptUri = sl.String(c.String("postinstall"))
	}

	if c.IsSet("userdata") {
		userData := []datatypes.Virtual_Guest_Attribute{
			datatypes.Virtual_Guest_Attribute{Value: sl.String(c.String("userdata"))},
		}
		autoSacaleGroupTemplate.VirtualGuestMemberTemplate.UserData = userData
	}

	if c.IsSet("key") {
		keys := c.IntSlice("key")
		sshkeys := []datatypes.Security_Ssh_Key{}
		for _, key := range keys {
			sshkeys = append(sshkeys,
				datatypes.Security_Ssh_Key{
					Id: sl.Int(key),
				},
			)
		}
		autoSacaleGroupTemplate.VirtualGuestMemberTemplate.SshKeys = sshkeys
	}

	if c.IsSet("policy-relative") || c.IsSet("policy-amount") || c.IsSet("policy-name") {

		if !c.IsSet("policy-relative") {
			return errors.NewMissingInputError("--policy-relative")
		}

		if !c.IsSet("policy-amount") {
			return errors.NewMissingInputError("--policy-amount")
		}

		if !c.IsSet("policy-name") {
			return errors.NewMissingInputError("--policy-name")
		}

		autoSacaleGroupTemplate.Policies = []datatypes.Scale_Policy{
			datatypes.Scale_Policy{
				Name: sl.String(c.String("policy-name")),
				ScaleActions: []datatypes.Scale_Policy_Action_Scale{
					datatypes.Scale_Policy_Action_Scale{
						Amount:    sl.Int(c.Int("policy-amount")),
						ScaleType: sl.String(c.String("policy-relative")),
					},
				},
			},
		}
	}

	if !c.IsSet("f") && !c.IsSet("force") {
		confirm, err := cmd.UI.Confirm(T("This action will incur charges on your account. Continue?"))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	autoScale, err := cmd.AutoScaleManager.CreateScaleGroup(&autoSacaleGroupTemplate)
	if err != nil {
		return cli.NewExitError(T("Failed to create Auto Scale Group.\n")+err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("ID"), utils.FormatIntPointer(autoScale.Id))
	table.Add(T("Created"), utils.FormatSLTimePointer(autoScale.CreateDate))
	table.Add(T("Name"), utils.FormatStringPointer(autoScale.Name))
	//Virtual Guests Table
	buf := new(bytes.Buffer)
	virtualGuests := autoScale.VirtualGuestMembers
	virtualGuestsTable := terminal.NewTable(buf, []string{T("Id"), T("Domain"), T("hostname")})
	for _, virtualGuest := range virtualGuests {
		virtualGuestsTable.Add(
			utils.FormatIntPointer(virtualGuest.VirtualGuest.Id),
			utils.FormatStringPointer(virtualGuest.VirtualGuest.Domain),
			utils.FormatStringPointer(virtualGuest.VirtualGuest.Hostname),
		)
	}
	virtualGuestsTable.Print()
	table.Add(T("Virtual Guest Members"), buf.String())

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

func AutoScaleCreateMetaData() cli.Command {
	return cli.Command{
		Category:    "autoscale",
		Name:        "create",
		Description: T("Order/Create a scale group."),
		Usage: T(`${COMMAND_NAME} sl autoscale create [OPTIONS]

EXAMPLE: 
   ${COMMAND_NAME} sl autoscale create --name testcreate --datacenter ams01 -- domain mydomain.com --hostname myhostname --cooldown 3600 --min 2 --max 3 
   --regional 142 --termination-policy 2 -os CENTOS_7_64 --cpu 2 --memory 1024 --disk 25

   ${COMMAND_NAME} sl autoscale create --name testcreate --datacenter ams01 --domain mydomain.com --hostname myhostname --cooldown 3600 --min 1 --max 3 
   --regional 142 --termination-policy 2 -os CENTOS_7_64 --cpu 2 --memory 1024 --disk 25  --disk 30 --userdata CENTOS --policy-relative ABSOLUTE 
   --policy-name mypolicy --policy-amount 3 --postinstall https://mypostinstallscript.com --key 1111111 --key 2222222`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "name",
				Usage:    T("TEXT Scale group's name.  [required]"),
				Required: true,
			},
			cli.IntFlag{
				Name:     "cooldown",
				Usage:    T("INTEGER The number of seconds this group will wait after lastActionDate before performing another action.  [required]"),
				Required: true,
			},
			cli.IntFlag{
				Name:     "min",
				Usage:    T("INTEGER Set the minimum number of guests  [required]"),
				Required: true,
			},
			cli.IntFlag{
				Name:     "max",
				Usage:    T("INTEGER Set the maximum number of guests  [required]"),
				Required: true,
			},
			cli.StringFlag{
				Name:     "regional",
				Usage:    T("INTEGER The identifier of the regional group this scaling group is assigned to.  [required]"),
				Required: true,
			},
			cli.StringFlag{
				Name:  "postinstall",
				Usage: T("TEXT Post-install script to download"),
			},
			cli.StringFlag{
				Name:     "os",
				Usage:    T("TEXT OS install code. Tip: you can specify <OS>_LATEST  [required]"),
				Required: true,
			},
			cli.StringFlag{
				Name:     "datacenter",
				Usage:    T("TEXT Datacenter shortname  [required]"),
				Required: true,
			},
			cli.StringFlag{
				Name:     "hostname",
				Usage:    T("TEXT Host portion of the FQDN  [required]"),
				Required: true,
			},
			cli.StringFlag{
				Name:     "domain",
				Usage:    T("TEXT Domain portion of the FQDN  [required]"),
				Required: true,
			},
			cli.IntFlag{
				Name:     "cpu",
				Usage:    T("INTEGER Number of CPUs for new guests (existing not effected)  [required]"),
				Required: true,
			},
			cli.IntFlag{
				Name:     "memory",
				Usage:    T("INTEGER RAM in MB or GB for new guests (existing not effected)  [required]"),
				Required: true,
			},
			cli.StringFlag{
				Name:  "policy-relative",
				Usage: T("TEXT The type of scale to perform(ABSOLUTE, PERCENT, RELATIVE)."),
			},
			cli.StringFlag{
				Name:     "termination-policy",
				Usage:    T("TEXT The termination policy for the group(CLOSEST_TO_NEXT_CHARGE=1, NEWEST=2, OLDEST=3).  [required]"),
				Required: true,
			},
			cli.StringFlag{
				Name:  "policy-name",
				Usage: T("TEXT Collection of policies for this group. This can be empty."),
			},
			cli.IntFlag{
				Name:  "policy-amount",
				Usage: T("TEXT The number to scale by. This number has different meanings based on type."),
			},
			cli.StringFlag{
				Name:  "userdata",
				Usage: T("TEXT User defined metadata string"),
			},
			cli.IntSliceFlag{
				Name:  "key",
				Usage: T("TEXT SSH keys to add to the root user (multiple occurrence permitted)"),
			},
			cli.IntSliceFlag{
				Name:     "disk",
				Usage:    T("INTEGER Disk sizes (multiple occurrence permitted)  [required]"),
				Required: true,
			},
			metadata.OutputFlag(),
			metadata.ForceFlag(),
		},
	}
}
