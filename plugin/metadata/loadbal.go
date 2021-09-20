package metadata

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

var (
	NS_LOADBAL_NAME                = "loadbal"
	CMD_LOADBAL_CANCLE_NAME        = "cancel"
	CMD_LOADBAL_ORDER_NAME         = "order"
	CMD_LOADBAL_ORDER_OPTIONS_NAME = "order-options"
	CMD_LOADBAL_DETAIL_NAME        = "detail"
	CMD_LOADBAL_LIST_NAME          = "list"
	CMD_LOADBAL_HEALTH_NAME        = "health-edit"

	//protocol list can be found in the loadbal detail
	CMD_LOADBAL_PROTOCOL_ADD_NAME    = "protocol-add"
	CMD_LOADBAL_PROTOCOL_DELETE_NAME = "protocol-delete"

	// membe list can be found in the LB list
	CMD_LOADBAL_MEMBER_ADD_NAME = "member-add"
	CMD_LOADBAL_MEMBER_DEL_NAME = "member-delete"

	//l7 pool list can be found in the lb details
	CMD_LOADBAL_L7POOL_ADD_NAME    = "l7pool-add" //add pool into LB
	CMD_LOADBAL_L7POOL_DETAIL_NAME = "l7pool-detail"
	CMD_LOADBAL_L7POOL_DELETE_NAME = "l7pool-delete"
	CMD_LOADBAL_L7POOL_EDIT_NAME   = "l7pool-edit"

	//l7 member list can be found in the L7 pool details
	CMD_LOADBAL_L7MEMBER_ADD_NAME    = "l7member-add"
	CMD_LOADBAL_L7MEMBER_DELETE_NAME = "l7member-delete"

	//policies list can be found in protocol details
	CMD_LOADBAL_L7POLICY_ADD_NAME    = "l7policy-add" //add policy to a protocol
	CMD_LOADBAL_L7POLICY_DELETE_NAME = "l7policy-delete"
	CMD_LOADBAL_L7POLICY_LIST_NAME   = "l7policies"

	//rule list can be found in the L7 policy details
	CMD_LOADBAL_L7RULE_ADD_NAME    = "l7rule-add"
	CMD_LOADBAL_L7RULE_DELETE_NAME = "l7rule-delete"
	CMD_LOADBAL_L7RULE_LIST_NAME   = "l7rules"
)

func LoadbalNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  NS_SL_NAME,
		Name:        NS_LOADBAL_NAME,
		Description: T("Classic infrastructure Load Balancers"),
	}
}

func LoadbalMetaData() cli.Command {
	return cli.Command{
		Category:    NS_SL_NAME,
		Name:        NS_LOADBAL_NAME,
		Description: T("Classic infrastructure Load Balancers"),
		Usage:       "${COMMAND_NAME} sl loadbal",
		Subcommands: []cli.Command{
			LoadbalCancelMetadata(),
			LoadbalOrderMetadata(),
			LoadbalOrderOptionsMetadata(),
			LoadbalDetailMetadata(),
			LoadbalListMetadata(),
			LoadbalHealthMetadata(),
			LoadbalMemberAddMetadata(),
			LoadbalMemberDelMetadata(),
			LoadbalProtocolAddMetadata(),
			LoadbalProtocolDelMetadata(),
			LoadbalProtocolEditMetadata(),
			LoadbalL7PoolDelMetadata(),
			LoadbalL7PoolAddMetadata(),
			LoadbalL7PoolDetailMetadata(),
			LoadbalL7PoolEditMetadata(),
			LoadbalL7MemberAddMetadata(),
			LoadbalL7MemberDeleteMetadata(),
			LoadbalL7PolicyAddMetadata(),
			LoadbalL7PolicyDeleteMetadata(),
			LoadbalL7PolicyListMetadata(),
			LoadbalL7RuleAddMetadata(),
			LoadbalL7RuleDelMetadata(),
			LoadbalL7RuleListMetadata(),
		},
	}
}

