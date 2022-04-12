package autoscale

import (
	"bytes"
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type DetailCommand struct {
	UI               terminal.UI
	AutoScaleManager managers.AutoScaleManager
	SecurityManager  managers.SecurityManager
}

func NewDetailCommand(ui terminal.UI, autoScaleManager managers.AutoScaleManager, securityManager managers.SecurityManager) (cmd *DetailCommand) {
	return &DetailCommand{
		UI:               ui,
		AutoScaleManager: autoScaleManager,
		SecurityManager:  securityManager,
	}
}

func (cmd *DetailCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	autoScaleGroupId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Autoscale Group ID")
	}

	autoScaleGroup, err := cmd.AutoScaleManager.GetScaleGroup(autoScaleGroupId, "")
	if err != nil {
		return cli.NewExitError(T("Failed to get AutoScale group.\n")+err.Error(), 2)
	}

	table := cmd.UI.Table([]string{T("Name"), T("Value")})
	table.Add(T("ID"), utils.FormatIntPointer(autoScaleGroup.Id))
	table.Add(T("Datacenter"), utils.FormatStringPointer(autoScaleGroup.RegionalGroup.Locations[0].LongName))
	table.Add(T("Termination"), utils.FormatStringPointer(autoScaleGroup.TerminationPolicy.Name))
	table.Add(T("Minimum Members"), utils.FormatIntPointer(autoScaleGroup.MinimumMemberCount))
	table.Add(T("Maximum Members"), utils.FormatIntPointer(autoScaleGroup.MaximumMemberCount))
	table.Add(T("Current Members"), utils.FormatUIntPointer(autoScaleGroup.VirtualGuestMemberCount))
	cooldown := strconv.Itoa(*autoScaleGroup.Cooldown) + " seconds"
	table.Add(T("Cooldown"), cooldown)
	table.Add(T("Last Action"), utils.FormatSLTimePointer(autoScaleGroup.LastActionDate))

	if autoScaleGroup.NetworkVlans != nil && len(autoScaleGroup.NetworkVlans) != 0 {
		buf := new(bytes.Buffer)
		networkVlanTable := terminal.NewTable(buf, []string{T("Network Type"), T("Vlan Name")})
		for _, networkVlan := range autoScaleGroup.NetworkVlans {
			vlanName := *networkVlan.NetworkVlan.PrimaryRouter.Hostname + "." + strconv.Itoa(*networkVlan.NetworkVlan.VlanNumber)
			networkVlanTable.Add(
				utils.FormatStringPointer(networkVlan.NetworkVlan.NetworkSpace),
				vlanName,
			)
		}
		networkVlanTable.Print()
		table.Add(T("Network Vlans"), buf.String())
	}

	//Virtual Guest Member Template Table
	buf := new(bytes.Buffer)
	virtualGuestMemberTemplate := autoScaleGroup.VirtualGuestMemberTemplate
	virtualGuestMemberTemplateTable := terminal.NewTable(buf, []string{T("Name"), T("Value")})
	virtualGuestMemberTemplateTable.Add(T("Hostname"), utils.FormatStringPointer(virtualGuestMemberTemplate.Hostname))
	virtualGuestMemberTemplateTable.Add(T("Domain"), utils.FormatStringPointer(virtualGuestMemberTemplate.Domain))
	virtualGuestMemberTemplateTable.Add(T("Core"), utils.FormatIntPointer(virtualGuestMemberTemplate.StartCpus))
	virtualGuestMemberTemplateTable.Add(T("Ram"), utils.FormatIntPointer(virtualGuestMemberTemplate.MaxMemory))
	if virtualGuestMemberTemplate.NetworkComponents != nil && len(virtualGuestMemberTemplate.NetworkComponents) != 0 {
		virtualGuestMemberTemplateTable.Add(T("Network"), utils.FormatIntPointer(virtualGuestMemberTemplate.NetworkComponents[0].MaxSpeed))
	}
	if virtualGuestMemberTemplate.SshKeys != nil && len(virtualGuestMemberTemplate.SshKeys) != 0 {
		for _, sshKey := range virtualGuestMemberTemplate.SshKeys {
			sshkeyData, err := cmd.SecurityManager.GetSSHKey(*sshKey.Id)
			if err != nil {
				return cli.NewExitError(T("Failed to get SSH key."), 2)
			}
			virtualGuestMemberTemplateTable.Add(T("SSH Key ")+strconv.Itoa(*sshKey.Id), utils.FormatStringPointer(sshkeyData.Label))
		}
	}
	if virtualGuestMemberTemplate.BlockDevices != nil && len(virtualGuestMemberTemplate.BlockDevices) != 0 {
		for _, disk := range virtualGuestMemberTemplate.BlockDevices {
			diskType := "SAN"
			if *virtualGuestMemberTemplate.LocalDiskFlag {
				diskType = "Local"
			}
			virtualGuestMemberTemplateTable.Add(diskType+T(" Disk ")+*disk.Device, utils.FormatIntPointer(disk.DiskImage.Capacity))
		}
	}
	virtualGuestMemberTemplateTable.Add(T("OS"), utils.FormatStringPointer(virtualGuestMemberTemplate.OperatingSystemReferenceCode))
	postInstall := "None"
	if virtualGuestMemberTemplate.PostInstallScriptUri != nil && *virtualGuestMemberTemplate.PostInstallScriptUri != "" {
		postInstall = *virtualGuestMemberTemplate.PostInstallScriptUri
	}
	virtualGuestMemberTemplateTable.Add(T("Post Install"), postInstall)
	virtualGuestMemberTemplateTable.Print()
	table.Add(T("Virtual Guest Member Template"), buf.String())

	//Policies Table
	if autoScaleGroup.Policies != nil && len(autoScaleGroup.Policies) != 0 {
		buf = new(bytes.Buffer)
		policiesTable := terminal.NewTable(buf, []string{T("Policy"), T("Cooldown")})
		for _, policy := range autoScaleGroup.Policies {
			if policy.Cooldown != nil {
				policiesTable.Add(*policy.Name, utils.FormatIntPointer(policy.Cooldown))
			} else {
				policiesTable.Add(*policy.Name, utils.FormatIntPointer(autoScaleGroup.Cooldown))
			}
		}
		policiesTable.Print()
		table.Add(T("Policies"), buf.String())
	}

	//Active Guests Table
	if autoScaleGroup.VirtualGuestMembers != nil && len(autoScaleGroup.VirtualGuestMembers) != 0 {
		buf = new(bytes.Buffer)
		activeGuestsTable := terminal.NewTable(buf, []string{T("Id"), T("Hostname"), T("Created")})
		for _, virtualGuest := range autoScaleGroup.VirtualGuestMembers {
			activeGuestsTable.Add(utils.FormatIntPointer(virtualGuest.VirtualGuest.Id), utils.FormatStringPointer(virtualGuest.VirtualGuest.Hostname), utils.FormatSLTimePointer(virtualGuest.VirtualGuest.ProvisionDate))
		}
		activeGuestsTable.Print()
		table.Add(T("Active Guests"), buf.String())
	}

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

func AutoScaleDetailMetaData() cli.Command {
	return cli.Command{
		Category:    "autoscale",
		Name:        "detail",
		Description: T("Get details of an Autoscale group."),
		Usage: T(`${COMMAND_NAME} sl autoscale detail IDENTIFIER

EXAMPLE: 
   ${COMMAND_NAME} sl autoscale detail IDENTIFIER`),
		Flags: []cli.Flag{
			metadata.OutputFlag(),
		},
	}
}
