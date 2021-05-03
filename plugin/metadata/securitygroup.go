package metadata

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
	. "github.ibm.com/cgallo/softlayer-cli/plugin/i18n"
)

var (
	NS_SECURITYGROUP_NAME  = "securitygroup"
	CMD_SECURITYGROUP_NAME = "securitygroup"

	CMD_SECURITYGROUP_CREATE_NAME           = "create"
	CMD_SECURITYGROUP_DELETE_NAME           = "delete"
	CMD_SECURITYGROUP_DETAIL_NAME           = "detail"
	CMD_SECURITYGROUP_EDIT_NAME             = "edit"
	CMD_SECURITYGROUP_INTERFACE_ADD_NAME    = "interface-add"
	CMD_SECURITYGROUP_INTERFACE_LIST_NAME   = "interface-list"
	CMD_SECURITYGROUP_INTERFACE_REMOVE_NAME = "interface-remove"
	CMD_SECURITYGROUP_LIST_NAME             = "list"
	CMD_SECURITYGROUP_RULE_ADD_NAME         = "rule-add"
	CMD_SECURITYGROUP_RULE_EDIT_NAME        = "rule-edit"
	CMD_SECURITYGROUP_RULE_LIST_NAME        = "rule-list"
	CMD_SECURITYGROUP_RULE_REMOVE_NAME      = "rule-remove"
)

func SecurityGroupNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  NS_SL_NAME,
		Name:        NS_SECURITYGROUP_NAME,
		Description: T("Classic infrastructure network security groups"),
	}
}

func SecurityGroupMetaData() cli.Command {
	return cli.Command{
		Category:    NS_SL_NAME,
		Name:        CMD_SECURITYGROUP_NAME,
		Description: T("Classic infrastructure network security groups"),
		Usage:       "${COMMAND_NAME} sl securitygroup",
		Subcommands: []cli.Command{
			SecurityGroupCreateMetaData(),
			SecurityGroupDeleteMetaData(),
			SecurityGroupDetailMetaData(),
			SecurityGroupEditMetaData(),
			SecurityGroupInterfaceAddMetaData(),
			SecurityGroupInterfaceListMetaData(),
			SecurityGroupInterfaceRemoveMetaData(),
			SecurityGroupListMetaData(),
			SecurityGroupRuleAddMetaData(),
			SecurityGroupRuleEditMetaData(),
			SecurityGroupRuleListMetaData(),
			SecurityGroupRuleRemoveMetaData(),
		},
	}
}

func SecurityGroupCreateMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_SECURITYGROUP_NAME,
		Name:        CMD_SECURITYGROUP_CREATE_NAME,
		Description: T("Create a security group"),
		Usage:       "${COMMAND_NAME} sl securitygroup create [OPTIONS]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "n,name",
				Usage: T("The name of the security group"),
			},
			cli.StringFlag{
				Name:  "d,description",
				Usage: T("The description of the security group"),
			},
			OutputFlag(),
		},
	}
}

func SecurityGroupDeleteMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_SECURITYGROUP_NAME,
		Name:        CMD_SECURITYGROUP_DELETE_NAME,
		Description: T("Delete the given security group"),
		Usage:       "${COMMAND_NAME} sl securitygroup delete SECURITYGROUP_ID [OPTIONS]",
		Flags: []cli.Flag{
			ForceFlag(),
		},
	}
}

func SecurityGroupDetailMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_SECURITYGROUP_NAME,
		Name:        CMD_SECURITYGROUP_DETAIL_NAME,
		Description: T("Get details about a security group"),
		Usage:       "${COMMAND_NAME} sl securitygroup detail SECURITYGROUP_ID [OPTIONS]",
		Flags: []cli.Flag{
			OutputFlag(),
		},
	}
}

func SecurityGroupEditMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_SECURITYGROUP_NAME,
		Name:        CMD_SECURITYGROUP_EDIT_NAME,
		Description: T("Edit details of a security group"),
		Usage:       "${COMMAND_NAME} sl securitygroup edit SECURITYGROUP_ID [OPTIONS]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "n,name",
				Usage: T("The name of the security group"),
			},
			cli.StringFlag{
				Name:  "d,description",
				Usage: T("The description of the security group"),
			},
		},
	}
}

func SecurityGroupInterfaceAddMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_SECURITYGROUP_NAME,
		Name:        CMD_SECURITYGROUP_INTERFACE_ADD_NAME,
		Description: T("Attach an interface to a security group"),
		Usage:       "${COMMAND_NAME} sl securitygroup interface-add SECURITYGROUP_ID [OPTIONS]",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "n,network-component",
				Usage: T("The network component ID to associate with the security group"),
			},
			cli.StringFlag{
				Name:  "s,server",
				Usage: T(" The server ID to associate with the security group"),
			},
			cli.StringFlag{
				Name:  "i,interface",
				Usage: T("The interface of the server to associate (public/private)"),
			},
		},
	}
}

func SecurityGroupInterfaceListMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_SECURITYGROUP_NAME,
		Name:        CMD_SECURITYGROUP_INTERFACE_LIST_NAME,
		Description: T("List interfaces associated with security group"),
		Usage:       "${COMMAND_NAME} sl securitygroup interface-list SECURITYGROUP_ID [OPTIONS]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "sortby",
				Usage: T("Column to sort by. Options are: id,virtualServerId,hostname"),
			},
			OutputFlag(),
		},
	}
}

func SecurityGroupInterfaceRemoveMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_SECURITYGROUP_NAME,
		Name:        CMD_SECURITYGROUP_INTERFACE_REMOVE_NAME,
		Description: T("Detach an interface from a security group"),
		Usage:       "${COMMAND_NAME} sl securitygroup interface-remove SECURITYGROUP_ID [OPTIONS]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "n,network-component",
				Usage: T("The network component to remove from the security group"),
			},
			cli.StringFlag{
				Name:  "s,server",
				Usage: T(" The server ID to remove from the security group"),
			},
			cli.StringFlag{
				Name:  "i,interface",
				Usage: T("The interface of the server to remove (public or private)"),
			},
		},
	}
}

func SecurityGroupListMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_SECURITYGROUP_NAME,
		Name:        CMD_SECURITYGROUP_LIST_NAME,
		Description: T("List security groups"),
		Usage:       "${COMMAND_NAME} sl securitygroup list [OPTIONS]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "sortby",
				Usage: T("Column to sort by. Options are: id,name,description,created"),
			},
			OutputFlag(),
		},
	}
}

func SecurityGroupRuleAddMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_SECURITYGROUP_NAME,
		Name:        CMD_SECURITYGROUP_RULE_ADD_NAME,
		Description: T("Add a security group rule to a security group"),
		Usage:       "${COMMAND_NAME} sl securitygroup rule-add SECURITYGROUP_ID [OPTIONS]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "r,remote-ip",
				Usage: T("The remote IP/CIDR to enforce"),
			},
			cli.IntFlag{
				Name:  "s,remote-group",
				Usage: T("The ID of the remote security group to enforce"),
			},
			cli.StringFlag{
				Name:  "d,direction",
				Usage: T("The direction of traffic to enforce (ingress or egress), required"),
			},
			cli.StringFlag{
				Name:  "e,ether-type",
				Usage: T("The ethertype (IPv4 or IPv6) to enforce, default is IPv4 if not specified"),
			},
			cli.IntFlag{
				Name:  "M,port-max",
				Usage: T("The upper port bound to enforce"),
			},
			cli.IntFlag{
				Name:  "m,port-min",
				Usage: T("The lower port bound to enforce"),
			},
			cli.StringFlag{
				Name:  "p,protocol",
				Usage: T("The protocol (icmp, tcp, udp) to enforce"),
			},
			OutputFlag(),
		},
	}
}

func SecurityGroupRuleEditMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_SECURITYGROUP_NAME,
		Name:        CMD_SECURITYGROUP_RULE_EDIT_NAME,
		Description: T("Edit a security group rule in a security group"),
		Usage:       "${COMMAND_NAME} sl securitygroup rule-edit SECURITYGROUP_ID RULE_ID [OPTIONS]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "r,remote-ip",
				Usage: T("The remote IP/CIDR to enforce"),
			},
			cli.IntFlag{
				Name:  "s,remote-group",
				Usage: T("The ID of the remote security group to enforce"),
			},
			cli.StringFlag{
				Name:  "d,direction",
				Usage: T("The direction of traffic to enforce (ingress or egress), required"),
			},
			cli.StringFlag{
				Name:  "e,ether-type",
				Usage: T("The ethertype (IPv4 or IPv6) to enforce, default is IPv4 if not specified"),
			},
			cli.IntFlag{
				Name:  "M,port-max",
				Usage: T("The upper port bound to enforce"),
			},
			cli.IntFlag{
				Name:  "m,port-min",
				Usage: T("The lower port bound to enforce"),
			},
			cli.StringFlag{
				Name:  "p,protocol",
				Usage: T("The protocol (icmp, tcp, udp) to enforce"),
			},
		},
	}
}

func SecurityGroupRuleListMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_SECURITYGROUP_NAME,
		Name:        CMD_SECURITYGROUP_RULE_LIST_NAME,
		Description: T("List security group rules"),
		Usage:       "${COMMAND_NAME} sl securitygroup rule-list SECURITYGROUP_ID [OPTIONS]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "sortby",
				Usage: T("Column to sort by. Options are: id,remoteIp,remoteGroupId,direction,ethertype,portRangeMin,portRangeMax,protocol"),
			},
			OutputFlag(),
		},
	}
}

func SecurityGroupRuleRemoveMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_SECURITYGROUP_NAME,
		Name:        CMD_SECURITYGROUP_RULE_REMOVE_NAME,
		Description: T("Remove a rule from a security group"),
		Usage:       "${COMMAND_NAME} sl securitygroup rule-remove SECURITYGROUP_ID RULE_ID [OPTIONS]",
		Flags: []cli.Flag{
			ForceFlag(),
		},
	}
}
