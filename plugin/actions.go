package plugin

import (
	"fmt"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/licenses"

	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration/core_config"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/account"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/callapi"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dedicatedhost"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dns"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/file"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/firewall"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/globalip"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/image"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/ipsec"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/loadbal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/order"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/placementgroup"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/security"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/securitygroup"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/subnet"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/tags"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/ticket"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/user"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/vlan"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

func GetCommandAcionBindings(context plugin.PluginContext, ui terminal.UI, session *session.Session) map[string]func(c *cli.Context) error {
	virtualServerManager := managers.NewVirtualServerManager(session)
	imageManager := managers.NewImageManager(session)
	networkManager := managers.NewNetworkManager(session)
	firewallManager := managers.NewFirewallManager(session)
	dnsManager := managers.NewDNSManager(session)
	ipsecManager := managers.NewIPSECManager(session)

	hardwareManager := managers.NewHardwareServerManager(session)
	orderManager := managers.NewOrderManager(session)
	callAPIManager := managers.NewCallAPIManager(session)
	licensesManager := managers.NewLicensesManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{

		// firewall - 5
		NS_FIREWALL_NAME + "-" + CMD_FW_ADD_NAME: func(c *cli.Context) error {
			return firewall.NewAddCommand(ui, firewallManager).Run(c)
		},
		NS_FIREWALL_NAME + "-" + CMD_FW_CANCEL_NAME: func(c *cli.Context) error {
			return firewall.NewCancelCommand(ui, firewallManager).Run(c)
		},
		NS_FIREWALL_NAME + "-" + CMD_FW_DETAIL_NAME: func(c *cli.Context) error {
			return firewall.NewDetailCommand(ui, firewallManager).Run(c)
		},
		NS_FIREWALL_NAME + "-" + CMD_FW_EDIT_NAME: func(c *cli.Context) error {
			return firewall.NewEditCommand(ui, firewallManager).Run(c)
		},
		NS_FIREWALL_NAME + "-" + CMD_FW_LIST_NAME: func(c *cli.Context) error {
			return firewall.NewListCommand(ui, firewallManager).Run(c)
		},

		// globalip - 5
		NS_GLOBALIP_NAME + "-" + CMD_GP_ASSIGN_NAME: func(c *cli.Context) error {
			return globalip.NewAssignCommand(ui, networkManager).Run(c)
		},
		NS_GLOBALIP_NAME + "-" + CMD_GP_CREATE_NAME: func(c *cli.Context) error {
			return globalip.NewCreateCommand(ui, networkManager).Run(c)
		},
		NS_GLOBALIP_NAME + "-" + CMD_GP_CANCEL_NAME: func(c *cli.Context) error {
			return globalip.NewCancelCommand(ui, networkManager).Run(c)
		},
		NS_GLOBALIP_NAME + "-" + CMD_GP_LIST_NAME: func(c *cli.Context) error {
			return globalip.NewListCommand(ui, networkManager).Run(c)
		},
		NS_GLOBALIP_NAME + "-" + CMD_GP_UNASSIGN_NAME: func(c *cli.Context) error {
			return globalip.NewUnassignCommand(ui, networkManager).Run(c)
		},

		//hardware -14
		NS_HARDWARE_NAME + "-" + CMD_HARDWARE_AUTHORIZE_STORAGE_NAME: func(c *cli.Context) error {
			return hardware.NewAuthorizeStorageCommand(ui, hardwareManager).Run(c)
		},
		NS_HARDWARE_NAME + "-" + CMD_HARDWARE_BILLING_NAME: func(c *cli.Context) error {
			return hardware.NewBillingCommand(ui, hardwareManager).Run(c)
		},
		NS_HARDWARE_NAME + "-" + CMD_HARDWARE_CANCEL_NAME: func(c *cli.Context) error {
			return hardware.NewCancelCommand(ui, hardwareManager).Run(c)
		},
		NS_HARDWARE_NAME + "-" + CMD_HARDWARE_CANCEL_REASONS_NAME: func(c *cli.Context) error {
			return hardware.NewCancelReasonsCommand(ui, hardwareManager).Run(c)
		},
		NS_HARDWARE_NAME + "-" + CMD_HARDWARE_CREATE_NAME: func(c *cli.Context) error {
			return hardware.NewCreateCommand(ui, hardwareManager, context).Run(c)
		},
		NS_HARDWARE_NAME + "-" + CMD_HARDWARE_CREATE_OPTIONS_NAME: func(c *cli.Context) error {
			return hardware.NewCreateOptionsCommand(ui, hardwareManager).Run(c)
		},
		NS_HARDWARE_NAME + "-" + CMD_HARDWARE_CREDENTIALS_NAME: func(c *cli.Context) error {
			return hardware.NewCredentialsCommand(ui, hardwareManager).Run(c)
		},
		NS_HARDWARE_NAME + "-" + CMD_HARDWARE_DETAIL_NAME: func(c *cli.Context) error {
			return hardware.NewDetailCommand(ui, hardwareManager).Run(c)
		},
		NS_HARDWARE_NAME + "-" + CMD_HARDWARE_EDIT_NAME: func(c *cli.Context) error {
			return hardware.NewEditCommand(ui, hardwareManager).Run(c)
		},
		NS_HARDWARE_NAME + "-" + CMD_HARDWARE_LIST_NAME: func(c *cli.Context) error {
			return hardware.NewListCommand(ui, hardwareManager).Run(c)
		},
		NS_HARDWARE_NAME + "-" + CMD_HARDWARE_POWER_CYCLE_NAME: func(c *cli.Context) error {
			return hardware.NewPowerCycleCommand(ui, hardwareManager).Run(c)
		},
		NS_HARDWARE_NAME + "-" + CMD_HARDWARE_POWER_OFF_NAME: func(c *cli.Context) error {
			return hardware.NewPowerOffCommand(ui, hardwareManager).Run(c)
		},
		NS_HARDWARE_NAME + "-" + CMD_HARDWARE_POWER_ON_NAME: func(c *cli.Context) error {
			return hardware.NewPowerOnCommand(ui, hardwareManager).Run(c)
		},
		NS_HARDWARE_NAME + "-" + CMD_HARDWARE_REBOOT_NAME: func(c *cli.Context) error {
			return hardware.NewRebootCommand(ui, hardwareManager).Run(c)
		},
		NS_HARDWARE_NAME + "-" + CMD_HARDWARE_RELOAD_NAME: func(c *cli.Context) error {
			return hardware.NewReloadCommand(ui, hardwareManager).Run(c)
		},
		NS_HARDWARE_NAME + "-" + CMD_HARDWARE_RESCUE_NAME: func(c *cli.Context) error {
			return hardware.NewRescueCommand(ui, hardwareManager).Run(c)
		},
		NS_HARDWARE_NAME + "-" + CMD_HARDWARE_UPDATE_FIRMWARE_NAME: func(c *cli.Context) error {
			return hardware.NewUpdateFirmwareCommand(ui, hardwareManager).Run(c)
		},
		NS_HARDWARE_NAME + "-toggle-ipmi": func(c *cli.Context) error {
			return hardware.NewToggleIPMICommand(ui, hardwareManager).Run(c)
		},
		NS_HARDWARE_NAME + "-bandwidth": func(c *cli.Context) error {
			return hardware.NewBandwidthCommand(ui, hardwareManager).Run(c)
		},
		NS_HARDWARE_NAME + "-storage": func(c *cli.Context) error {
			return hardware.NewStorageCommand(ui, hardwareManager).Run(c)
		},
		NS_HARDWARE_NAME + "-guests": func(c *cli.Context) error {
			return hardware.NewGuestsCommand(ui, hardwareManager).Run(c)
		},

		//ipsec - 11
		NS_IPSEC_NAME + "-" + CMD_IPSEC_CONFIG_NAME: func(c *cli.Context) error {
			return ipsec.NewConfigCommand(ui, ipsecManager).Run(c)
		},
		NS_IPSEC_NAME + "-" + CMD_IPSEC_CANCEL_NAME: func(c *cli.Context) error {
			return ipsec.NewCancelCommand(ui, ipsecManager).Run(c)
		},
		NS_IPSEC_NAME + "-" + CMD_IPSEC_ORDER_NAME: func(c *cli.Context) error {
			return ipsec.NewOrderCommand(ui, ipsecManager, context).Run(c)
		},
		NS_IPSEC_NAME + "-" + CMD_IPSEC_DETAIL_NAME: func(c *cli.Context) error {
			return ipsec.NewDetailCommand(ui, ipsecManager).Run(c)
		},
		NS_IPSEC_NAME + "-" + CMD_IPSEC_LIST_NAME: func(c *cli.Context) error {
			return ipsec.NewListCommand(ui, ipsecManager).Run(c)
		},
		NS_IPSEC_NAME + "-" + CMD_IPSEC_SUBNET_ADD_NAME: func(c *cli.Context) error {
			return ipsec.NewAddSubnetCommand(ui, ipsecManager).Run(c)
		},
		NS_IPSEC_NAME + "-" + CMD_IPSEC_SUBNET_REMOVE_NAME: func(c *cli.Context) error {
			return ipsec.NewRemoveSubnetCommand(ui, ipsecManager).Run(c)
		},
		NS_IPSEC_NAME + "-" + CMD_IPSEC_TRANS_ADD_NAME: func(c *cli.Context) error {
			return ipsec.NewAddTranslationCommand(ui, ipsecManager).Run(c)
		},
		NS_IPSEC_NAME + "-" + CMD_IPSEC_TRANS_REMOVE_NAME: func(c *cli.Context) error {
			return ipsec.NewRemoveTranslationCommand(ui, ipsecManager).Run(c)
		},
		NS_IPSEC_NAME + "-" + CMD_IPSEC_TRANS_UPDATE_NAME: func(c *cli.Context) error {
			return ipsec.NewUpdateTranslationCommand(ui, ipsecManager).Run(c)
		},
		NS_IPSEC_NAME + "-" + CMD_IPSEC_UPDATE_NAME: func(c *cli.Context) error {
			return ipsec.NewUpdateCommand(ui, ipsecManager).Run(c)
		},

		//securitygroup 12
		NS_SECURITYGROUP_NAME + "-" + CMD_SECURITYGROUP_CREATE_NAME: func(c *cli.Context) error {
			return securitygroup.NewCreateCommand(ui, networkManager).Run(c)
		},
		NS_SECURITYGROUP_NAME + "-" + CMD_SECURITYGROUP_DELETE_NAME: func(c *cli.Context) error {
			return securitygroup.NewDeleteCommand(ui, networkManager).Run(c)
		},
		NS_SECURITYGROUP_NAME + "-" + CMD_SECURITYGROUP_DETAIL_NAME: func(c *cli.Context) error {
			return securitygroup.NewDetailCommand(ui, networkManager).Run(c)
		},
		NS_SECURITYGROUP_NAME + "-" + CMD_SECURITYGROUP_EDIT_NAME: func(c *cli.Context) error {
			return securitygroup.NewEditCommand(ui, networkManager).Run(c)
		},
		NS_SECURITYGROUP_NAME + "-" + CMD_SECURITYGROUP_INTERFACE_ADD_NAME: func(c *cli.Context) error {
			return securitygroup.NewInterfaceAddCommand(ui, networkManager, virtualServerManager).Run(c)
		},
		NS_SECURITYGROUP_NAME + "-" + CMD_SECURITYGROUP_INTERFACE_LIST_NAME: func(c *cli.Context) error {
			return securitygroup.NewInterfaceListCommand(ui, networkManager).Run(c)
		},
		NS_SECURITYGROUP_NAME + "-" + CMD_SECURITYGROUP_INTERFACE_REMOVE_NAME: func(c *cli.Context) error {
			return securitygroup.NewInterfaceRemoveCommand(ui, networkManager, virtualServerManager).Run(c)
		},
		NS_SECURITYGROUP_NAME + "-" + CMD_SECURITYGROUP_LIST_NAME: func(c *cli.Context) error {
			return securitygroup.NewListCommand(ui, networkManager).Run(c)
		},
		NS_SECURITYGROUP_NAME + "-" + CMD_SECURITYGROUP_RULE_ADD_NAME: func(c *cli.Context) error {
			return securitygroup.NewRuleAddCommand(ui, networkManager).Run(c)
		},
		NS_SECURITYGROUP_NAME + "-" + CMD_SECURITYGROUP_RULE_EDIT_NAME: func(c *cli.Context) error {
			return securitygroup.NewRuleEditCommand(ui, networkManager).Run(c)
		},
		NS_SECURITYGROUP_NAME + "-" + CMD_SECURITYGROUP_RULE_LIST_NAME: func(c *cli.Context) error {
			return securitygroup.NewRuleListCommand(ui, networkManager).Run(c)
		},
		NS_SECURITYGROUP_NAME + "-" + CMD_SECURITYGROUP_RULE_REMOVE_NAME: func(c *cli.Context) error {
			return securitygroup.NewRuleRemoveCommand(ui, networkManager).Run(c)
		},

		//subnet 5
		NS_SUBNET_NAME + "-" + CMD_SUBNET_CANCEL_NAME: func(c *cli.Context) error {
			return subnet.NewCancelCommand(ui, networkManager).Run(c)
		},
		NS_SUBNET_NAME + "-" + CMD_SUBNET_CREATE_NAME: func(c *cli.Context) error {
			return subnet.NewCreateCommand(ui, networkManager).Run(c)
		},
		NS_SUBNET_NAME + "-" + CMD_SUBNET_DETAIL_NAME: func(c *cli.Context) error {
			return subnet.NewDetailCommand(ui, networkManager).Run(c)
		},
		NS_SUBNET_NAME + "-" + CMD_SUBNET_LIST_NAME: func(c *cli.Context) error {
			return subnet.NewListCommand(ui, networkManager).Run(c)
		},
		NS_SUBNET_NAME + "-" + CMD_SUBNET_LOOKUP_NAME: func(c *cli.Context) error {
			return subnet.NewLookupCommand(ui, networkManager).Run(c)
		},
		NS_SUBNET_NAME + "-" + CMD_SUBNET_ROUTE_NAME: func(c *cli.Context) error {
			return subnet.NewRouteCommand(ui, networkManager).Run(c)
		},
		NS_SUBNET_NAME + "-" + CMD_SUBNET_CLEAR_ROUTE_NAME: func(c *cli.Context) error {
			return subnet.NewClearRouteCommand(ui, networkManager).Run(c)
		},

		//virual server - 20
		NS_VIRTUAL_NAME + "-" + CMD_VS_AUTHORIZE_STORAGE_NAME: func(c *cli.Context) error {
			return virtual.NewAuthorizeStorageCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_CANCEL_NAME: func(c *cli.Context) error {
			return virtual.NewCancelCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_CAPTURE_NAME: func(c *cli.Context) error {
			return virtual.NewCaptureCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_CREATE_NAME: func(c *cli.Context) error {
			return virtual.NewCreateCommand(ui, virtualServerManager, imageManager, context).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_CREATE_HOST_NAME: func(c *cli.Context) error {
			return virtual.NewCreateHostCommand(ui, virtualServerManager, networkManager, context).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_CREATE_OPTIONS_NAME: func(c *cli.Context) error {
			return virtual.NewCreateOptionsCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_CREDENTIALS_NAME: func(c *cli.Context) error {
			return virtual.NewCredentialsCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_DETAIL_NAME: func(c *cli.Context) error {
			return virtual.NewDetailCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_DNS_SYNC_NAME: func(c *cli.Context) error {
			return virtual.NewDnsSyncCommand(ui, virtualServerManager, dnsManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_EDIT_NAME: func(c *cli.Context) error {
			return virtual.NewEditCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_LIST_NAME: func(c *cli.Context) error {
			return virtual.NewListCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_LIST_HOST_NAME: func(c *cli.Context) error {
			return virtual.NewListHostCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_MIGRATE_NAME: func(c *cli.Context) error {
			return virtual.NewMigrageCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_PAUSE_NAME: func(c *cli.Context) error {
			return virtual.NewPauseCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_POWER_OFF_NAME: func(c *cli.Context) error {
			return virtual.NewPowerOffCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_POWER_ON_NAME: func(c *cli.Context) error {
			return virtual.NewPowerOnCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_READY_NAME: func(c *cli.Context) error {
			return virtual.NewReadyCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_BILLING_NAME: func(c *cli.Context) error {
			return virtual.NewBillingCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_REBOOT_NAME: func(c *cli.Context) error {
			return virtual.NewRebootCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_RELOAD_NAME: func(c *cli.Context) error {
			return virtual.NewReloadCommand(ui, virtualServerManager, context).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_RESCUE_NAME: func(c *cli.Context) error {
			return virtual.NewRescueCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_RESUME_NAME: func(c *cli.Context) error {
			return virtual.NewResumeCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_UPGRADE_NAME: func(c *cli.Context) error {
			return virtual.NewUpgradeCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_CAPACITY_CREATE_OPTIONS: func(c *cli.Context) error {
			return virtual.NewCapacityCreateOptiosCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_CAPACITY_DETAIL_NAME: func(c *cli.Context) error {
			return virtual.NewCapacityDetailCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-bandwidth": func(c *cli.Context) error {
			return virtual.NewBandwidthCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-storage": func(c *cli.Context) error {
			return virtual.NewStorageCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-placementgroup-list": func(c *cli.Context) error {
			return virtual.NewPlacementGroupListCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-placementgroup-create-options": func(c *cli.Context) error {
			return virtual.NewPlacementGruopCreateOptionsCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-placementgroup-create": func(c *cli.Context) error {
			return virtual.NewVSPlacementGroupCreateCommand(ui, virtualServerManager, context).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_CAPACITY_LIST_NAME: func(c *cli.Context) error {
			return virtual.NewCapacityListCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_CAPACITY_CREATE_NAME: func(c *cli.Context) error {
			return virtual.NewCapacityCreateCommand(ui, virtualServerManager, context).Run(c)
		},
		NS_VIRTUAL_NAME + "-usage": func(c *cli.Context) error {
			return virtual.NewUsageCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-" + CMD_VS_PLACEMENT_DETAIL_NAME: func(c *cli.Context) error {
			return virtual.NewPlacementGroupDetailsCommand(ui, virtualServerManager).Run(c)
		},

		//order
		NS_ORDER_NAME + "-" + CMD_ORDER_CATEGORY_LIST_NAME: func(c *cli.Context) error {
			return order.NewCategoryListCommand(ui, orderManager).Run(c)
		},
		NS_ORDER_NAME + "-" + CMD_ORDER_ITEM_LIST_NAME: func(c *cli.Context) error {
			return order.NewItemListCommand(ui, orderManager).Run(c)
		},
		NS_ORDER_NAME + "-" + CMD_ORDER_PACKAGE_LIST_NAME: func(c *cli.Context) error {
			return order.NewPackageListCommand(ui, orderManager).Run(c)
		},
		NS_ORDER_NAME + "-" + CMD_ORDER_PACKAGE_LOCATION_NAME: func(c *cli.Context) error {
			return order.NewPackageLocationCommand(ui, orderManager).Run(c)
		},
		NS_ORDER_NAME + "-" + CMD_ORDER_PLACE_NAME: func(c *cli.Context) error {
			return order.NewPlaceCommand(ui, orderManager, context).Run(c)
		},
		NS_ORDER_NAME + "-" + CMD_ORDER_PLACE_QUOTE_NAME: func(c *cli.Context) error {
			return order.NewPlaceQuoteCommand(ui, orderManager, context).Run(c)
		},
		NS_ORDER_NAME + "-" + CMD_ORDER_PRESET_LIST_NAME: func(c *cli.Context) error {
			return order.NewPresetListCommand(ui, orderManager).Run(c)
		},


		//callapi
		NS_SL_NAME + "-" + CMD_CALLAPI_NAME: func(c *cli.Context) error {
			return callapi.NewCallAPICommand(ui, callAPIManager).Run(c)
		},
		
		//license
		NS_LICENSES_NAME + "-" + CMD_LICENSES_CREATE_OPTIONS_NAME: func(c *cli.Context) error {
			return licenses.NewLicensesOptionsCommand(ui, licensesManager).Run(c)
		},
	}

	// ibmcloud sl account
	accountCommands := account.GetCommandAcionBindings(context, ui, session)
	for name, action := range accountCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl dedicatedhost
	dedicatedhostCommands := dedicatedhost.GetCommandActionBindings(context, ui, session)
	for name, action := range dedicatedhostCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl dns
	dnsCommands := dns.GetCommandActionBindings(context, ui, session)
	for name, action := range dnsCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl vlan
	vlanCommands := vlan.GetCommandActionBindings(context, ui, session)
	for name, action := range vlanCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl block
	blockCommands := block.GetCommandAcionBindings(context, ui, session)
	for name, action := range blockCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl file
	fileCommands := file.GetCommandAcionBindings(context, ui, session)
	for name, action := range fileCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl image
	imageCommands := image.GetCommandActionBindings(context, ui, session)
	for name, action := range imageCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl tags
	tagsCommands := tags.GetCommandActionBindings(context, ui, session)
	for name, action := range tagsCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl loadbal
	loadbalCommands := loadbal.GetCommandAcionBindings(ui, session)
	for name, action := range loadbalCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl placement-group
	placementgroupCommands := placementgroup.GetCommandActionBindings(context, ui, session)
	for name, action := range placementgroupCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl ticket
	ticketCommands := ticket.GetCommandActionBindings(context, ui, session)
	for name, action := range ticketCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl security
	// ibmcloud sl sshkey
	// ibmcloud sl ssl
	userCommands := user.GetCommandActionBindings(ui, session)
	for name, action := range userCommands {
		CommandActionBindings[name] = action
	}

	for name, action := range security.GetCommandActionBindings(ui, session) {
		CommandActionBindings[name] = action
	}

	actionWithPreCheck := make(map[string]func(c *cli.Context) error)

	for name, action := range CommandActionBindings {
		actionCopy := action
		actionWithPreCheck[name] = func(c *cli.Context) (err error) {
			err = PreChecktRequirement(context, ui)
			if err != nil {
				return err
			}

			defer func() {
				// catch panic
				if recoverErr := recover(); recoverErr != nil {
					err = cli.NewExitError(fmt.Sprintf("%v", recoverErr), 1)
				}
				switch err.(type) {
				case *errors.InvalidUsageError:
					ui.Failed("%v", err)
					showCmdErr := cli.ShowCommandHelp(c, c.Command.Name)
					if showCmdErr != nil {
						fmt.Println(showCmdErr.Error())
					}
					err = cli.NewExitError("", 2)
				}
			}()
			err = actionCopy(c)
			return err
		}
	}
	return actionWithPreCheck
}

func PreChecktRequirement(context plugin.PluginContext, ui terminal.UI) error {
	var errorMessage error
	switch {
	case !context.IsLoggedIn():
		errorMessage = fmt.Errorf(T("Not logged in. Use '{{.Command}}' to log in.",
			map[string]interface{}{"Command": terminal.CommandColor(context.CLIName() + " login")}))
	case context.IAMToken() == "":
		errorMessage = fmt.Errorf(T("IAM token is required. Use '{{.Command}}' to log in.",
			map[string]interface{}{"Command": terminal.CommandColor(context.CLIName() + " login")}))
	case context.IMSAccountID() == "":
		errorMessage = fmt.Errorf(T("Current account is not linked to a Softlayer account. Use '{{.Command}}' to switch account.",
			map[string]interface{}{"Command": terminal.CommandColor(context.CLIName() + " target -c")}))
	case !core_config.NewIAMTokenInfo(context.IAMToken()).Accounts.Valid:
		errorMessage = fmt.Errorf(T("The linked Softlayer account is not validated. Use '{{.Command}}' to re-login.",
			map[string]interface{}{"Command": terminal.CommandColor(context.CLIName() + " login")}))
	}
	if errorMessage != nil {
		return utils.FailWithError(errorMessage.Error(), ui)
	}
	return nil
}
