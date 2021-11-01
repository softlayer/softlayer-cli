package plugin

import (
	"fmt"

	"github.com/softlayer/softlayer-go/session"
	"github.com/urfave/cli"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/configuration/core_config"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/callapi"
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
	userManager := managers.NewUserManager(session)
	callAPIManager := managers.NewCallAPIManager(session)
	ticketManager := managers.NewTicketManager(session)
	placeGroupManager := managers.NewPlaceGroupManager(session)

	CommandActionBindings := map[string]func(c *cli.Context) error{

		//dns - 9
		NS_DNS_NAME + "-" + CMD_DNS_IMPORT_NAME: func(c *cli.Context) error {
			return dns.NewImportCommand(ui, dnsManager).Run(c)
		},
		NS_DNS_NAME + "-" + CMD_DNS_RECORD_ADD_NAME: func(c *cli.Context) error {
			return dns.NewRecordAddCommand(ui, dnsManager).Run(c)
		},
		NS_DNS_NAME + "-" + CMD_DNS_RECORD_EDIT_NAME: func(c *cli.Context) error {
			return dns.NewRecordEditCommand(ui, dnsManager).Run(c)
		},
		NS_DNS_NAME + "-" + CMD_DNS_RECORD_LIST_NAME: func(c *cli.Context) error {
			return dns.NewRecordListCommand(ui, dnsManager).Run(c)
		},
		NS_DNS_NAME + "-" + CMD_DNS_RECORD_REMOVE_NAME: func(c *cli.Context) error {
			return dns.NewRecordRemoveCommand(ui, dnsManager).Run(c)
		},
		NS_DNS_NAME + "-" + CMD_DNS_ZONE_CREATE_NAME: func(c *cli.Context) error {
			return dns.NewZoneCreateCommand(ui, dnsManager).Run(c)
		},
		NS_DNS_NAME + "-" + CMD_DNS_ZONE_DELETE_NAME: func(c *cli.Context) error {
			return dns.NewZoneDeleteCommand(ui, dnsManager).Run(c)
		},
		NS_DNS_NAME + "-" + CMD_DNS_ZONE_LIST_NAME: func(c *cli.Context) error {
			return dns.NewZoneListCommand(ui, dnsManager).Run(c)
		},
		NS_DNS_NAME + "-" + CMD_DNS_ZONE_PRINT_NAME: func(c *cli.Context) error {
			return dns.NewZonePrintCommand(ui, dnsManager).Run(c)
		},

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

		// image - 6
		NS_IMAGE_NAME + "-" + CMD_IMG_DELETE_NAME: func(c *cli.Context) error {
			return image.NewDeleteCommand(ui, imageManager).Run(c)
		},
		NS_IMAGE_NAME + "-" + CMD_IMG_DETAIL_NAME: func(c *cli.Context) error {
			return image.NewDetailCommand(ui, imageManager).Run(c)
		},
		NS_IMAGE_NAME + "-" + CMD_IMG_EDIT_NAME: func(c *cli.Context) error {
			return image.NewEditCommand(ui, imageManager).Run(c)
		},
		NS_IMAGE_NAME + "-" + CMD_IMG_EXPORT_NAME: func(c *cli.Context) error {
			return image.NewExportCommand(ui, imageManager).Run(c)
		},
		NS_IMAGE_NAME + "-" + CMD_IMG_IMPORT_NAME: func(c *cli.Context) error {
			return image.NewImportCommand(ui, imageManager).Run(c)
		},
		NS_IMAGE_NAME + "-" + CMD_IMG_LIST_NAME: func(c *cli.Context) error {
			return image.NewListCommand(ui, imageManager).Run(c)
		},
		NS_IMAGE_NAME + "-" + CMD_IMG_DATACENTER_NAME: func(c *cli.Context) error {
			return image.NewDatacenterCommand(ui, imageManager).Run(c)
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
		NS_VIRTUAL_NAME + "-" + CMD_VS_CAPACITY_DETAIL_NAME: func(c *cli.Context) error {
			return virtual.NewCapacityDetailCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-bandwidth": func(c *cli.Context) error {
			return virtual.NewBandwidthCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-storage": func(c *cli.Context) error {
			return virtual.NewStorageCommand(ui, virtualServerManager).Run(c)
		},
		NS_VIRTUAL_NAME + "-" +CMD_VS_CAPACITY_LIST_NAME: func(c *cli.Context) error {
			return virtual.NewCapacityListCommand(ui, virtualServerManager).Run(c)
		},

		//Placement group
		NS_PLACEMENT_GROUP_NAME + "-" + CMD_PLACEMENT_GROUP_CREATE_NAME: func(c *cli.Context) error {
			return placementgroup.NewPlacementGroupCreateCommand(ui, placeGroupManager).Run(c)
		},
		NS_PLACEMENT_GROUP_NAME + "-" + CMD_PLACEMENT_GROUP_LIST_NAME: func(c *cli.Context) error {
			return placementgroup.NewPlacementGroupListCommand(ui, placeGroupManager).Run(c)
		},
		NS_PLACEMENT_GROUP_NAME + "-" + CMD_PLACEMENT_GROUP_DELETE_NAME: func(c *cli.Context) error {
			return placementgroup.NewPlacementGroupDeleteCommand(ui, placeGroupManager, virtualServerManager).Run(c)
		},
		NS_PLACEMENT_GROUP_NAME + "-" + CMD_PLACEMENT_GROUP_CREATE_OPTIONS_NAME: func(c *cli.Context) error {
			return placementgroup.NewPlacementGroupCreateOptionsCommand(ui, placeGroupManager).Run(c)
		},
		NS_PLACEMENT_GROUP_NAME + "-" + CMD_PLACEMENT_GROUP_DETAIL_NAME: func(c *cli.Context) error {
			return placementgroup.NewPlacementGroupDetailCommand(ui, placeGroupManager).Run(c)
		},

		//vlan 6
		NS_VLAN_NAME + "-" + CMD_VLAN_CREATE_NAME: func(c *cli.Context) error {
			return vlan.NewCreateCommand(ui, networkManager, context).Run(c)
		},
		NS_VLAN_NAME + "-" + CMD_VLAN_CANCEL_NAME: func(c *cli.Context) error {
			return vlan.NewCancelCommand(ui, networkManager).Run(c)
		},
		NS_VLAN_NAME + "-" + CMD_VLAN_DETAIL_NAME: func(c *cli.Context) error {
			return vlan.NewDetailCommand(ui, networkManager).Run(c)
		},
		NS_VLAN_NAME + "-" + CMD_VLAN_EDIT_NAME: func(c *cli.Context) error {
			return vlan.NewEditCommand(ui, networkManager).Run(c)
		},
		NS_VLAN_NAME + "-" + CMD_VLAN_LIST_NAME: func(c *cli.Context) error {
			return vlan.NewListCommand(ui, networkManager).Run(c)
		},
		NS_VLAN_NAME + "-" + CMD_VLAN_OPTIONS_NAME: func(c *cli.Context) error {
			return vlan.NewOptionsCommand(ui, networkManager).Run(c)
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

		//user
		NS_USER_NAME + "-" + CMD_USER_CREATE_NAME: func(c *cli.Context) error {
			return user.NewCreateCommand(ui, userManager).Run(c)
		},
		NS_USER_NAME + "-" + CMD_USER_LIST_NAME: func(c *cli.Context) error {
			return user.NewListCommand(ui, userManager).Run(c)
		},
		NS_USER_NAME + "-" + CMD_USER_DELETE_NAME: func(c *cli.Context) error {
			return user.NewDeleteCommand(ui, userManager).Run(c)
		},
		NS_USER_NAME + "-" + CMD_USER_DETAIL_NAME: func(c *cli.Context) error {
			return user.NewDetailsCommand(ui, userManager).Run(c)
		},
		NS_USER_NAME + "-" + CMD_USER_PERMISSIONS_NAME: func(c *cli.Context) error {
			return user.NewPermissionsCommand(ui, userManager).Run(c)
		},
		NS_USER_NAME + "-" + CMD_USER_EDIT_DETAILS_NAME: func(c *cli.Context) error {
			return user.NewEditCommand(ui, userManager).Run(c)
		},
		NS_USER_NAME + "-" + CMD_USER_EDIT_PERMISSIONS_NAME: func(c *cli.Context) error {
			return user.NewEditPermissionCommand(ui, userManager).Run(c)
		},

		//callapi
		NS_SL_NAME + "-" + CMD_CALLAPI_NAME: func(c *cli.Context) error {
			return callapi.NewCallAPICommand(ui, callAPIManager).Run(c)
		},

		//ticket
		NS_TICKET_NAME + "-" + CMD_TICKET_CREATE_NAME: func(c *cli.Context) error {
			return ticket.NewCreateStandardTicketCommand(ui, ticketManager).Run(c)
		},

		NS_TICKET_NAME + "-" + CMD_TICKET_ATTACH_NAME: func(c *cli.Context) error {
			return ticket.NewAttachDeviceTicketCommand(ui, ticketManager).Run(c)
		},

		NS_TICKET_NAME + "-" + CMD_TICKET_DETACH_NAME: func(c *cli.Context) error {
			return ticket.NewDetachDeviceTicketCommand(ui, ticketManager).Run(c)
		},

		NS_TICKET_NAME + "-" + CMD_TICKET_DETAIL_NAME: func(c *cli.Context) error {
			return ticket.NewDetailTicketCommand(ui, ticketManager, userManager).Run(c)
		},

		NS_TICKET_NAME + "-" + CMD_TICKET_UPDATE_NAME: func(c *cli.Context) error {
			return ticket.NewUpdateTicketCommand(ui, ticketManager).Run(c)
		},

		NS_TICKET_NAME + "-" + CMD_TICKET_SUBJECTS_NAME: func(c *cli.Context) error {
			return ticket.NewSubjectsTicketCommand(ui, ticketManager).Run(c)
		},

		NS_TICKET_NAME + "-" + CMD_TICKET_LIST_NAME: func(c *cli.Context) error {
			return ticket.NewListTicketCommand(ui, ticketManager).Run(c)
		},

		NS_TICKET_NAME + "-" + CMD_TICKET_UPLOAD_NAME: func(c *cli.Context) error {
			return ticket.NewUploadFileTicketCommand(ui, ticketManager).Run(c)
		},

		NS_TICKET_NAME + "-" + CMD_TICKET_SUMMARY_NAME: func(c *cli.Context) error {
			return ticket.NewSummaryTicketCommand(ui, ticketManager).Run(c)
		},
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

	// ibmcloud sl tags
	tagCommands := tags.GetCommandAcionBindings(ui, session)
	for name, action := range tagCommands {
		CommandActionBindings[name] = action
	}

	// ibmcloud sl loadbal
	loadbalCommands := loadbal.GetCommandAcionBindings(ui, session)
	for name, action := range loadbalCommands {
		CommandActionBindings[name] = action
	}


	// ibmcloud sl security
	// ibmcloud sl sshkey
	// ibmcloud sl ssl 
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
