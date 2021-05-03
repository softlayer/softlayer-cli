package metadata

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
)

var (
	NS_IPSEC_NAME  = "ipsec"
	CMD_IPSEC_NAME = "ipsec"

	CMD_IPSEC_CANCEL_NAME        = "cancel"
	CMD_IPSEC_CONFIG_NAME        = "config"
	CMD_IPSEC_DETAIL_NAME        = "detail"
	CMD_IPSEC_LIST_NAME          = "list"
	CMD_IPSEC_ORDER_NAME         = "order"
	CMD_IPSEC_SUBNET_ADD_NAME    = "subnet-add"
	CMD_IPSEC_SUBNET_REMOVE_NAME = "subnet-remove"
	CMD_IPSEC_TRANS_ADD_NAME     = "translation-add"
	CMD_IPSEC_TRANS_REMOVE_NAME  = "translation-remove"
	CMD_IPSEC_TRANS_UPDATE_NAME  = "translation-update"
	CMD_IPSEC_UPDATE_NAME        = "update"
)

func IpsecNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  NS_SL_NAME,
		Name:        NS_IPSEC_NAME,
		Description: T("Classic infrastructure IPSEC VPN"),
	}
}

func IpsecMetaData() cli.Command {
	return cli.Command{
		Category:    NS_SL_NAME,
		Name:        CMD_IPSEC_NAME,
		Description: T("Classic infrastructure IPSEC VPN"),
		Usage:       "${COMMAND_NAME} sl ipsec",
		Subcommands: []cli.Command{
			IpsecCancelMetaData(),
			IpsecConfigMetaData(),
			IpsecOrderMetaData(),
			IpsecDetailMetaData(),
			IpsecListMetaData(),
			IpsecSubnetAddMetaData(),
			IpsecSubnetRemoveMetaData(),
			IpsecTransAddMetaData(),
			IpsecTransRemoveMetaData(),
			IpsecTransUpdataMetaData(),
			IpsecUpdateMetaData(),
		},
	}
}

func IpsecCancelMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_IPSEC_NAME,
		Name:        CMD_IPSEC_CANCEL_NAME,
		Description: T("Cancel a IPSec VPN tunnel context"),
		Usage:       T(`${COMMAND_NAME} sl ipsec cancel CONTEXT_ID [OPTIONS]`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "immediate",
				Usage: T("Cancel the IPSec immediately instead of on the billing anniversary"),
			},
			cli.StringFlag{
				Name:  "reason",
				Usage: T("An optional reason for cancellation"),
			},
			ForceFlag(),
		},
	}
}

func IpsecConfigMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_IPSEC_NAME,
		Name:        CMD_IPSEC_CONFIG_NAME,
		Description: T("Request configuration of a tunnel context"),
		Usage: T(`${COMMAND_NAME} sl ipsec config CONTEXT_ID [OPTIONS]

  Request configuration of a tunnel context.

  This action will update the advancedConfigurationFlag on the context
  instance and further modifications against the context will be prevented
  until all changes can be propagated to network devices.`),
		Flags: []cli.Flag{},
	}
}

func IpsecOrderMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_IPSEC_NAME,
		Name:        CMD_IPSEC_ORDER_NAME,
		Description: T("Order a IPSec VPN tunnel"),
		Usage:       T(`${COMMAND_NAME} sl ipsec order [OPTIONS]`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Short name of the datacenter for the IPSec. For example, dal09[required]"),
			},
			OutputFlag(),
		},
	}
}

func IpsecDetailMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_IPSEC_NAME,
		Name:        CMD_IPSEC_DETAIL_NAME,
		Description: T("List IPSec VPN tunnel context details"),
		Usage: T(`${COMMAND_NAME} sl ipsec detail CONTEXT_ID [OPTIONS]

  List IPSEC VPN tunnel context details.

  Additional resources can be joined using multiple instances of the include
  option, for which the following choices are available.

  at: address translations
  is: internal subnets
  rs: remote subnets
  sr: statically routed subnets
  ss: service subnets`),
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:  "i,include",
				Usage: T("Include extra resources. Options are: at,is,rs,sr,ss"),
			},
			OutputFlag(),
		},
	}
}

func IpsecListMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_IPSEC_NAME,
		Name:        CMD_IPSEC_LIST_NAME,
		Description: T("List IPSec VPN tunnel contexts"),
		Usage:       "${COMMAND_NAME} sl ipsec list [OPTIONS]",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "order",
				Usage: T("Filter by ID of the order that purchased the IPSec"),
			},
			OutputFlag(),
		},
	}
}

func IpsecSubnetAddMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_IPSEC_NAME,
		Name:        CMD_IPSEC_SUBNET_ADD_NAME,
		Description: T("Add a subnet to an IPSec tunnel context"),
		Usage: T(`${COMMAND_NAME} sl ipsec subnet-add CONTEXT_ID [OPTIONS] 

  Add a subnet to an IPSEC tunnel context.

  A subnet id may be specified to link to the existing tunnel context.

  Otherwise, a network identifier in CIDR notation should be specified,
  indicating that a subnet resource should first be created before
  associating it with the tunnel context. Note that this is only supported
  for remote subnets, which are also deleted upon failure to attach to a
  context.

  A separate configuration request should be made to realize changes on
  network devices.`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "s,subnet-id",
				Usage: T("Subnet identifier to add, required"),
			},
			cli.StringFlag{
				Name:  "t,subnet-type",
				Usage: T("Subnet type to add. Options are: internal,remote,service[required]"),
			},
			cli.StringFlag{
				Name:  "n,network",
				Usage: T("Subnet network identifier to create"),
			},
		},
	}
}

func IpsecSubnetRemoveMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_IPSEC_NAME,
		Name:        CMD_IPSEC_SUBNET_REMOVE_NAME,
		Description: T("Remove a subnet from an IPSEC tunnel context"),
		Usage: T(`${COMMAND_NAME} sl ipsec subnet-remove CONTEXT_ID SUBNET_ID SUBNET_TYPE 

  Remove a subnet from an IPSEC tunnel context.

  The subnet id to remove must be specified.

  Remote subnets are deleted upon removal from a tunnel context.

  A separate configuration request should be made to realize changes on
  network devices.`),
		Flags: []cli.Flag{},
	}
}

func IpsecTransAddMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_IPSEC_NAME,
		Name:        CMD_IPSEC_TRANS_ADD_NAME,
		Description: T("Add an address translation to an IPSec tunnel"),
		Usage: T(`${COMMAND_NAME} sl ipsec translation-add CONTEXT_ID [OPTIONS]

  Add an address translation to an IPSEC tunnel context.

  A separate configuration request should be made to realize changes on
  network devices.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "s,static-ip",
				Usage: T("Static IP address[required]"),
			},
			cli.StringFlag{
				Name:  "r,remote-ip",
				Usage: T("Remote IP address[required]"),
			},
			cli.StringFlag{
				Name:  "n,note",
				Usage: T("Note value"),
			},
			OutputFlag(),
		},
	}
}

func IpsecTransRemoveMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_IPSEC_NAME,
		Name:        CMD_IPSEC_TRANS_REMOVE_NAME,
		Description: T("Remove a translation entry from an IPSec"),
		Usage: T(`${COMMAND_NAME} sl ipsec translation-remove CONTEXT_ID TRANSLATION_ID 

  Remove a translation entry from an IPSEC tunnel context.

  A separate configuration request should be made to realize changes on
  network devices.`),
		Flags: []cli.Flag{},
	}
}

func IpsecTransUpdataMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_IPSEC_NAME,
		Name:        CMD_IPSEC_TRANS_UPDATE_NAME,
		Description: T("Update an address translation for an IPSec"),
		Usage: T(`${COMMAND_NAME} sl ipsec translation-update CONTEXT_ID TRANSLATION_ID [OPTIONS]

  Update an address translation for an IPSEC tunnel context.

  A separate configuration request should be made to realize changes on
  network devices.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "s,static-ip",
				Usage: T("Static IP address[required]"),
			},
			cli.StringFlag{
				Name:  "r,remote-ip",
				Usage: T("Remote IP address[required]"),
			},
			cli.StringFlag{
				Name:  "n,note",
				Usage: T("Note"),
			},
			OutputFlag(),
		},
	}
}

func IpsecUpdateMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_IPSEC_NAME,
		Name:        CMD_IPSEC_UPDATE_NAME,
		Description: T("Update tunnel context properties"),
		Usage: T(`${COMMAND_NAME} sl ipsec update CONTEXT_ID [OPTIONS]

  Update tunnel context properties.

  Updates are made atomically, so either all are accepted or none are.

  Key life values must be in the range 120-172800.

  Phase 2 perfect forward secrecy must be in the range 0-1.

  A separate configuration request should be made to realize changes on
  network devices.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "n,name",
				Usage: T("Friendly name"),
			},
			cli.StringFlag{
				Name:  "r,remote-peer",
				Usage: T("Remote peer IP address"),
			},
			cli.StringFlag{
				Name:  "k,preshared-key",
				Usage: T("Preshared key"),
			},
			cli.StringFlag{
				Name:  "a,phase1-auth",
				Usage: T("Phase 1 authentication. Options are: MD5,SHA1,SHA256"),
			},
			cli.StringFlag{
				Name:  "c,phase1-crypto",
				Usage: T("Phase 1 encryption. Options are: DES,3DES,AES128,AES192,AES256"),
			},
			cli.IntFlag{
				Name:  "d,phase1-dh",
				Usage: T("Phase 1 Diffie-Hellman group. Options are: 0,1,2,5"),
			},
			cli.IntFlag{
				Name:  "t,phase1-key-ttl",
				Usage: T("Phase 1 key life. Range is 120-172800"),
			},
			cli.StringFlag{
				Name:  "u,phase2-auth",
				Usage: T("Phase 2 authentication. Options are: MD5,SHA1,SHA256"),
			},
			cli.StringFlag{
				Name:  "y,phase2-crypto",
				Usage: T("Phase 2 encryption. Options are: DES,3DES,AES128,AES192,AES256"),
			},
			cli.IntFlag{
				Name:  "e,phase2-dh",
				Usage: T("Phase 2 Diffie-Hellman group. Options are: 0,1,2,5"),
			},
			cli.IntFlag{
				Name:  "f,phase2-forward-secrecy",
				Usage: T("Phase 2 perfect forward secrecy. Range is 0-1"),
			},
			cli.IntFlag{
				Name:  "l,phase2-key-ttl",
				Usage: T("Phase 2 key life. Range is 120-172800"),
			},
			OutputFlag(),
		},
	}
}
