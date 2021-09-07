package metadata

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

var (
	NS_VIRTUAL_NAME  = "vs"
	CMD_VIRTUAL_NAME = "vs"

	CMD_VS_AUTHORIZE_STORAGE_NAME = "authorize-storage"
	CMD_VS_CANCEL_NAME            = "cancel"
	CMD_VS_CAPTURE_NAME           = "capture"
	CMD_VS_CREATE_NAME            = "create"
	CMD_VS_CREATE_HOST_NAME       = "host-create"
	CMD_VS_CREATE_OPTIONS_NAME    = "options"
	CMD_VS_CREDENTIALS_NAME       = "credentials"
	CMD_VS_DETAIL_NAME            = "detail"
	CMD_VS_DNS_SYNC_NAME          = "dns-sync"
	CMD_VS_EDIT_NAME              = "edit"
	CMD_VS_LIST_NAME              = "list"
	CMD_VS_LIST_HOST_NAME         = "host-list"
	CMD_VS_PAUSE_NAME             = "pause"
	CMD_VS_POWER_OFF_NAME         = "power-off"
	CMD_VS_POWER_ON_NAME          = "power-on"
	CMD_VS_READY_NAME             = "ready"
	CMD_VS_REBOOT_NAME            = "reboot"
	CMD_VS_RELOAD_NAME            = "reload"
	CMD_VS_RESCUE_NAME            = "rescue"
	CMD_VS_RESUME_NAME            = "resume"
	CMD_VS_UPGRADE_NAME           = "upgrade"
	CMD_VS_MIGRATE_NAME           = "migrate"
)

func VSNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  NS_SL_NAME,
		Name:        NS_VIRTUAL_NAME,
		Description: T("Classic infrastructure Virtual Servers"),
	}
}

func VSMetaData() cli.Command {
	return cli.Command{
		Category:    NS_SL_NAME,
		Name:        CMD_VIRTUAL_NAME,
		Description: T("Classic infrastructure Virtual Servers"),
		Usage:       "${COMMAND_NAME} sl vs",
		Subcommands: []cli.Command{
			VSCancelMetaData(),
			VSCaptureMetaData(),
			VSCreateHostMetaData(),
			VSCreateMetaData(),
			VSCreateOptionsMetaData(),
			VSCredentialsMetaData(),
			VSDetailMetaData(),
			VSDNSSyncMetaData(),
			VSEditMetaData(),
			VSListHostMetaData(),
			VSListMetaData(),
			VSMigrateMetaData(),
			VSPauseMetaData(),
			VSPowerOffMetaData(),
			VSPowerOnMetaData(),
			VSReadyMetaData(),
			VSRebootMetaData(),
			VSReloadMetaData(),
			VSRescueMetaData(),
			VSResumeMetaData(),
			VSUpgradeMetaData(),
			VSAuthorizeStorageMetaData(),
			VSBandwidthMetaData(),
		},
	}
}

func VSAuthorizeStorageMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VIRTUAL_NAME,
		Name:        CMD_VS_AUTHORIZE_STORAGE_NAME,
		Description: T("Authorize File, Block and Portable Storage to a Virtual Server"),
		Usage: T(`${COMMAND_NAME} sl vs authorize-storage [OPTIONS] IDENTIFIER

EXAMPLE:
   ${COMMAND_NAME} sl vs authorize-storage --username-storage SL01SL30-37 1234567
   Authorize File, Block and Portable Storage to a Virtual Server.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "u, username-storage",
				Usage: T("The storage username to be added to the virtual server."),
			},
			cli.IntFlag{
				Name:  "p, portable-id",
				Usage: T("The portable storage id to be added to the virtual server"),
			},
			OutputFlag(),
		},
	}
}

func VSCancelMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VIRTUAL_NAME,
		Name:        CMD_VS_CANCEL_NAME,
		Description: T("Cancel virtual server instance"),
		Usage: T(`${COMMAND_NAME} sl vs cancel IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vs cancel 12345678
   This command cancels virtual server instance with ID of 12345678.`),
		Flags: []cli.Flag{
			ForceFlag(),
		},
	}
}

func VSMigrateMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VIRTUAL_NAME,
		Name:        CMD_VS_MIGRATE_NAME,
		Description: T("Manage VSIs that require migration"),
		Usage: T(`${COMMAND_NAME} sl vs migrate [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vs migrate --guest 1234567
   Manage VSIs that require migration. Can migrate Dedicated Instance from one dedicated host to another dedicated host as well.`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "g, guest",
				Usage: T("Guest ID to immediately migrate."),
			},
			cli.BoolFlag{
				Name:  "a, all",
				Usage: T("Migrate ALL guests that require migration immediately."),
			},
			cli.IntFlag{
				Name:  "H, host",
				Usage: T("Dedicated Host ID to migrate to. Only works on guests that are already on a dedicated host."),
			},
			OutputFlag(),
		},
	}
}

func VSCaptureMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VIRTUAL_NAME,
		Name:        CMD_VS_CAPTURE_NAME,
		Description: T("Capture virtual server instance into an image"),
		Usage: T(`${COMMAND_NAME} sl vs capture IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vs capture 12345678 -n mycloud --all --note testing
   This command captures virtual server instance with ID of 12345678 with all disks into an image named "mycloud" with note "testing".`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "n,name",
				Usage: T("Name of the image [required]"),
			},
			cli.BoolFlag{
				Name:  "all",
				Usage: T("Capture all disks that belong to the virtual server"),
			},
			cli.StringFlag{
				Name:  "note",
				Usage: T("Add a note to be associated with the image"),
			},
			OutputFlag(),
		},
	}
}

func VSCreateHostMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VIRTUAL_NAME,
		Name:        CMD_VS_CREATE_HOST_NAME,
		Description: T("Create a host for dedicated virtual servers"),
		Usage:       "${COMMAND_NAME} sl vs host-create [OPTIONS]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "H,hostname",
				Usage: T("Host portion of the FQDN [required]"),
			},
			cli.StringFlag{
				Name:  "D,domain",
				Usage: T("Domain portion of the FQDN [required]"),
			},
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Datacenter shortname [required]"),
			},
			cli.StringFlag{
				Name:  "s,size",
				Usage: T("Size of the dedicated host, currently only one size is available: 56_CORES_X_242_RAM_X_1_4_TB"),
			},
			cli.StringFlag{
				Name:  "b,billing",
				Usage: T("Billing rate. Default is: hourly. Options are: hourly, monthly"),
			},
			cli.StringFlag{
				Name:  "v,vlan-private",
				Usage: T("The ID of the private VLAN on which you want the dedicated host placed. See: '${COMMAND_NAME} sl vlan list' for reference"),
			},
			ForceFlag(),
			OutputFlag(),
		},
	}
}

func VSCreateMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VIRTUAL_NAME,
		Name:        CMD_VS_CREATE_NAME,
		Description: T("Create virtual server instance"),
		Usage: T(`${COMMAND_NAME} sl vs create [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vs create -H myvsi -D ibm.com -c 4 -m 4096 -d dal10 -o UBUNTU_16_64 --disk 100 --disk 1000 --vlan-public 413
	This command orders a virtual server instance with hostname is myvsi, domain is ibm.com, 4 cpu cores, 4096M memory, located at datacenter: dal10,
	operation system is UBUNTU 16 64 bits, 2 disks, one is 100G, the other is 1000G, and placed at public vlan with ID 413.
	${COMMAND_NAME} sl vs create -H myvsi -D ibm.com -c 4 -m 4096 -d dal10 -o UBUNTU_16_64 --disk 100 --disk 1000 --vlan-public 413 --test
	This command tests whether the order is valid with above options before the order is actually placed.
	${COMMAND_NAME} sl vs create -H myvsi -D ibm.com -c 4 -m 4096 -d dal10 -o UBUNTU_16_64 --disk 100 --disk 1000 --vlan-public 413 --export ~/myvsi.txt
	This command exports above options to a file: myvsi.txt under user home directory for later use.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "H,hostname",
				Usage: T("Host portion of the FQDN [required]"),
			},
			cli.StringFlag{
				Name:  "D,domain",
				Usage: T("Domain portion of the FQDN [required]"),
			},
			cli.IntFlag{
				Name:  "c,cpu",
				Usage: T("Number of CPU cores [required]"),
			},
			cli.IntFlag{
				Name:  "m,memory",
				Usage: T("Memory in megabytes [required]"),
			},
			cli.StringFlag{
				Name:  "flavor",
				Usage: T("Public Virtual Server flavor key name"),
			},
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Datacenter shortname [required]"),
			},
			cli.StringFlag{
				Name:  "o,os",
				Usage: T("OS install code. Tip: you can specify <OS>_LATEST"),
			},
			cli.IntFlag{
				Name:  "image",
				Usage: T("Image ID. See: '${COMMAND_NAME} sl image list' for reference"),
			},
			cli.StringFlag{
				Name:  "billing",
				Usage: T("Billing rate. Default is: hourly. Options are: hourly, monthly"),
			},
			cli.BoolFlag{
				Name:  "dedicated",
				Usage: T("Create a dedicated Virtual Server (Private Node)"),
			},
			cli.IntFlag{
				Name:  "host-id",
				Usage: T("Host Id to provision a Dedicated Virtual Server onto"),
			},
			cli.BoolFlag{
				Name:  "san",
				Usage: T("Use SAN storage instead of local disk"),
			},
			cli.BoolFlag{
				Name:  "test",
				Usage: T("Do not actually create the virtual server"),
			},
			cli.StringFlag{
				Name:  "export",
				Usage: T("Exports options to a template file"),
			},
			cli.StringFlag{
				Name:  "i,postinstall",
				Usage: T("Post-install script to download"),
			},
			cli.IntSliceFlag{
				Name:  "k,key",
				Usage: T("The IDs of the SSH keys to add to the root user (multiple occurrence permitted)"),
			},
			cli.IntSliceFlag{
				Name:  "disk",
				Usage: T("Disk sizes (multiple occurrence permitted)"),
			},
			cli.BoolFlag{
				Name:  "private",
				Usage: T("Forces the virtual server to only have access the private network"),
			},
			cli.StringFlag{
				Name:  "like",
				Usage: T("Use the configuration from an existing virtual server"),
			},
			cli.IntFlag{
				Name:  "n,network",
				Usage: T("Network port speed in Mbps"),
			},
			cli.StringSliceFlag{
				Name:  "g,tag",
				Usage: T("Tags to add to the instance (multiple occurrence permitted)"),
			},
			cli.StringFlag{
				Name:  "t,template",
				Usage: T("A template file that defaults the command-line options"),
			},
			cli.StringFlag{
				Name:  "u,userdata",
				Usage: T("User defined metadata string"),
			},
			cli.StringFlag{
				Name:  "F,userfile",
				Usage: T("Read userdata from file"),
			},
			cli.StringFlag{
				Name:  "vlan-public",
				Usage: T("The ID of the public VLAN on which you want the virtual server placed"),
			},
			cli.StringFlag{
				Name:  "vlan-private",
				Usage: T("The ID of the private VLAN on which you want the virtual server placed"),
			},
			cli.IntSliceFlag{
				Name:  "S,public-security-group",
				Usage: T("Security group ID to associate with the public interface (multiple occurrence permitted)"),
			},
			cli.IntSliceFlag{
				Name:  "s,private-security-group",
				Usage: T("Security group ID to associate with the private interface (multiple occurrence permitted)"),
			},
			cli.IntFlag{
				Name:  "wait",
				Usage: T("Wait until the virtual server is finished provisioning for up to X seconds before returning. It's not compatible with option --quantity"),
			},
			cli.IntFlag{
				Name:  "placement-group-id",
				Usage: T("Placement Group Id to order this guest on."),
			},
			cli.StringFlag{
				Name:  "boot-mode",
				Usage: T("Specify the mode to boot the OS in. Supported modes are HVM and PV."),
			},
			cli.IntFlag{
				Name:  "subnet-public",
				Usage: T("The ID of the public SUBNET on which you want the virtual server placed"),
			},
			cli.IntFlag{
				Name:  "subnet-private",
				Usage: T("The ID of the private SUBNET on which you want the virtual server placed"),
			},
			cli.BoolFlag{
				Name:  "transient",
				Usage: T("Create a transient virtual server"),
			},
			cli.IntFlag{
				Name:  "quantity",
				Usage: T("The quantity of virtual server be created. It should be greater or equal to 1. This value defaults to 1."),
				Value: 1,
			},
			ForceFlag(),
		},
	}
}

func VSCreateOptionsMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VIRTUAL_NAME,
		Name:        CMD_VS_CREATE_OPTIONS_NAME,
		Description: T("List options for creating virtual server instance"),
		Usage: T(`${COMMAND_NAME} sl vs options [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vs options
   This command lists all the options for creating a virtual server instance, eg.datacenters, cpu, memory, os, disk, network speed, etc.`),
		Flags: []cli.Flag{
			OutputFlag(),
		},
	}
}

func VSCredentialsMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VIRTUAL_NAME,
		Name:        CMD_VS_CREDENTIALS_NAME,
		Description: T("List virtual server instance credentials"),
		Usage: T(`${COMMAND_NAME} sl vs credentials IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vs credentials 12345678
   This command lists all username and password pairs of virtual server instance with ID 12345678.`),
		Flags: []cli.Flag{
			OutputFlag(),
		},
	}
}

func VSDetailMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VIRTUAL_NAME,
		Name:        CMD_VS_DETAIL_NAME,
		Description: T("Get details for a virtual server instance"),
		Usage: T(`${COMMAND_NAME} sl vs detail IDENTIFIER [OPTIONS] 
	
EXAMPLE:
   ${COMMAND_NAME} sl vs details 12345678
   This command lists detailed information about virtual server instance with ID 12345678.`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "passwords",
				Usage: T("Show passwords (check over your shoulder!)"),
			},
			cli.BoolFlag{
				Name:  "price",
				Usage: T("Show associated prices"),
			},
			OutputFlag(),
		},
	}
}

func VSDNSSyncMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VIRTUAL_NAME,
		Name:        CMD_VS_DNS_SYNC_NAME,
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
			ForceFlag(),
		},
	}
}

func VSEditMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VIRTUAL_NAME,
		Name:        CMD_VS_EDIT_NAME,
		Description: T("Edit a virtual server instance's details"),
		Usage: T(`${COMMAND_NAME} sl vs edit IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vs edit 12345678 -D ibm.com -H myapp --tag testcli --public-speed 1000
   This command updates virtual server instance with ID 12345678 and set its domain to be "ibm.com", hostname to "myapp", tag to "testcli", 
   and public network port speed to 1000 Mbps.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "D,domain",
				Usage: T("Domain portion of the FQDN"),
			},
			cli.StringFlag{
				Name:  "H,hostname",
				Usage: T("Host portion of the FQDN. example: server"),
			},
			cli.StringSliceFlag{
				Name:  "g,tag",
				Usage: T("Tags to set or empty string to remove all"),
			},
			cli.StringFlag{
				Name:  "u,userdata",
				Usage: T("User defined metadata string"),
			},
			cli.StringFlag{
				Name:  "F,userfile",
				Usage: T("Read userdata from file"),
			},
			cli.IntFlag{
				Name:  "public-speed",
				Usage: T("Public port speed, options are: 0,10,100,1000,10000"),
			},
			cli.IntFlag{
				Name:  "private-speed",
				Usage: T("Private port speed, options are: 0,10,100,1000,10000"),
			},
		},
	}
}

func VSListHostMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VIRTUAL_NAME,
		Name:        CMD_VS_LIST_HOST_NAME,
		Description: T("List dedicated hosts on your account"),
		Usage:       "${COMMAND_NAME} sl vs host-list [OPTIONS]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "n,name",
				Usage: T("Filter by name of the dedicated host"),
			},
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Filter by datacenter of the dedicated host"),
			},
			cli.StringFlag{
				Name:  "owner",
				Usage: T("Filter by owner of the dedicated host"),
			},
			cli.IntFlag{
				Name:  "order",
				Usage: T("Filter by ID of the order which purchased this dedicated host"),
			},
			OutputFlag(),
		},
	}
}

func VSListMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VIRTUAL_NAME,
		Name:        CMD_VS_LIST_NAME,
		Description: T("List virtual server instances on your account"),
		Usage: T(`${COMMAND_NAME} sl vs list [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl vs list --domain ibm.com --hourly --sortby memory
   This command lists all hourly-billing virtual server instances on current account filtering domain equals to "ibm.com" and sort them by memory.`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "c,cpu",
				Usage: T("Filter by number of CPU cores"),
			},
			cli.StringFlag{
				Name:  "D,domain",
				Usage: T("Filter by domain portion of the FQDN"),
			},
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Filter by datacenter shortname"),
			},
			cli.StringFlag{
				Name:  "H,hostname",
				Usage: T("Filter by host portion of the FQDN"),
			},
			cli.IntFlag{
				Name:  "m,memory",
				Usage: T("Filter by memory in megabytes"),
			},
			cli.IntFlag{
				Name:  "n,network",
				Usage: T("Filter by network port speed in Mbps"),
			},
			cli.StringFlag{
				Name:  "P,public-ip",
				Usage: T("Filter by public IP address"),
			},
			cli.StringFlag{
				Name:  "p,private-ip",
				Usage: T("Filter by private IP address"),
			},
			cli.BoolFlag{
				Name:  "hourly",
				Usage: T("Show only hourly instances"),
			},
			cli.BoolFlag{
				Name:  "monthly",
				Usage: T("Show only monthly instances"),
			},
			cli.StringSliceFlag{
				Name:  "g,tag",
				Usage: T("Filter by tags (multiple occurrence permitted)"),
			},
			cli.IntFlag{
				Name:  "o,order",
				Usage: T("Filter by ID of the order which purchased this instance"),
			},
			cli.StringFlag{
				Name:  "owner",
				Usage: T("Filtered by Id of user who owns the instances"),
			},
			cli.StringFlag{
				Name:  "sortby",
				Usage: T("Column to sort by, default is:hostname, options are:id,hostname,domain,datacenter,cpu,memory,public_ip,private_ip"),
			},
			cli.StringSliceFlag{
				Name:  "column",
				Usage: T("Column to display. Options are: id,hostname,domain,cpu,memory,public_ip,private_ip,datacenter,action,guid,power_state,created_by,tags. This option can be specified multiple times"),
			},
			cli.StringSliceFlag{
				Name:   "columns",
				Hidden: true,
			},
			OutputFlag(),
		},
	}
}

func VSPauseMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VIRTUAL_NAME,
		Name:        CMD_VS_PAUSE_NAME,
		Description: T("Pause an active virtual server instance"),
		Usage: T(`${COMMAND_NAME} sl vs pause IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl vs pause 12345678 -f
   This command pauses virtual server instance with ID 12345678 without asking for confirmation.`),
		Flags: []cli.Flag{
			ForceFlag(),
		},
	}
}

func VSPowerOffMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VIRTUAL_NAME,
		Name:        CMD_VS_POWER_OFF_NAME,
		Description: T("Power off an active virtual server instance"),
		Usage: T(`${COMMAND_NAME} sl vs power-off IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl vs power-off 12345678 --soft
   This command performs a soft power off for virtual server instance with ID 12345678.`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "hard",
				Usage: T("Perform a hard shutdown"),
			},
			cli.BoolFlag{
				Name:  "soft",
				Usage: T("Perform a soft shutdown"),
			},
			ForceFlag(),
		},
	}
}

func VSPowerOnMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VIRTUAL_NAME,
		Name:        CMD_VS_POWER_ON_NAME,
		Description: T("Power on a virtual server instance"),
		Usage: T(`${COMMAND_NAME} sl vs power-on IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl vs power-on 12345678
   This command performs a power on for virtual server instance with ID 12345678.`),
		Flags: []cli.Flag{
			ForceFlag(),
		},
	}
}

func VSReadyMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VIRTUAL_NAME,
		Name:        CMD_VS_READY_NAME,
		Description: T("Check if a virtual server instance is ready for use"),
		Usage: T(`${COMMAND_NAME} sl vs ready IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vs ready 12345678 --wait 30
   This command checks virtual server instance with ID 12345678 status to see if it is ready for use continuously and waits up to 30 seconds.`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "wait",
				Usage: T("Wait until the virtual server is finished provisioning for up to X seconds before returning"),
			},
		},
	}
}

func VSRebootMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VIRTUAL_NAME,
		Name:        CMD_VS_REBOOT_NAME,
		Description: T("Reboot an active virtual server instance"),
		Usage: T(`${COMMAND_NAME} sl vs reboot IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl vs reboot 12345678 --hard
   This command performs a hard reboot for virtual server instance with ID 12345678.`),
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

func VSReloadMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VIRTUAL_NAME,
		Name:        CMD_VS_RELOAD_NAME,
		Description: T("Reload operating system on a virtual server instance"),
		Usage: T(`${COMMAND_NAME} sl vs reload IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl vs reload 12345678
   This command reloads current operating system for virtual server instance with ID 12345678.
   ${COMMAND_NAME} sl vs reload 12345678 --image 1234
   This command reloads operating system from image with ID 1234 for virtual server instance with ID 12345678.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "i,postinstall",
				Usage: T("Post-install script to download"),
			},
			cli.IntFlag{
				Name:  "image",
				Usage: T("Image ID. The default is to use the current operating system.\nSee: '${COMMAND_NAME} sl image list' for reference"),
			},
			cli.IntSliceFlag{
				Name:  "k,key",
				Usage: T("The IDs of the SSH keys to add to the root user (multiple occurrence permitted)"),
			},
			ForceFlag(),
		},
	}
}

func VSRescueMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VIRTUAL_NAME,
		Name:        CMD_VS_RESCUE_NAME,
		Description: T("Reboot a virtual server instance into a rescue image"),
		Usage: T(`${COMMAND_NAME} sl vs rescue IDENTIFIER [OPTIONS]
	
EXAMPLE:
   ${COMMAND_NAME} sl vs rescue 12345678
   This command reboots virtual server instance with ID 12345678 into a rescue image.`),
		Flags: []cli.Flag{
			ForceFlag(),
		},
	}
}

func VSResumeMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VIRTUAL_NAME,
		Name:        CMD_VS_RESUME_NAME,
		Description: T("Resume a paused virtual server instance"),
		Usage: T(`${COMMAND_NAME} sl vs resume IDENTIFIER [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl vs resume 12345678
   This command resumes virtual server instance with ID 12345678.`),
		Flags: []cli.Flag{
			ForceFlag(),
		},
	}
}

func VSUpgradeMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_VIRTUAL_NAME,
		Name:        CMD_VS_UPGRADE_NAME,
		Description: T("Upgrade a virtual server instance"),
		Usage: T(`${COMMAND_NAME} sl vs upgrade IDENTIFIER [OPTIONS]
	Note: Classic infrastructure service automatically reboots the instance once upgrade request is
  	placed. The instance is halted until the upgrade transaction is completed.
  	However for Network, no reboot is required.

EXAMPLE:
   ${COMMAND_NAME} sl vs upgrade 12345678 -c 8 -m 8192 --network 1000
   This commands upgrades virtual server instance with ID 12345678 and set number of CPU cores to 8, memory to 8192M, network port speed to 1000 Mbps.`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "c,cpu",
				Usage: T("Number of CPU cores"),
			},
			cli.BoolFlag{
				Name:  "private",
				Usage: T("CPU core will be on a dedicated host server"),
			},
			cli.IntFlag{
				Name:  "m,memory",
				Usage: T("Memory in megabytes"),
			},
			cli.IntFlag{
				Name:  "network",
				Usage: T("Network port speed in Mbps"),
			},
			cli.StringFlag{
				Name:  "flavor",
				Usage: T("Flavor key name"),
			},
			ForceFlag(),
			OutputFlag(),
		},
	}
}

