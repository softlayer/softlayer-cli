package metadata

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

var (
	NS_ORDER_NAME  = "order"
	CMD_ORDER_NAME = "order"

	CMD_ORDER_CATEGORY_LIST_NAME    = "category-list"
	CMD_ORDER_ITEM_LIST_NAME        = "item-list"
	CMD_ORDER_PACKAGE_LIST_NAME     = "package-list"
	CMD_ORDER_PACKAGE_LOCATION_NAME = "package-locations"
	CMD_ORDER_PLACE_NAME            = "place"
	CMD_ORDER_PLACE_QUOTE_NAME      = "place-quote"
	CMD_ORDER_PRESET_LIST_NAME      = "preset-list"
)

func OrderNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  NS_SL_NAME,
		Name:        NS_ORDER_NAME,
		Description: T("Classic infrastructure Orders"),
	}
}

func OrderMetaData() cli.Command {
	return cli.Command{
		Category:    NS_SL_NAME,
		Name:        CMD_ORDER_NAME,
		Usage:       "${COMMAND_NAME} sl order",
		Description: T("Classic infrastructure Orders"),
		Subcommands: []cli.Command{
			OrderCategoryListMetaData(),
			OrderItemListMetaData(),
			OrderPackageListMetaData(),
			OrderPackageLocaionMetaData(),
			OrderPlaceMetaData(),
			OrderPlaceQuoteMetaData(),
			OrderPresetListMetaData(),
		},
	}
}

func OrderCategoryListMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_ORDER_NAME,
		Name:        CMD_ORDER_CATEGORY_LIST_NAME,
		Description: T("List the categories of a package"),
		Usage: T(`${COMMAND_NAME} sl order category-list [OPTIONS] PACKAGE_KEYNAME
	
EXAMPLE: 
   ${COMMAND_NAME} sl order category-list BARE_METAL_SERVER
   This command lists the categories of Bare Metal servers.`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "required",
				Usage: T("List only the required categories for the package"),
			},
			OutputFlag(),
		},
	}
}

func OrderItemListMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_ORDER_NAME,
		Name:        CMD_ORDER_ITEM_LIST_NAME,
		Description: T("List package items that are used for ordering"),
		Usage: T(`${COMMAND_NAME} sl order item-list [OPTIONS] PACKAGE_KEYNAME
	
EXAMPLE: 
   ${COMMAND_NAME} sl order item-list CLOUD_SERVER
   This command lists all items in the VSI package.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "keyword",
				Usage: T("A word (or string) that is used to filter item names"),
			},
			cli.StringFlag{
				Name:  "category",
				Usage: T("Category code that is used to filter items"),
			},
			OutputFlag(),
		},
	}
}

func OrderPackageListMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_ORDER_NAME,
		Name:        CMD_ORDER_PACKAGE_LIST_NAME,
		Description: T("List packages that can be ordered with the placeOrder API"),
		Usage: T(`${COMMAND_NAME} sl order package-list [OPTIONS]
		
EXAMPLE: 
   ${COMMAND_NAME} sl order package-list
   This command list out all packages for ordering.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "keyword",
				Usage: T("A word (or string) that is used to filter package names"),
			},
			cli.StringFlag{
				Name:  "package-type ",
				Usage: T("The keyname for the type of package. For example, BARE_METAL_CPU"),
			},
			OutputFlag(),
		},
	}
}

func OrderPackageLocaionMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_ORDER_NAME,
		Name:        CMD_ORDER_PACKAGE_LOCATION_NAME,
		Description: T("List datacenters a package can be ordered in"),
		Usage:       "${COMMAND_NAME} sl order package-locations PACKAGE_KEYNAME [OPTIONS]",
		Flags: []cli.Flag{
			OutputFlag(),
		},
	}
}

func OrderPlaceMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_ORDER_NAME,
		Name:        CMD_ORDER_PLACE_NAME,
		Description: T("Place or verify an order"),
		Usage: T(`${COMMAND_NAME} sl order place PACKAGE_KEYNAME LOCATION ORDER_ITEM1,ORDER_ITEM2,ORDER_ITEM3,ORDER_ITEM4... [OPTIONS]
	
	EXAMPLE: 
	${COMMAND_NAME} sl order place CLOUD_SERVER DALLAS13 GUEST_CORES_4,RAM_16_GB,REBOOT_REMOTE_CONSOLE,1_GBPS_PUBLIC_PRIVATE_NETWORK_UPLINKS,BANDWIDTH_0_GB_2,1_IP_ADDRESS,GUEST_DISK_100_GB_SAN,OS_UBUNTU_16_04_LTS_XENIAL_XERUS_MINIMAL_64_BIT_FOR_VSI,MONITORING_HOST_PING,NOTIFICATION_EMAIL_AND_TICKET,AUTOMATED_NOTIFICATION,UNLIMITED_SSL_VPN_USERS_1_PPTP_VPN_USER_PER_ACCOUNT,NESSUS_VULNERABILITY_ASSESSMENT_REPORTING --billing hourly --extras '{"virtualGuests": [{"hostname": "test", "domain": "softlayer.com"}]}' --complex-type SoftLayer_Container_Product_Order_Virtual_Guest
	This command orders an hourly VSI with 4 CPU, 16 GB RAM, 100 GB SAN disk, Ubuntu 16.04, and 1 Gbps public & private uplink in dal13`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "preset",
				Usage: T("The order preset (if required by the package)"),
			},
			cli.BoolFlag{
				Name:  "verify",
				Usage: T("Flag denoting whether to verify the order, or not place it"),
			},
			cli.IntFlag{
				Name:  "quantity",
				Usage: T("The quantity of the item being ordered. This value defaults to 1"),
			},
			cli.StringFlag{
				Name:  "billing",
				Usage: T("Billing rate [hourly|monthly], [default: hourly]"),
			},
			cli.StringFlag{
				Name:  "complex-type",
				Usage: T("The complex type of the order. The type begins with 'SoftLayer_Container_Product_Order_'"),
			},
			cli.StringFlag{
				Name:  "extras",
				Usage: T("JSON string that denotes extra data needs to be sent with the order"),
			},
			ForceFlag(),
			OutputFlag(),
		},
	}
}

func OrderPlaceQuoteMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_ORDER_NAME,
		Name:        CMD_ORDER_PLACE_QUOTE_NAME,
		Description: T("Place a quote"),
		Usage: T(`${COMMAND_NAME} sl order place-quote PACKAGE_KEYNAME LOCATION ORDER_ITEM1,ORDER_ITEM2,ORDER_ITEM3,ORDER_ITEM4... [OPTIONS]

    EXAMPLE: 
    ${COMMAND_NAME} sl order place-quote CLOUD_SERVER DALLAS13 GUEST_CORES_4,RAM_16_GB,REBOOT_REMOTE_CONSOLE,1_GBPS_PUBLIC_PRIVATE_NETWORK_UPLINKS,BANDWIDTH_0_GB_2,1_IP_ADDRESS,GUEST_DISK_100_GB_SAN,OS_UBUNTU_16_04_LTS_XENIAL_XERUS_MINIMAL_64_BIT_FOR_VSI,MONITORING_HOST_PING,NOTIFICATION_EMAIL_AND_TICKET,AUTOMATED_NOTIFICATION,UNLIMITED_SSL_VPN_USERS_1_PPTP_VPN_USER_PER_ACCOUNT,NESSUS_VULNERABILITY_ASSESSMENT_REPORTING --extras '{"virtualGuests": [{"hostname": "test", "domain": "softlayer.com"}]}' --complex-type SoftLayer_Container_Product_Order_Virtual_Guest --name "foobar" --send-email
    This command places a quote for a VSI with 4 CPU, 16 GB RAM, 100 GB SAN disk, Ubuntu 16.04, and 1 Gbps public & private uplink in datacenter dal13`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "preset",
				Usage: T("The order preset (if required by the package)"),
			},
			cli.StringFlag{
				Name:  "name",
				Usage: T("A custom name to be assigned to the quote (optional)"),
			},
			cli.BoolFlag{
				Name:  "send-email",
				Usage: T("The quote will be sent to the associated email address"),
			},
			cli.StringFlag{
				Name:  "complex-type",
				Usage: T("The complex type of the order. The type begins with 'SoftLayer_Container_Product_Order_'"),
			},
			cli.StringFlag{
				Name:  "extras",
				Usage: T("JSON string that denotes extra data needs to be sent with the order"),
			},
			OutputFlag(),
		},
	}
}

func OrderPresetListMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_ORDER_NAME,
		Name:        CMD_ORDER_PRESET_LIST_NAME,
		Description: T("List package presets"),
		Usage: T(`${COMMAND_NAME} sl order preset-list [OPTIONS] PACKAGE_KEYNAME

   EXAMPLE: 
	  ${COMMAND_NAME} sl order preset-list BARE_METAL_SERVER
	  This command lists the presets for Bare Metal servers.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "keyword",
				Usage: T("A word (or string) used to filter presets"),
			},
			OutputFlag(),
		},
	}
}