func LoadbalCancelMetadata() cli.Command {
	return cli.Command{
		Category:    NS_LOADBAL_NAME,
		Name:        CMD_LOADBAL_CANCLE_NAME,
		Description: T("Cancel an existing load balancer"),
		Usage:       "${COMMAND_NAME} sl loadbal cancel (--id LOADBAL_ID)",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "id",
				Usage: T("ID for the load balancer [required]"),
			},
			ForceFlag(),
		},
	}
}

func LoadbalOrderMetadata() cli.Command {
	return cli.Command{
		Category:    NS_LOADBAL_NAME,
		Name:        CMD_LOADBAL_ORDER_NAME,
		Description: T("Order a load balancer"),
		Usage:       "${COMMAND_NAME} sl loadbal order (-n, --name NAME) (-d, --datacenter DATACENTER) (-t, --type PublicToPrivate | PrivateToPrivate | PublicToPublic ) [-l, --label LABEL] [ -s, --subnet SUBNET_ID] [--frontend-protocol PROTOCOL] [--frontend-port PORT] [--backend-protocol PROTOCOL] [--backend-port PORT] [-m, --method METHOD] [-c, --connections CONNECTIONS] [--sticky cookie | source-ip] [--use-public-subnet] [--verify]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "n,name",
				Usage: T("Name for this load balancer [required]"),
			},
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Datacenter name. It can be found from the keyName in the command '${COMMAND_NAME} sl order package-locations LBAAS' output. [required]"),
			},
			cli.StringFlag{
				Name:  "t,type",
				Usage: T("Load balancer type: PublicToPrivate | PrivateToPrivate | PublicToPublic [required]"),
			},
			cli.IntFlag{
				Name:  "s,subnet",
				Usage: T("Private subnet Id to order the load balancer. See '${COMMAND_NAME} sl loadbal order-options'. Only available in PublicToPrivate and PrivateToPrivate load balancer type"),
			},
			cli.StringFlag{
				Name:  "l,label",
				Usage: T("A descriptive label for this load balancer"),
			},
			cli.StringFlag{
				Name:  "frontend-protocol",
				Usage: T("Frontend protocol [default: HTTP]"),
			},
			cli.IntFlag{
				Name:  "frontend-port",
				Usage: T("Frontend port [default: 80]"),
			},
			cli.StringFlag{
				Name:  "backend-protocol",
				Usage: T("Backend protocol [default: HTTP]"),
			},
			cli.IntFlag{
				Name:  "backend-port",
				Usage: T("Backend port [default: 80]"),
			},
			cli.StringFlag{
				Name:  "m,method",
				Usage: T("Balancing Method: ROUNDROBIN | LEASTCONNECTION | WEIGHTED_RR, default: ROUNDROBIN"),
			},
			cli.IntFlag{
				Name:  "c, connections",
				Usage: T("Maximum number of connections to allow"),
			},
			cli.StringFlag{
				Name:  "sticky",
				Usage: T("Use 'cookie' or 'source-ip' to stick"),
			},
			cli.BoolFlag{
				Name:  "use-public-subnet",
				Usage: T("If this option is specified, the public ip will be allocated from a public subnet in this account. Otherwise, it will be allocated form IBM system pool. Only available in PublicToPrivate load balancer type."),
			},
			cli.BoolFlag{
				Name:  "verify",
				Usage: T("Only verify an order, dont actually create one"),
			},
			ForceFlag(),
		},
	}
}

func LoadbalOrderOptionsMetadata() cli.Command {
	return cli.Command{
		Category:    NS_LOADBAL_NAME,
		Name:        CMD_LOADBAL_ORDER_OPTIONS_NAME,
		Description: T("List options for order a load balancer"),
		Usage:       "${COMMAND_NAME} sl loadbal order-options [-d, --datacenter DATACENTER]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Show only selected datacenter, use shortname (dal13) format"),
			},
		},
	}
}