func VSBandwidthMetaData() cli.Command {
    return cli.Command{
        Category:    CMD_VIRTUAL_NAME,
        Name:        "bandwidth",
        Description: T("Bandwidth data over date range."),
        Usage: T(`${COMMAND_NAME} sl {{.Command}} bandwidth upgrade IDENTIFIER [OPTIONS]
Time formats that are either '2006-01-02', '2006-01-02T15:04' or '2006-01-02T15:04-07:00'

Due to some rounding and date alignment details, results here might be slightly different than results in the control portal.
Bandwidth is listed in GB, if no time zone is specified, GMT+0 is assumed.

Example::

   ${COMMAND_NAME} sl {{.Command}} bandwidth 1234 -s 2006-01-02T15:04 -e 2006-01-02T15:04-07:00`, map[string]interface{}{"Command": "vs"}),
        Flags: []cli.Flag{
            cli.StringFlag{
                Name:  "s,start",
                Usage: T("Start date for bandwdith reporting"),
            },
            cli.StringFlag{
                Name:  "e,end",
                Usage: T("End date for bandwidth reporting"),
            },
            cli.IntFlag{
                Name:  "r,rollup",
                Usage: T("Number of seconds to report as one data point. 300, 600, 1800, 3600 (default), 43200 or 86400 seconds"),
            },
            cli.BoolFlag{
                Name:  "q,quite",
                Usage: T("Only show the summary table."),
            },
            OutputFlag(),
        },
    }
}