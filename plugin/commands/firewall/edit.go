package firewall

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

const DELIMITER = "=========================================\n"

type EditCommand struct {
	*metadata.SoftlayerCommand
	FirewallManager managers.FirewallManager
	Command         *cobra.Command
}

func NewEditCommand(sl *metadata.SoftlayerCommand) (cmd *EditCommand) {
	thisCmd := &EditCommand{
		SoftlayerCommand: sl,
		FirewallManager:  managers.NewFirewallManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "edit " + T("IDENTIFIER"),
		Short: T("Edit firewall rules"),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *EditCommand) Run(args []string) error {
	firewallType, firewallID, err := cmd.FirewallManager.ParseFirewallID(args[0])
	if err != nil {
		return errors.NewAPIError(T("Failed to parse firewall ID : {{.FirewallID}}.\n", map[string]interface{}{"FirewallID": args[0]}), err.Error(), 1)
	}

	if firewallType == "multiVlan" {
		cmd.UI.Print(T("All multi vlan rules must be managed through the FortiGate dashboard using the provided credentials."))
		return nil
	}

	file, err := ioutil.TempFile("", "rules")
	if err != nil {
		log.Fatal(err)
	}

	if firewallType == "vlan" {
		origRules, err := cmd.FirewallManager.GetDedicatedFirewallRules(firewallID)
		if err != nil {
			return errors.NewAPIError(T("Failed to get dedicated firewall rules for {{.FirewallID}}.\n", map[string]interface{}{"FirewallID": firewallID}), err.Error(), 2)
		}
		_, err = OpenEditorForVlanRules(origRules, file.Name())
		if err != nil {
			return errors.NewAPIError(T("Failed to open editor for vlan rules: {{.FirewallID}}.\n", map[string]interface{}{"FirewallID": firewallID}), err.Error(), 2)
		}
		b, err := ioutil.ReadFile(file.Name())
		cmd.UI.Print(string(b))
		confirm, err := cmd.UI.Confirm(T("Would you like to submit the rules. Continue?"))
		if err != nil {
			return err
		}
		if !confirm {
			deleteTempFile(*file)
			cmd.UI.Print(T("Aborted."))
			return nil
		}
		editedRules, err := ParseVlanRulefile(string(b))
		if err != nil {
			return errors.NewAPIError(T("Failed to parse vlan rule file.\n"), err.Error(), 2)
		}
		_, err = cmd.FirewallManager.EditDedicatedFirewallRules(firewallID, editedRules)
		if err != nil {
			return errors.NewAPIError(T("Failed to edit dedicated firewall rules.\n"), err.Error(), 2)
		}
		deleteTempFile(*file)
		cmd.UI.Ok()
		cmd.UI.Print(T("Firewall {{.FirewallID}} was updated.", map[string]interface{}{"FirewallID": firewallID}))
	} else {
		origRules, err := cmd.FirewallManager.GetStandardFirewallRules(firewallID)
		if err != nil {
			return errors.NewAPIError(T("Failed to get standard firewall rules for {{.FirewallID}}.\n", map[string]interface{}{"FirewallID": firewallID}), err.Error(), 2)
		}
		_, err = OpenEditorForComponentRules(origRules, file.Name())
		if err != nil {
			return errors.NewAPIError(T("Failed to open editor for component rules:  {{.FirewallID}}.\n", map[string]interface{}{"FirewallID": firewallID}), err.Error(), 2)
		}
		b, err := ioutil.ReadFile(file.Name())
		cmd.UI.Print(string(b))
		confirm, err := cmd.UI.Confirm(T("Would you like to submit the rules. Continue?"))
		if err != nil {
			return err
		}
		if !confirm {
			deleteTempFile(*file)
			cmd.UI.Print(T("Aborted."))
			return nil
		}
		editedRules, err := ParseComponentRulefile(string(b))
		if err != nil {
			return errors.NewAPIError(T("Failed to parse component rule file.\n"), err.Error(), 2)
		}
		_, err = cmd.FirewallManager.EditStandardFirewallRules(firewallID, editedRules)
		if err != nil {
			return errors.NewAPIError(T("Failed to edit standard firewall rules.\n"), err.Error(), 2)
		}
		deleteTempFile(*file)
		cmd.UI.Ok()
		cmd.UI.Print(T("Firewall {{.FirewallID}} was updated.", map[string]interface{}{"FirewallID": firewallID}))
	}

	return nil
}

func openEditor(file string) error {
	cmd := exec.Command("nano", file) // #nosec
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	return err
}

func OpenEditorForComponentRules(origRules []datatypes.Network_Component_Firewall_Rule, filePath string) (*os.File, error) {
	tempFile := filePath
	f, err := os.Create(tempFile) // #nosec
	if err != nil {
		return nil, err
	}
	defer func() {
		if fErr := f.Close(); fErr != nil {
			log.Fatal(fErr)
		}
	}()

	if len(origRules) == 0 {
		_, writeErr := f.WriteString(DELIMITER)
		if writeErr != nil {
			return nil, writeErr
		}
		_, writeErr = f.WriteString(GetFormattedComponentRule(datatypes.Network_Component_Firewall_Rule{}))
		if writeErr != nil {
			return nil, writeErr
		}
	} else {
		for _, rule := range origRules {
			_, writeErr := f.WriteString(DELIMITER)
			if writeErr != nil {
				return nil, writeErr
			}
			_, writeErr = f.WriteString(GetFormattedComponentRule(rule))
			if writeErr != nil {
				return nil, writeErr
			}
		}
	}
	_, writeErr := f.WriteString(DELIMITER)
	if writeErr != nil {
		return nil, writeErr
	}
	editorErr := openEditor(tempFile)
	if editorErr != nil {
		return nil, editorErr
	}
	return f, err
}
func OpenEditorForVlanRules(origRules []datatypes.Network_Vlan_Firewall_Rule, filePath string) (*os.File, error) {
	tempFile := filePath
	f, err := os.Create(tempFile) // #nosec
	if err != nil {
		return nil, err
	}
	defer func() {
		if fErr := f.Close(); fErr != nil {
			log.Fatal(fErr)
		}
	}()
	if len(origRules) == 0 {
		_, writeErr := f.WriteString(DELIMITER)
		if writeErr != nil {
			return nil, writeErr
		}
		_, writeErr = f.WriteString(GetFormattedVlanRule(datatypes.Network_Vlan_Firewall_Rule{}))
		if writeErr != nil {
			return nil, writeErr
		}
	} else {
		for _, rule := range origRules {
			_, writeErr := f.WriteString(DELIMITER)
			if writeErr != nil {
				return nil, writeErr
			}
			_, writeErr = f.WriteString(GetFormattedVlanRule(rule))
			if writeErr != nil {
				return nil, writeErr
			}
		}
	}
	_, writeErr := f.WriteString(DELIMITER)
	if writeErr != nil {
		return nil, writeErr
	}
	editorErr := openEditor(tempFile)
	if editorErr != nil {
		return nil, editorErr
	}
	return f, err
}

func GetFormattedComponentRule(rule datatypes.Network_Component_Firewall_Rule) string {
	if rule.Action == nil {
		rule.Action = sl.String("permit")
	}
	if rule.Protocol == nil {
		rule.Protocol = sl.String("tcp")
	}
	if rule.SourceIpAddress == nil {
		rule.SourceIpAddress = sl.String("any")
	}
	if rule.SourceIpSubnetMask == nil {
		rule.SourceIpSubnetMask = sl.String("255.255.255.255")
	}
	if rule.DestinationIpAddress == nil {
		rule.DestinationIpAddress = sl.String("any")
	}
	if rule.DestinationIpSubnetMask == nil {
		rule.DestinationIpSubnetMask = sl.String("255.255.255.255")
	}
	if rule.DestinationPortRangeStart == nil {
		rule.DestinationPortRangeStart = sl.Int(1)
	}
	if rule.DestinationPortRangeEnd == nil {
		rule.DestinationPortRangeEnd = sl.Int(1)
	}
	if rule.Version == nil {
		rule.Version = sl.Int(4)
	}
	return fmt.Sprintf("action: %s\nprotocol: %s\nsource_ip_address: %s\nsource_ip_subnet_mask: %s\ndestination_ip_address: %s\ndestination_ip_subnet_mask: %s\ndestination_port_range_start: %d\ndestination_port_range_end: %d\nversion: %d\n",
		*rule.Action, *rule.Protocol, *rule.SourceIpAddress, *rule.SourceIpSubnetMask, *rule.DestinationIpAddress, *rule.DestinationIpSubnetMask, *rule.DestinationPortRangeStart, *rule.DestinationPortRangeEnd, *rule.Version)
}

func GetFormattedVlanRule(rule datatypes.Network_Vlan_Firewall_Rule) string {
	if rule.Action == nil {
		rule.Action = sl.String("permit")
	}
	if rule.Protocol == nil {
		rule.Protocol = sl.String("tcp")
	}
	if rule.SourceIpAddress == nil {
		rule.SourceIpAddress = sl.String("any")
	}
	if rule.SourceIpSubnetMask == nil {
		rule.SourceIpSubnetMask = sl.String("255.255.255.255")
	}
	if rule.DestinationIpAddress == nil {
		rule.DestinationIpAddress = sl.String("any")
	}
	if rule.DestinationIpSubnetMask == nil {
		rule.DestinationIpSubnetMask = sl.String("255.255.255.255")
	}
	if rule.DestinationPortRangeStart == nil {
		rule.DestinationPortRangeStart = sl.Int(1)
	}
	if rule.DestinationPortRangeEnd == nil {
		rule.DestinationPortRangeEnd = sl.Int(1)
	}
	if rule.Version == nil {
		rule.Version = sl.Int(4)
	}
	return fmt.Sprintf("action: %s\nprotocol: %s\nsource_ip_address: %s\nsource_ip_subnet_mask: %s\ndestination_ip_address: %s\ndestination_ip_subnet_mask: %s\ndestination_port_range_start: %d\ndestination_port_range_end: %d\nversion: %d\n",
		*rule.Action, *rule.Protocol, *rule.SourceIpAddress, *rule.SourceIpSubnetMask, *rule.DestinationIpAddress, *rule.DestinationIpSubnetMask, *rule.DestinationPortRangeStart, *rule.DestinationPortRangeEnd, *rule.Version)
}

func ParseVlanRulefile(content string) ([]datatypes.Network_Vlan_Firewall_Rule, error) {
	rules := strings.Split(content, DELIMITER)
	parsedRules := []datatypes.Network_Vlan_Firewall_Rule{}
	order := 1
	for _, rule := range rules {
		if strings.Trim(rule, " ") == "" {
			continue
		}
		parsedRule := datatypes.Network_Vlan_Firewall_Rule{}
		parsedRule.OrderValue = sl.Int(order)
		order = order + 1
		lines := strings.Split(rule, "\n")
		for _, line := range lines {
			if strings.Trim(line, " ") == "" {
				continue
			}
			keyValue := strings.Split(line, ":")
			key := strings.Trim(keyValue[0], " ")
			value := strings.Trim(keyValue[1], " ")
			if key == "action" {
				parsedRule.Action = sl.String(value)
			} else if key == "protocol" {
				parsedRule.Protocol = sl.String(value)
			} else if key == "source_ip_address" {
				parsedRule.SourceIpAddress = sl.String(value)
			} else if key == "source_ip_subnet_mask" {
				parsedRule.SourceIpSubnetMask = sl.String(value)
			} else if key == "destination_ip_address" {
				parsedRule.DestinationIpAddress = sl.String(value)
			} else if key == "destination_ip_subnet_mask" {
				parsedRule.DestinationIpSubnetMask = sl.String(value)
			} else if key == "destination_port_range_start" {
				startPort, err := strconv.Atoi(value)
				if err != nil {
					return nil, errors.New(T("Failed to parse destination port range start. \n") + err.Error())
				}
				parsedRule.DestinationPortRangeStart = sl.Int(startPort)
			} else if key == "destination_port_range_end" {
				endPort, err := strconv.Atoi(value)
				if err != nil {
					return nil, errors.New(T("Failed to parse destination port range end. \n") + err.Error())
				}
				parsedRule.DestinationPortRangeEnd = sl.Int(endPort)
			} else if key == "version" {
				version, err := strconv.Atoi(value)
				if err != nil {
					return nil, errors.New(T("Failed to parse version. \n") + err.Error())
				}
				parsedRule.Version = sl.Int(version)
			}
		}
		parsedRules = append(parsedRules, parsedRule)
	}
	return parsedRules, nil
}

func ParseComponentRulefile(content string) ([]datatypes.Network_Component_Firewall_Rule, error) {
	rules := strings.Split(content, DELIMITER)
	parsedRules := []datatypes.Network_Component_Firewall_Rule{}
	order := 1
	for _, rule := range rules {
		if strings.Trim(rule, " ") == "" {
			continue
		}
		parsedRule := datatypes.Network_Component_Firewall_Rule{}
		parsedRule.OrderValue = sl.Int(order)
		order = order + 1
		lines := strings.Split(rule, "\n")
		for _, line := range lines {
			if strings.Trim(line, " ") == "" {
				continue
			}
			keyValue := strings.Split(line, ":")
			key := strings.Trim(keyValue[0], " ")
			value := strings.Trim(keyValue[1], " ")
			if key == "action" {
				parsedRule.Action = sl.String(value)
			} else if key == "protocol" {
				parsedRule.Protocol = sl.String(value)
			} else if key == "source_ip_address" {
				parsedRule.SourceIpAddress = sl.String(value)
			} else if key == "source_ip_subnet_mask" {
				parsedRule.SourceIpSubnetMask = sl.String(value)
			} else if key == "destination_ip_address" {
				parsedRule.DestinationIpAddress = sl.String(value)
			} else if key == "destination_ip_subnet_mask" {
				parsedRule.DestinationIpSubnetMask = sl.String(value)
			} else if key == "destination_port_range_start" {
				startPort, err := strconv.Atoi(value)
				if err != nil {
					return nil, errors.New(T("Failed to parse destination port range start. \n") + err.Error())
				}
				parsedRule.DestinationPortRangeStart = sl.Int(startPort)
			} else if key == "destination_port_range_end" {
				endPort, err := strconv.Atoi(value)
				if err != nil {
					return nil, errors.New(T("Failed to parse destination port range end. \n") + err.Error())
				}
				parsedRule.DestinationPortRangeEnd = sl.Int(endPort)
			} else if key == "version" {
				version, err := strconv.Atoi(value)
				if err != nil {
					return nil, errors.New(T("Failed to parse version. \n") + err.Error())
				}
				parsedRule.Version = sl.Int(version)
			}
		}
		parsedRules = append(parsedRules, parsedRule)
	}
	return parsedRules, nil
}

func deleteTempFile(file os.File) {
	if fErr := file.Close(); fErr != nil {
		log.Fatal(fErr)
	}
	fErr := os.Remove(file.Name())
	if fErr != nil {
		log.Fatal(fErr)
	}
}