func LoadbalDetailMetadata() cli.Command {
	return cli.Command{
		Category:    NS_LOADBAL_NAME,
		Name:        CMD_LOADBAL_DETAIL_NAME,
		Description: T("Get load balancer details"),
		Usage:       "${COMMAND_NAME} sl loadbal detail (--id LOADBAL_ID)",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "id",
				Usage: T("ID for the load balancer [required]"),
			},
		},
	}
}

func LoadbalListMetadata() cli.Command {
	return cli.Command{
		Category:    NS_LOADBAL_NAME,
		Name:        CMD_LOADBAL_LIST_NAME,
		Description: T("List active load balancers"),
		Usage:       "${COMMAND_NAME} sl loadbal list",
		Flags:       []cli.Flag{},
	}
}

func LoadbalHealthMetadata() cli.Command {
	return cli.Command{
		Category:    NS_LOADBAL_NAME,
		Name:        CMD_LOADBAL_HEALTH_NAME,
		Description: T("Edit load balancer health check"),
		Usage:       "${COMMAND_NAME} sl loadbal health-edit (--lb-id LOADBAL_ID)  (--health-uuid HEALTH_CHECK_UUID) [-i, --interval INTERVAL] [-r, --retry RETRY] [-t, --timeout TIMEOUT] [-u, --url URL]",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "lb-id",
				Usage: T("ID for the load balancer [required]"),
			},
			cli.StringFlag{
				Name:  "health-uuid",
				Usage: T("Health check UUID to modify [required]"),
			},
			cli.IntFlag{
				Name:  "i,interval",
				Usage: T("Seconds between checks. [2-60]"),
			},
			cli.IntFlag{
				Name:  "r,retry",
				Usage: T("Number of times before marking as DOWN. [1-10]"),
			},
			cli.IntFlag{
				Name:  "t,timeout",
				Usage: T("Seconds to wait for a connection. [1-59]"),
			},
			cli.StringFlag{
				Name:  "u,url",
				Usage: T("Url path for HTTP/HTTPS checks"),
			},
		},
	}
}

func LoadbalMemberAddMetadata() cli.Command {
	return cli.Command{
		Category:    NS_LOADBAL_NAME,
		Name:        CMD_LOADBAL_MEMBER_ADD_NAME,
		Description: T("Add a new load balancer member"),
		Usage:       "${COMMAND_NAME} sl loadbal member-add (--id LOADBAL_ID) (--ip PRIVATE_IP)",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "id",
				Usage: T("ID for the load balancer [required]"),
			},
			cli.StringFlag{
				Name:  "ip",
				Usage: T("Private IP of the new member [required]"),
			},
		},
	}
}

func LoadbalMemberDelMetadata() cli.Command {
	return cli.Command{
		Category:    NS_LOADBAL_NAME,
		Name:        CMD_LOADBAL_MEMBER_DEL_NAME,
		Description: T("Remove a load balancer member"),
		Usage:       "${COMMAND_NAME} sl loadbal member-del (--lb-id LOADBAL_ID) (-m, --member-uuid MEMBER_UUID)",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "lb-id",
				Usage: T("ID for the load balancer [required]"),
			},
			cli.StringFlag{
				Name:  "m,member-uuid",
				Usage: T("Member UUID [required]"),
			},
		},
	}
}

