package metadata

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

var (
	NS_HARDWARE_NAME  = "hardware"
	CMD_HARDWARE_NAME = "hardware"

	//sl-hardware
	CMD_HARDWARE_AUTHORIZE_STORAGE_NAME = "authorize-storage"
	CMD_HARDWARE_CANCEL_NAME            = "cancel"
	CMD_HARDWARE_CANCEL_REASONS_NAME    = "cancel-reasons"
	CMD_HARDWARE_CREATE_NAME            = "create"
	CMD_HARDWARE_CREATE_OPTIONS_NAME    = "create-options"
	CMD_HARDWARE_CREDENTIALS_NAME       = "credentials"
	CMD_HARDWARE_DETAIL_NAME            = "detail"
	CMD_HARDWARE_EDIT_NAME              = "edit"
	CMD_HARDWARE_LIST_NAME              = "list"
	CMD_HARDWARE_POWER_CYCLE_NAME       = "power-cycle"
	CMD_HARDWARE_POWER_OFF_NAME         = "power-off"
	CMD_HARDWARE_POWER_ON_NAME          = "power-on"
	CMD_HARDWARE_REBOOT_NAME            = "reboot"
	CMD_HARDWARE_RELOAD_NAME            = "reload"
	CMD_HARDWARE_RESCUE_NAME            = "rescue"
	CMD_HARDWARE_UPDATE_FIRMWARE_NAME   = "update-firmware"
)

func HardwareNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  NS_SL_NAME,
		Name:        NS_HARDWARE_NAME,
		Description: T("Classic infrastructure hardware servers"),
	}
}

func HardwareMetaData() cli.Command {
	return cli.Command{
		Category:    NS_SL_NAME,
		Name:        CMD_HARDWARE_NAME,
		Description: T("Classic infrastructure hardware servers"),
		Usage:       "${COMMAND_NAME} sl hardware",
		Subcommands: []cli.Command{
			HardwareAuthorizeStorageMataData(),
			HardwareCancelMetaData(),
			HardwareCancelReasonsMetaData(),
			HardwareCreateMetaData(),
			HardwareCreateOptionsMetaData(),
			HardwareCredentialsMetaData(),
			HardwareDetailMetaData(),
			HardwareEditMetaData(),
			HardwareListMetaData(),
			HardwarePowerCycleMetaData(),
			HardwarePowerOffMetaData(),
			HardwarePowerOnMetaData(),
			HardwarePowerRebootMetaData(),
			HardwareReloadMetaData(),
			HardwareRescueMetaData(),
			HardwareUpdateFirmwareMetaData(),
			HardwareToggleIPMIMetaData(),
		},
	}
}