func LoadbalProtocolAddMetadata() cli.Command {
	return cli.Command{
		Category:    NS_LOADBAL_NAME,
		Name:        CMD_LOADBAL_PROTOCOL_ADD_NAME,
		Description: T("Add a new load balancer protocol"),
		Usage:       "${COMMAND_NAME} sl loadbal protocol-add (--id LOADBAL_ID) [--front-protocol PROTOCOL] [back-protocol PROTOCOL] [--front-port PORT] [--back-port PORT] [-m, --method METHOD] [-c, --connections CONNECTIONS] [--sticky cookie | source-ip] [--client-timeout SECONDS] [--server-timeout SECONDS]",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "id",
				Usage: T("ID for the load balancer [required]"),
			},
			cli.StringFlag{
				Name:  "front-protocol",
				Usage: T("Protocol type to use for incoming connections: [HTTP|HTTPS|TCP]. Default: HTTP"),
			},
			cli.StringFlag{
				Name:  "back-protocol",
				Usage: T("Protocol type to use when connecting to backend servers: [HTTP|HTTPS|TCP]. Defaults to whatever --front-protocol is"),
			},
			cli.IntFlag{
				Name:  "front-port",
				Usage: T("Internet side port. Default: 80"),
			},
			cli.IntFlag{
				Name:  "back-port",
				Usage: T("Private side port. Default: 80"),
			},
			cli.StringFlag{
				Name:  "m, method",
				Usage: T("Balancing Method: [ROUNDROBIN|LEASTCONNECTION|WEIGHTED_RR]. Default: ROUNDROBIN"),
			},
			cli.IntFlag{
				Name:  "c, connections",
				Usage: T("Maximum number of connections to allow"),
			},
			cli.StringFlag{
				Name:  "sticky",
				Usage: T("Use 'cookie' or 'source-ip' to stick"),
			},
			cli.IntFlag{
				Name:  "client-timeout",
				Usage: T("Client side timeout setting, in seconds"),
			},
			cli.IntFlag{
				Name:  "server-timeout",
				Usage: T("Server side timeout setting, in seconds"),
			},
		},
	}
}

func LoadbalProtocolEditMetadata() cli.Command {
	return cli.Command{
		Category:    NS_LOADBAL_NAME,
		Name:        "protocol-edit",
		Description: T("Edit load balancer protocol"),
		Usage:       "${COMMAND_NAME} sl loadbal protocol-edit (--id LOADBAL_ID) (--protocol-uuid PROTOCOL_UUID) [--front-protocol PROTOCOL] [back-protocol PROTOCOL] [--front-port PORT] [--back-port PORT] [-m, --method METHOD] [-c, --connections CONNECTIONS] [--sticky cookie | source-ip] [--client-timeout SECONDS] [--server-timeout SECONDS]",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "id",
				Usage: T("ID for the load balancer [required]"),
			},
			cli.StringFlag{
				Name:  "protocol-uuid",
				Usage: T("UUID of the protocol you want to edit."),
			},
			cli.StringFlag{
				Name:  "front-protocol",
				Usage: T("Protocol type to use for incoming connections: [HTTP|HTTPS|TCP]. Default: HTTP"),
			},
			cli.StringFlag{
				Name:  "back-protocol",
				Usage: T("Protocol type to use when connecting to backend servers: [HTTP|HTTPS|TCP]. Defaults to whatever --front-protocol is"),
			},
			cli.IntFlag{
				Name:  "front-port",
				Usage: T("Internet side port. Default: 80"),
			},
			cli.IntFlag{
				Name:  "back-port",
				Usage: T("Private side port. Default: 80"),
			},
			cli.StringFlag{
				Name:  "m, method",
				Usage: T("Balancing Method: [ROUNDROBIN|LEASTCONNECTION|WEIGHTED_RR]. Default: ROUNDROBIN"),
			},
			cli.IntFlag{
				Name:  "c, connections",
				Usage: T("Maximum number of connections to allow"),
			},
			cli.StringFlag{
				Name:  "sticky",
				Usage: T("Use 'cookie' or 'source-ip' to stick"),
			},
			cli.IntFlag{
				Name:  "client-timeout",
				Usage: T("Client side timeout setting, in seconds"),
			},
			cli.IntFlag{
				Name:  "server-timeout",
				Usage: T("Server side timeout setting, in seconds"),
			},
		},
	}
}

func LoadbalProtocolDelMetadata() cli.Command {
	return cli.Command{
		Category:    NS_LOADBAL_NAME,
		Name:        CMD_LOADBAL_PROTOCOL_DELETE_NAME,
		Description: T("Delete a protocol"),
		Usage:       "${COMMAND_NAME} sl loadbal protocol-delete (--lb-id LOADBAL_ID) (--protocol-uuid PROTOCOL_UUID)",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "lb-id",
				Usage: T("ID for the load balancer [required]"),
			},
			cli.StringFlag{
				Name:  "protocol-uuid",
				Usage: T("UUID for the protocol [required]"),
			},
		},
	}
}

func LoadbalL7PoolDelMetadata() cli.Command {
	return cli.Command{
		Category:    NS_LOADBAL_NAME,
		Name:        CMD_LOADBAL_L7POOL_DELETE_NAME,
		Description: T("Delete a L7 pool"),
		Usage:       "${COMMAND_NAME} sl loadbal l7pool-delete (--pool-id L7POOL_ID)",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "pool-id",
				Usage: T("ID for the load balancer pool [required]"),
			},
		},
	}
}

func LoadbalL7PoolAddMetadata() cli.Command {
	return cli.Command{
		Category:    NS_LOADBAL_NAME,
		Name:        CMD_LOADBAL_L7POOL_ADD_NAME,
		Description: T("Add a new L7 pool"),
		Usage:       "${COMMAND_NAME} sl loadbal l7pool-add (--id LOADBAL_ID) (-n, --name NAME) [-m, --method METHOD] [-s, --server BACKEND_IP:PORT] [-p, --protocol PROTOCOL] [--health-path PATH] [--health-interval INTERVAL] [--health-retry RETRY] [--health-timeout TIMEOUT] [--sticky cookie | source-ip]",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "id",
				Usage: T("ID for the load balancer [required]"),
			},
			cli.StringFlag{
				Name:  "n, name",
				Usage: T("Name for this L7 pool. [required]"),
			},
			cli.StringFlag{
				Name:  "m, method",
				Usage: T("Balancing Method: [ROUNDROBIN|LEASTCONNECTION|WEIGHTED_RR]. [default: ROUNDROBIN]"),
			},
			cli.StringFlag{
				Name:  "p, protocol",
				Usage: T("Protocol type to use for incoming connections. [default: HTTP]"),
			},
			cli.StringSliceFlag{
				Name:  "s, server",
				Usage: T("Backend servers that are part of this pool. Format: BACKEND_IP:PORT. eg. 10.0.0.1:80 (multiple occurrence permitted)"),
			},
			cli.StringFlag{
				Name:  "health-path",
				Usage: T("Health check path.  [default: /]"),
			},
			cli.IntFlag{
				Name:  "health-interval",
				Usage: T("Health check interval between checks. [default: 5]"),
			},
			cli.IntFlag{
				Name:  "health-retry",
				Usage: T("Health check number of times before marking as DOWN. [default: 2]"),
			},
			cli.IntFlag{
				Name:  "health-timeout",
				Usage: T("Health check timeout. [default: 2]"),
			},
			cli.StringFlag{
				Name:  "sticky",
				Usage: T("Use 'cookie' or 'source-ip' to stick"),
			},
		},
	}

}

func LoadbalL7PoolDetailMetadata() cli.Command {
	return cli.Command{
		Category:    NS_LOADBAL_NAME,
		Name:        CMD_LOADBAL_L7POOL_DETAIL_NAME,
		Description: T("Show L7 pool details"),
		Usage:       "${COMMAND_NAME} sl loadbal l7pool-detail (--pool-id L7POOL_ID)",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "pool-id",
				Usage: T("ID for the load balancer pool [required]"),
			},
		},
	}
}