func HardwareAuthorizeStorageMataData() cli.Command {
	return cli.Command{
		Category:    CMD_HARDWARE_NAME,
		Name:        CMD_HARDWARE_AUTHORIZE_STORAGE_NAME,
		Description: T("Authorize File and Block Storage to a Hardware Server"),
		Usage: T(`${COMMAND_NAME} sl hardware authorize-storage [OPTIONS] IDENTIFIER
	
EXAMPLE:
   ${COMMAND_NAME} sl hardware authorize-storage --username-storage SL01SL30-37 1234567
   Authorize File and Block Storage to a Hardware Server.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "u, username-storage",
				Usage: T("The storage username to be added to the hardware server."),
			},
			OutputFlag(),
		},
	}
}

func HardwareCancelMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_HARDWARE_NAME,
		Name:        CMD_HARDWARE_CANCEL_NAME,
		Description: T("Cancel a hardware server"),
		Usage:       "${COMMAND_NAME} sl hardware cancel IDENTIFIER [OPTIONS]",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "i,immediate",
				Usage: T("Cancels the server immediately (instead of on the billing anniversary)"),
			},
			cli.StringFlag{
				Name:  "r,reason",
				Usage: T("An optional cancellation reason. See '${COMMAND_NAME} sl hardware cancel-reasons' for a list of available options"),
			},
			cli.StringFlag{
				Name:  "c,comment",
				Usage: T("An optional comment to add to the cancellation ticket"),
			},
			ForceFlag(),
		},
	}
}

func HardwareCancelReasonsMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_HARDWARE_NAME,
		Name:        CMD_HARDWARE_CANCEL_REASONS_NAME,
		Description: T("Display a list of cancellation reasons"),
		Usage:       "${COMMAND_NAME} sl hardware cancel-reasons",
	}
}

func HardwareCreateMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_HARDWARE_NAME,
		Name:        CMD_HARDWARE_CREATE_NAME,
		Description: T("Order/create a hardware server"),
		Usage: `${COMMAND_NAME} sl hardware create [OPTIONS] 
	See '${COMMAND_NAME} sl hardware create-options' for valid options.`,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "H,hostname",
				Usage: T("Host portion of the FQDN[required]"),
			},
			cli.StringFlag{
				Name:  "D,domain",
				Usage: T("Domain portion of the FQDN[required]"),
			},
			cli.StringFlag{
				Name:  "s,size",
				Usage: T("Hardware size[required]"),
			},
			cli.StringFlag{
				Name:  "o,os",
				Usage: T("OS install code[required]"),
			},
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Datacenter shortname[required]"),
			},
			cli.IntFlag{
				Name:  "p,port-speed",
				Usage: T("Port speed[required]"),
			},
			cli.StringFlag{
				Name:  "b,billing",
				Usage: T("Billing rate, either hourly or monthly, default is hourly if not specified"),
			},
			cli.StringFlag{
				Name:  "i,post-install",
				Usage: T("Post-install script to download"),
			},
			cli.IntSliceFlag{
				Name:  "k,key",
				Usage: T("SSH keys to add to the root user, multiple occurrence allowed"),
			},
			cli.BoolFlag{
				Name:  "n,no-public",
				Usage: T("Private network only"),
			},
			cli.StringSliceFlag{
				Name:  "e,extra",
				Usage: T("Extra options, multiple occurrence allowed"),
			},
			cli.BoolFlag{
				Name:  "t,test",
				Usage: T("Do not actually create the virtual server"),
			},
			cli.StringFlag{
				Name:  "m,template",
				Usage: T("A template file that defaults the command-line options"),
			},
			cli.StringFlag{
				Name:  "x,export",
				Usage: T("Exports options to a template file"),
			},
			ForceFlag(),
		},
	}
}

func HardwareCreateOptionsMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_HARDWARE_NAME,
		Name:        CMD_HARDWARE_CREATE_OPTIONS_NAME,
		Description: T("Server order options for a given chassis"),
		Usage:       "${COMMAND_NAME} sl hardware create-options",
	}
}

func HardwareCredentialsMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_HARDWARE_NAME,
		Name:        CMD_HARDWARE_CREDENTIALS_NAME,
		Description: T("List hardware server credentials"),
		Usage:       "${COMMAND_NAME} sl hardware credentials IDENTIFIER",
		Flags: []cli.Flag{
			OutputFlag(),
		},
	}
}

func HardwareDetailMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_HARDWARE_NAME,
		Name:        CMD_HARDWARE_DETAIL_NAME,
		Description: T("Get details for a hardware server"),
		Usage:       "${COMMAND_NAME} sl hardware detail IDENTIFIER [OPTIONS]",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "p,passwords",
				Usage: T("Show passwords (check over your shoulder!)"),
			},
			cli.BoolFlag{
				Name:  "c,price",
				Usage: T("Show associated prices"),
			},
			OutputFlag(),
		},
	}
}

func HardwareEditMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_HARDWARE_NAME,
		Name:        CMD_HARDWARE_EDIT_NAME,
		Description: T("Edit hardware server details"),
		Usage:       "${COMMAND_NAME} sl hardware edit IDENTIFIER [OPTIONS]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "H,hostname",
				Usage: T("Host portion of the FQDN"),
			},
			cli.StringFlag{
				Name:  "D,domain",
				Usage: T("Domain portion of the FQDN"),
			},
			cli.StringSliceFlag{
				Name:  "g,tag",
				Usage: T("Tags to set or empty string to remove all (multiple occurrence permitted)."),
			},
			cli.StringFlag{
				Name:  "F,userfile",
				Usage: T("Read userdata from file"),
			},
			cli.StringFlag{
				Name:  "u,userdata",
				Usage: T("User defined metadata string"),
			},
			cli.IntFlag{
				Name:  "p,public-speed",
				Usage: T("Public port speed, options are: 0,10,100,1000,10000"),
			},
			cli.IntFlag{
				Name:  "v,private-speed",
				Usage: T("Private port speed, options are: 0,10,100,1000,10000"),
			},
		},
	}
}

func HardwareListMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_HARDWARE_NAME,
		Name:        CMD_HARDWARE_LIST_NAME,
		Description: T("List hardware servers"),
		Usage:       "${COMMAND_NAME} sl hardware list [OPTIONS]",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "c,cpu",
				Usage: T("Filter by number of CPU cores"),
			},
			cli.StringFlag{
				Name:  "D,domain",
				Usage: T("Filter by domain"),
			},
			cli.StringFlag{
				Name:  "H,hostname",
				Usage: T("Filter by hostname"),
			},
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Filter by datacenter"),
			},
			cli.IntFlag{
				Name:  "m,memory",
				Usage: T("Filter by memory in gigabytes"),
			},
			cli.IntFlag{
				Name:  "n,network",
				Usage: T("Filter by network port speed in Mbps"),
			},
			cli.StringSliceFlag{
				Name:  "g,tag",
				Usage: T("Filter by tags, multiple occurrence allowed"),
			},
			cli.StringFlag{
				Name:  "p,public-ip",
				Usage: T("Filter by public IP address"),
			},
			cli.StringFlag{
				Name:  "v,private-ip",
				Usage: T("Filter by private IP address"),
			},
			cli.IntFlag{
				Name:  "o,order",
				Usage: T("Filter by ID of the order which purchased hardware server"),
			},
			cli.StringFlag{
				Name:  "owner",
				Usage: T("Filter by ID of the owner"),
			},
			cli.StringFlag{
				Name:  "sortby",
				Usage: T("Column to sort by, default:hostname, option:id,guid,hostname,domain,public_ip,private_ip,cpu,memory,os,datacenter,status,ipmi_ip,created,created_by"),
			},
			cli.StringSliceFlag{
				Name:  "column",
				Usage: T("Column to display,  options are: id,hostname,domain,public_ip,private_ip,datacenter,status,guid,cpu,memory,os,ipmi_ip,created,created_by,tags. This option can be specified multiple times"),
			},
			cli.StringSliceFlag{
				Name:   "columns",
				Hidden: true,
			},
			OutputFlag(),
		},
	}
}

func HardwarePowerCycleMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_HARDWARE_NAME,
		Name:        CMD_HARDWARE_POWER_CYCLE_NAME,
		Description: T("Power cycle a server"),
		Usage:       "${COMMAND_NAME} sl hardware power-cycle IDENTIFIER [OPTIONS]",
		Flags: []cli.Flag{
			ForceFlag(),
		},
	}
}

func HardwarePowerOffMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_HARDWARE_NAME,
		Name:        CMD_HARDWARE_POWER_OFF_NAME,
		Description: T("Power off an active server"),
		Usage:       "${COMMAND_NAME} sl hardware power-off IDENTIFIER [OPTIONS]",
		Flags: []cli.Flag{
			ForceFlag(),
		},
	}
}

func HardwarePowerOnMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_HARDWARE_NAME,
		Name:        CMD_HARDWARE_POWER_ON_NAME,
		Description: T("Power on a server"),
		Usage:       "${COMMAND_NAME} sl hardware power-on IDENTIFIER",
	}
}

func HardwarePowerRebootMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_HARDWARE_NAME,
		Name:        CMD_HARDWARE_REBOOT_NAME,
		Description: T("Reboot an active server"),
		Usage:       "${COMMAND_NAME} sl hardware reboot IDENTIFIER [OPTIONS]",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "hard",
				Usage: T("Perform a hard reboot"),
			},
			cli.BoolFlag{
				Name:  "soft",
				Usage: T("Perform a soft reboot"),
			},
			ForceFlag(),
		},
	}
}

func HardwareReloadMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_HARDWARE_NAME,
		Name:        CMD_HARDWARE_RELOAD_NAME,
		Description: T("Reload operating system on a server"),
		Usage:       "${COMMAND_NAME} sl hardware reload IDENTIFIER [OPTIONS]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "i,postinstall",
				Usage: T("Post-install script to download, only HTTPS executes, HTTP leaves file in /root"),
			},
			cli.IntSliceFlag{
				Name:  "k,key",
				Usage: T("IDs of SSH key to add to the root user, multiple occurrence allowed"),
			},
			cli.BoolFlag{
				Name:  "b,upgrade-bios",
				Usage: T("Upgrade BIOS"),
			},
			cli.BoolFlag{
				Name:  "w,upgrade-firmware",
				Usage: T("Upgrade all hard drives' firmware"),
			},
			ForceFlag(),
		},
	}
}

func HardwareRescueMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_HARDWARE_NAME,
		Name:        CMD_HARDWARE_RESCUE_NAME,
		Description: T("Reboot server into a rescue image"),
		Usage:       "${COMMAND_NAME} sl hardware rescue IDENTIFIER [OPTIONS]",
		Flags: []cli.Flag{
			ForceFlag(),
		},
	}
}

func HardwareUpdateFirmwareMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_HARDWARE_NAME,
		Name:        CMD_HARDWARE_UPDATE_FIRMWARE_NAME,
		Description: T("Update server firmware"),
		Usage:       "${COMMAND_NAME} sl hardware update-firmware IDENTIFIER [OPTIONS]",
		Flags: []cli.Flag{
			ForceFlag(),
		},
	}
}

// HardwareToggleIPMIMetaData returns the metatadata of command `hardware toggle-ipmi`
func HardwareToggleIPMIMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_HARDWARE_NAME,
		Name:        "toggle-ipmi",
		Description: T("Toggle the IPMI interface on and off. This command is asynchronous."),
		Usage:       "${COMMAND_NAME} sl hardware toggle-ipmi IDENTIFIER [OPTIONS]",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "enable",
				Usage: T("Enable the IPMI interface."),
			},
			cli.BoolFlag{
				Name:  "disable",
				Usage: T("Disable the IPMI interface."),
			},
			QuietFlag(),
		},
	}
}