func LoadbalL7PoolEditMetadata() cli.Command {
	return cli.Command{
		Category:    NS_LOADBAL_NAME,
		Name:        CMD_LOADBAL_L7POOL_EDIT_NAME,
		Description: T("Edit a L7 pool"),
		Usage:       "${COMMAND_NAME} sl loadbal l7pool-edit (--pool-uuid L7POOL_UUID) [-m, --method METHOD] [-s, --server BACKEND_IP:PORT] [-p, --protocol PROTOCOL] [--health-path PATH] [--health-interval INTERVAL] [--health-retry RETRY] [--health-timeout TIMEOUT] [--sticky cookie | source-ip]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "pool-uuid",
				Usage: T("UUID for the load balancer pool [required]"),
			},
			cli.StringFlag{
				Name:  "m, method",
				Usage: T("Balancing Method: [ROUNDROBIN|LEASTCONNECTION|WEIGHTED_RR]"),
			},
			cli.StringFlag{
				Name:  "p, protocol",
				Usage: T("Protocol type to use for incoming connections"),
			},
			cli.StringSliceFlag{
				Name:  "s, server",
				Usage: T("Backend servers that are part of this pool. Format: BACKEND_IP:PORT. eg. 10.0.0.1:80 (multiple occurrence permitted)"),
			},
			cli.StringFlag{
				Name:  "health-path",
				Usage: T("Health check path"),
			},
			cli.IntFlag{
				Name:  "health-interval",
				Usage: T("Health check interval between checks"),
			},
			cli.IntFlag{
				Name:  "health-retry",
				Usage: T("Health check number of times before marking as DOWN"),
			},
			cli.IntFlag{
				Name:  "health-timeout",
				Usage: T("Health check timeout"),
			},
			cli.StringFlag{
				Name:  "sticky",
				Usage: T("Use 'cookie' or 'source-ip' to stick"),
			},
		},
	}
}

func LoadbalL7MemberAddMetadata() cli.Command {
	return cli.Command{
		Category:    NS_LOADBAL_NAME,
		Name:        CMD_LOADBAL_L7MEMBER_ADD_NAME,
		Description: T("Add a new L7 pool member"),
		Usage:       "${COMMAND_NAME} sl loadbal member-add (--pool-uuid L7POOL_UUID) (--address IP_ADDRESS) (--port PORT)",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "pool-uuid",
				Usage: T("UUID for the load balancer pool [required]"),
			},
			cli.StringFlag{
				Name:  "address",
				Usage: T("Backend servers IP address. [required]"),
			},
			cli.IntFlag{
				Name:  "port",
				Usage: T("Backend servers port. [required]"),
			},
		},
	}
}

func LoadbalL7MemberDeleteMetadata() cli.Command {
	return cli.Command{
		Category:    NS_LOADBAL_NAME,
		Name:        CMD_LOADBAL_L7MEMBER_DELETE_NAME,
		Description: T("Remove a load balancer member"),
		Usage:       "${COMMAND_NAME} sl loadbal l7member-del (--pool-uuid L7POOL_UUID) (--member-uuid L7MEMBER_UUID)",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "pool-uuid",
				Usage: T("UUID for the load balancer pool [required]"),
			},
			cli.StringFlag{
				Name:  "member-uuid",
				Usage: T("UUID for the load balancer member [required]"),
			},
		},
	}
}

func LoadbalL7PolicyAddMetadata() cli.Command {
	return cli.Command{
		Category:    NS_LOADBAL_NAME,
		Name:        CMD_LOADBAL_L7POLICY_ADD_NAME,
		Description: T("Add a new L7 policy"),
		Usage:       "${COMMAND_NAME} sl loadbal l7policy-add (--protocol-uuid PROTOCOL_UUID) (-n, --name NAME) (-a,--action REJECT | REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS) [-r,--redirect REDIRECT] [-p,--priority PRIORITY]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "protocol-uuid",
				Usage: T("UUID for the load balancer protocol [required]"),
			},
			cli.StringFlag{
				Name:  "n,name",
				Usage: T("Policy name"),
			},
			cli.StringFlag{
				Name:  "a,action",
				Usage: T("Policy action: REJECT | REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS"),
			},
			cli.StringFlag{
				Name:  "r,redirect",
				Usage: T("POOL_UUID, URL or HTTPS_PROTOCOL_UUID . It's only available in REDIRECT_POOL | REDIRECT_URL | REDIRECT_HTTPS action"),
			},
			cli.IntFlag{
				Name:  "p,priority",
				Usage: T("Policy priority"),
				Value: 1,
			},
		},
	}
}

func LoadbalL7PolicyDeleteMetadata() cli.Command {
	return cli.Command{
		Category:    NS_LOADBAL_NAME,
		Name:        CMD_LOADBAL_L7POLICY_DELETE_NAME,
		Description: T("Delete a L7 policy"),
		Usage:       "${COMMAND_NAME} sl loadbal l7policy-delete (--policy-id POLICY_ID) [-f, --force]",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "policy-id",
				Usage: T("ID for the load balancer policy [required]"),
			},
			ForceFlag(),
		},
	}
}

func LoadbalL7PolicyListMetadata() cli.Command {
	return cli.Command{
		Category:    NS_LOADBAL_NAME,
		Name:        CMD_LOADBAL_L7POLICY_LIST_NAME,
		Description: T("List L7 policies"),
		Usage:       "${COMMAND_NAME} sl loadbal l7policies (--protocol-id PROTOCOL_ID)",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "protocol-id",
				Usage: T("ID for the load balancer protocol [required]"),
			},
		},
	}
}

func LoadbalL7RuleAddMetadata() cli.Command {
	return cli.Command{
		Category:    NS_LOADBAL_NAME,
		Name:        CMD_LOADBAL_L7RULE_ADD_NAME,
		Description: T("Add a new L7 rule"),
		Usage:       "${COMMAND_NAME} sl loadbal l7rule-add (--policy-uuid L7POLICY_UUID) (-t, --type HOST_NAME | FILE_TYPE | HEADER | COOKIE | PATH ) (-c, --compare-type EQUAL_TO | ENDS_WITH | STARTS_WITH | REGEX | CONTAINS) (-v,--value VALUE) [-k,--key KEY] [--invert 0 | 1]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "policy-uuid",
				Usage: T("UUID for the load balancer policy [required]"),
			},
			cli.StringFlag{
				Name:  "t,type",
				Usage: T("Rule type: HOST_NAME | FILE_TYPE | HEADER | COOKIE | PATH. [required]"),
			},
			cli.StringFlag{
				Name:  "c,compare-type",
				Usage: T("Compare type: EQUAL_TO | ENDS_WITH | STARTS_WITH | REGEX | CONTAINS. [required]"),
			},
			cli.StringFlag{
				Name:  "v,value",
				Usage: T("Compared Value [required]"),
			},
			cli.StringFlag{
				Name:  "k,key",
				Usage: T("Key name. It's only available in HEADER or COOKIE type"),
			},
			cli.IntFlag{
				Name:  "invert",
				Usage: T("Invert rule: 0 | 1."),
			},
		},
	}
}

func LoadbalL7RuleDelMetadata() cli.Command {
	return cli.Command{
		Category:    NS_LOADBAL_NAME,
		Name:        CMD_LOADBAL_L7RULE_DELETE_NAME,
		Description: T("Delete a L7 rule"),
		Usage:       "${COMMAND_NAME} sl loadbal l7rule-delete (--policy-uuid L7POLICY_UUID) (--rule-uuid L7RULE_UUID) [-f, --force]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "policy-uuid",
				Usage: T("UUID for the load balancer policy [required]"),
			},
			cli.StringFlag{
				Name:  "rule-uuid",
				Usage: T("UUID for the load balancer rule [required]"),
			},
			ForceFlag(),
		},
	}
}

func LoadbalL7RuleListMetadata() cli.Command {
	return cli.Command{
		Category:    NS_LOADBAL_NAME,
		Name:        CMD_LOADBAL_L7RULE_LIST_NAME,
		Description: T("List l7 rules"),
		Usage:       "${COMMAND_NAME} sl loadbal l7rules (--policy-id Policy_ID)",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "policy-id",
				Usage: T("ID for the load balancer policy [required]"),
			},
		},
	}
}
