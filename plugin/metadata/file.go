package metadata

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

var (
	NS_FILE_NAME  = "file"
	CMD_FILE_NAME = "file"

	//sl file
	CMD_FILE_ACCESS_AUTHORIZE_NAME       = "access-authorize"
	CMD_FILE_ACCESS_LIST_NAME            = "access-list"
	CMD_FILE_ACCESS_REVOKE_NAME          = "access-revoke"
	CMD_FILE_REPLICA_FAILBACK_NAME       = "replica-failback"
	CMD_FILE_REPLICA_FAILOVER_NAME       = "replica-failover"
	CMD_FILE_REPLICA_LOCATIONS_NAME      = "replica-locations"
	CMD_FILE_REPLICA_ORDER_NAME          = "replica-order"
	CMD_FILE_REPLICA_PARTNERS_NAME       = "replica-partners"
	CMD_FILE_SNAPSHOT_CANCEL_NAME        = "snapshot-cancel"
	CMD_FILE_SNAPSHOT_CREATE_NAME        = "snapshot-create"
	CMD_FILE_SNAPSHOT_DELETE_NAME        = "snapshot-delete"
	CMD_FILE_SNAPSHOT_DISABLE_NAME       = "snapshot-disable"
	CMD_FILE_SNAPSHOT_ENABLE_NAME        = "snapshot-enable"
	CMD_FILE_SNAPSHOT_LIST_NAME          = "snapshot-list"
	CMD_FILE_SNAPSHOT_ORDER_NAME         = "snapshot-order"
	CMD_FILE_SNAPSHOT_RESTORE_NAME       = "snapshot-restore"
	CMD_FILE_SNAPSHOT_SCHEDULE_LIST_NAME = "snapshot-schedule-list"
	CMD_FILE_VOLUME_CANCEL_NAME          = "volume-cancel"
	CMD_FILE_VOLUME_COUNT_NAME           = "volume-count"
	CMD_FILE_VOLUME_DETAIL_NAME          = "volume-detail"
	CMD_FILE_VOLUME_DUPLICATE_NAME       = "volume-duplicate"
	CMD_FILE_VOLUME_LIST_NAME            = "volume-list"
	CMD_FILE_VOLUME_ORDER_NAME           = "volume-order"
	CMD_FILE_VOLUME_MODIFY_NAME          = "volume-modify"
	CMD_FILE_VOLUME_OPTIONS_NAME         = "volume-options"
	CMD_FILE_VOLUME_LIMITS_NAME          = "volume-limits"
	CMD_FILE_VOLUME_REFRESH_NAME         = "volume-refresh"
	CMD_FILE_VOLUME_CONVERT_NAME         = "volume-convert"
)

func FileNamespace() plugin.Namespace {
	return plugin.Namespace{
		ParentName:  NS_SL_NAME,
		Name:        NS_FILE_NAME,
		Description: T("Classic infrastructure File Storage"),
	}
}

func FileMetaData() cli.Command {
	return cli.Command{
		Category:    NS_SL_NAME,
		Name:        CMD_FILE_NAME,
		Description: T("Classic infrastructure File Storage"),
		Usage:       "${COMMAND_NAME} sl file",
		Subcommands: []cli.Command{
			FileAccessAuthorizeMetaData(),
			FileAccessListMetaData(),
			FileAccessRevokeMetaData(),
			FileReplicaFailbackMetaData(),
			FileReplicaFailoverMetaData(),
			FileReplicaLocationsMetaData(),
			FileReplicaOrderMetaData(),
			FileReplicaPartnersMetaData(),
			FileSnapshotCancelMetaData(),
			FileSnapshotCreateMetaData(),
			FileSnapshotDisableMetaData(),
			FileSnapshotEnableMetaData(),
			FileSnapshotDeleteMetaData(),
			FileSnapshotListMetaData(),
			FileSnapshotOrderMetaData(),
			FileSnapshotScheduleListMetaData(),
			FileSnapshotRestoreMetaData(),
			FileVolumeCancelMetaData(),
			FileVolumeCountMetaData(),
			FileVolumeListMetaData(),
			FileVolumeDetailMetaData(),
			FileVolumeDeplicateMetaData(),
			FileVolumeModifyMetaData(),
			FileVolumeOrderMetaData(),
			FileVolumeOptionsMetaData(),
			FileVolumeLimitsMetaData(),
			FileVolumeRefreshMetaData(),
			FileVolumeConvertMetaData(),
			FileDisasterRecoveryFailoverMetaData(),
		},
	}
}

func FileAccessAuthorizeMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_ACCESS_AUTHORIZE_NAME,
		Description: T("Authorize hosts to access a given volume"),
		Usage: T(`${COMMAND_NAME} sl file access-authorize VOLUME_ID [OPTIONS]
		
EXAMPLE:
   ${COMMAND_NAME} sl file access-authorize 12345678 --virtual-id 87654321
   This command authorizes virtual server with ID 87654321 to access volume with ID 12345678.`),
		Flags: []cli.Flag{
			cli.IntSliceFlag{
				Name:  "d,hardware-id",
				Usage: T("The ID of one hardware server to authorize"),
			},
			cli.IntSliceFlag{
				Name:  "v,virtual-id",
				Usage: T("The ID of one virtual server to authorize"),
			},
			cli.IntSliceFlag{
				Name:  "i,ip-address-id",
				Usage: T("The ID of one IP address to authorize"),
			},
			cli.StringSliceFlag{
				Name:  "p,ip-address",
				Usage: T("An IP address to authorize"),
			},
			cli.IntSliceFlag{
				Name:  "s,subnet-id",
				Usage: T("An ID of one subnet to authorize"),
			},
			OutputFlag(),
		},
	}
}

func FileAccessListMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_ACCESS_LIST_NAME,
		Description: T("List hosts that are authorized to access the volume"),
		Usage: T(`${COMMAND_NAME} sl file access-list VOLUME_ID [OPTIONS]
		
EXAMPLE:
   ${COMMAND_NAME} sl file access-list 12345678 --sortby id 
   This command lists all hosts that are authorized to access volume with ID 12345678 and sorts them by ID.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "sortby",
				Usage: T("Column to sort by. Options are: id,name,type,private_ip_address,source_subnet,host_iqn,username,password,allowed_host_id"),
			},
			cli.StringSliceFlag{
				Name:  "column",
				Usage: T("Column to display. Options are: id,name,type,private_ip_address,source_subnet,host_iqn,username,password,allowed_host_id. This option can be specified multiple times"),
			},
			cli.StringSliceFlag{
				Name:   "columns",
				Hidden: true,
			},
			OutputFlag(),
		},
	}
}

func FileAccessRevokeMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_ACCESS_REVOKE_NAME,
		Description: T("Revoke authorization for hosts that are accessing a specific volume"),
		Usage: T(`${COMMAND_NAME} sl file access-revoke VOLUME_ID [OPTIONS]
		
EXAMPLE:
   ${COMMAND_NAME} sl file access-revoke 12345678 --virtual-id 87654321
   This command revokes access of virtual server with ID 87654321 to volume with ID 12345678.`),
		Flags: []cli.Flag{
			cli.IntSliceFlag{
				Name:  "d,hardware-id",
				Usage: T("The ID of one hardware server to revoke"),
			},
			cli.IntSliceFlag{
				Name:  "v,virtual-id",
				Usage: T("The ID of one virtual server to revoke"),
			},
			cli.IntSliceFlag{
				Name:  "i,ip-address-id",
				Usage: T("The ID of one IP address to revoke"),
			},
			cli.StringSliceFlag{
				Name:  "p,ip-address",
				Usage: T("An IP address to revoke"),
			},
			cli.IntSliceFlag{
				Name:  "s,subnet-id",
				Usage: T("An ID of one subnet to revoke"),
			},
		},
	}
}

func FileReplicaFailbackMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_REPLICA_FAILBACK_NAME,
		Description: T("Failback a file volume from replica"),
		Usage: T(`${COMMAND_NAME} sl file replica-failback VOLUME_ID
		
EXAMPLE:
   ${COMMAND_NAME} sl file replica-failback 12345678
   This command performs failback operation for volume with ID 12345678.`),
	}
}

func FileReplicaFailoverMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_REPLICA_FAILOVER_NAME,
		Description: T("Failover a file volume to the given replica volume"),
		Usage: T(`${COMMAND_NAME} sl file replica-failover VOLUME_ID REPLICA_ID
		
EXAMPLE:
   ${COMMAND_NAME} sl file replica-failover 12345678 87654321
   This command performs failover operation for volume with ID 12345678 to replica volume with ID 87654321.`),
	}
}

func FileReplicaLocationsMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_REPLICA_LOCATIONS_NAME,
		Description: T("List suitable replication datacenters for the given volume"),
		Usage: T(`${COMMAND_NAME} sl file replica-locations VOLUME_ID [OPTIONS]
		
EXAMPLE:
   ${COMMAND_NAME} sl file replica-locations 12345678
   This command lists suitable replication data centers for file volume with ID 12345678.`),
		Flags: []cli.Flag{
			OutputFlag(),
		},
	}
}

func FileReplicaOrderMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_REPLICA_ORDER_NAME,
		Description: T("Order a file storage replica volume"),
		Usage: T(`${COMMAND_NAME} sl file replica-order VOLUME_ID [OPTIONS]
		
EXAMPLE:
   ${COMMAND_NAME} sl file replica-order 12345678 -s DAILY -d dal09 --tier 4 
   This command orders a replica for volume with ID 12345678, which performs DAILY replication, is located at dal09, tier level is 4.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "s,snapshot-schedule",
				Usage: T("Snapshot schedule to use for replication. Options are: HOURLY,DAILY,WEEKLY [required]"),
			},
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Short name of the datacenter for the replica. For example, dal09 [required]"),
			},
			cli.Float64Flag{
				Name:  "t,tier",
				Usage: T("Endurance Storage Tier (IOPS per GB) of the primary volume for which a replica is ordered [optional], options are: 0.25,2,4,10,if no tier is specified, the tier of the original volume will be used"),
			},
			cli.IntFlag{
				Name:  "i,iops",
				Usage: T("Performance Storage IOPs, between 100 and 6000 in multiples of 100,if no IOPS value is specified, the IOPS value of the original volume will be used"),
			},
			ForceFlag(),
			OutputFlag(),
		},
	}
}

func FileReplicaPartnersMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_REPLICA_PARTNERS_NAME,
		Description: T("List existing replicant volumes for a file volume"),
		Usage: T(`${COMMAND_NAME} sl file replica-partners VOLUME_ID [OPTIONS]
		
EXAMPLE:
   ${COMMAND_NAME} sl file replica-partners 12345678
   This command lists existing replicant volumes for file volume with ID 12345678.`),
		Flags: []cli.Flag{
			OutputFlag(),
		},
	}
}

func FileSnapshotCancelMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_SNAPSHOT_CANCEL_NAME,
		Description: T("Cancel existing snapshot space for a given volume"),
		Usage: T(`${COMMAND_NAME} sl file snapshot-cancel SNAPSHOT_ID [OPTIONS]
		
EXAMPLE:
   ${COMMAND_NAME} sl file snapshot-cancel 12345678 --immediate -f 
   This command cancels snapshot with ID 12345678 immediately without asking for confirmation.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "reason",
				Usage: T("An optional reason for cancellation"),
			},
			cli.BoolFlag{
				Name:  "immediate",
				Usage: T("Cancel the snapshot space immediately instead of on the billing anniversary"),
			},
			ForceFlag(),
		},
	}
}

func FileSnapshotCreateMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_SNAPSHOT_CREATE_NAME,
		Description: T("Create a snapshot on a given volume"),
		Usage: T(`${COMMAND_NAME} sl file snapshot-create VOLUME_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl file snapshot-create 12345678 --note snapshotforibmcloud
   This command creates a snapshot for volume with ID 12345678 and with addition note as snapshotforibmcloud.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "n,note",
				Usage: T("Notes to set on the new snapshot"),
			},
			OutputFlag(),
		},
	}
}

func FileSnapshotDisableMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_SNAPSHOT_DISABLE_NAME,
		Description: T("Disable snapshots on the specified schedule for a given volume"),
		Usage: T(`${COMMAND_NAME} sl file snapshot-disable VOLUME_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl file snapshot-disable 12345678 -s DAILY
   This command disables daily snapshot for volume with ID 12345678.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "s,schedule-type",
				Usage: T("Snapshot schedule [required], options are: HOURLY,DAILY,WEEKLY"),
			},
		},
	}
}

func FileSnapshotEnableMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_SNAPSHOT_ENABLE_NAME,
		Description: T("Enable snapshots for a given volume on the specified schedule"),
		Usage: T(`${COMMAND_NAME} sl file snapshot-enable VOLUME_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl file snapshot-enable 12345678 -s WEEKLY -c 5 -m 0 --hour 2 -d 0
   This command enables snapshot for volume with ID 12345678, snapshot is taken weekly on every Sunday at 2:00, and up to 5 snapshots are retained.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "s,schedule-type",
				Usage: T("Snapshot schedule [required], options are: HOURLY,DAILY,WEEKLY"),
			},
			cli.IntFlag{
				Name:  "c,retention-count",
				Usage: T("Number of snapshots to retain [required]"),
			},
			cli.IntFlag{
				Name:  "m,minute",
				Usage: T("Minute of the hour when snapshots should be taken, integer between 0 to 59"),
			},
			cli.IntFlag{
				Name:  "r,hour",
				Usage: T("Hour of the day when snapshots should be taken, integer between 0 to 23"),
			},
			cli.IntFlag{
				Name:  "d,day-of-week",
				Usage: T("Day of the week when snapshots should be taken, integer between 0 to 6. \n      0 means Sunday,1 means Monday,2 means Tuesday,3 means Wendesday,4 means Thursday,5 means Friday,6 means Saturday"),
			},
		},
	}
}

func FileSnapshotDeleteMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_SNAPSHOT_DELETE_NAME,
		Description: T("Delete a snapshot on a given volume"),
		Usage: T(`${COMMAND_NAME} sl file snapshot-delete SNAPSHOT_ID

EXAMPLE:
   ${COMMAND_NAME} sl file snapshot-delete 12345678 
   This command deletes snapshot with ID 12345678.`),
	}
}

func FileSnapshotListMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_SNAPSHOT_LIST_NAME,
		Description: T("List file storage snapshots"),
		Usage: T(`${COMMAND_NAME} sl file snapshot-list VOLUME_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl file snapshot-list 12345678 --sortby id 
   This command lists all snapshots of volume with ID 12345678 and sorts them by ID.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "sortby",
				Usage: T("Column to sort by. Options are: id,name,created,size_bytes"),
			},
			// remove this flag because this command only has 4 columns no need to select
			// cli.StringSliceFlag{
			// 	Name:  CMD_FILE_SNAPSHOT_LIST_OPT2,
			// 	Usage: CMD_FILE_SNAPSHOT_LIST_OPT2_DESC,
			// },
			OutputFlag(),
		},
	}
}

func FileSnapshotOrderMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_SNAPSHOT_ORDER_NAME,
		Description: T("Order snapshot space for a file storage volume"),
		Usage: T(`${COMMAND_NAME} sl file snapshot-order VOLUME_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl file snapshot-order 12345678 -s 1000 -t 4 
   This commands order snapshot space for volume with ID 12345678, the size is 1000GB, the tier level is 4 IOPS per GB.`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "s,size",
				Usage: T("Size of snapshot space to create in GB  [required]"),
			},
			cli.Float64Flag{
				Name:  "t,tier",
				Usage: T("Endurance Storage Tier (IOPS per GB) of the file volume for which space is ordered [optional], options are: 0.25,2,4,10"),
			},
			cli.IntFlag{
				Name:  "i,iops",
				Usage: T("Performance Storage IOPs, between 100 and 6000 in multiples of 100"),
			},
			cli.BoolFlag{
				Name:  "u,upgrade",
				Usage: T("Flag to indicate that the order is an upgrade"),
			},
			ForceFlag(),
			OutputFlag(),
		},
	}
}

func FileSnapshotScheduleListMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_SNAPSHOT_SCHEDULE_LIST_NAME,
		Description: T("List snapshot schedules for a given volume"),
		Usage: T(`${COMMAND_NAME} sl file snapshot-schedule-list VOLUME_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl file snapshot-schedule-list 12345678
   This command list snapshot schedules for volume with ID 12345678`),
		Flags: []cli.Flag{
			OutputFlag(),
		},
	}
}

func FileSnapshotRestoreMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_SNAPSHOT_RESTORE_NAME,
		Description: T("Restore file volume using a given snapshot"),
		Usage: T(`${COMMAND_NAME} sl file snapshot-restore VOLUME_ID SNAPSHOT_ID
	
EXAMPLE:
   ${COMMAND_NAME} sl file snapshot-restore 12345678 87654321
   This command restores volume with ID 12345678 from snapshot with ID 87654321.`),
	}
}

func FileVolumeCancelMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_VOLUME_CANCEL_NAME,
		Description: T("Cancel an existing file storage volume"),
		Usage: T(`${COMMAND_NAME} sl file volume-cancel VOLUME_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl file volume-cancel 12345678 --immediate -f 
   This command cancels volume with ID 12345678 immediately and without asking for confirmation.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "reason",
				Usage: T("An optional reason for cancellation"),
			},
			cli.BoolFlag{
				Name:  "immediate",
				Usage: T("Cancel the file storage volume immediately instead of on the billing anniversary"),
			},
			ForceFlag(),
		},
	}
}

func FileVolumeCountMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_VOLUME_COUNT_NAME,
		Description: T("List number of file storage volumes per datacenter"),
		Usage:       "${COMMAND_NAME} sl file volume-count [OPTIONS]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Filter by datacenter shortname"),
			},
		},
	}
}

func FileVolumeListMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_VOLUME_LIST_NAME,
		Description: T("List file storage"),
		Usage: T(`${COMMAND_NAME} sl file volume-list [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl file volume-list -d dal09 -t endurance --sortby capacity_gb
   This command lists all endurance volumes on current account that are located at dal09, and sorts them by capacity.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "u,username",
				Usage: T("Filter by volume username"),
			},
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Filter by datacenter shortname"),
			},
			cli.StringFlag{
				Name:  "t,storage-type",
				Usage: T("Filter by type of storage volume, options are: performance,endurance"),
			},
			cli.IntFlag{
				Name:  "o,order",
				Usage: T("Filter by ID of the order that purchased the file storage"),
			},
			cli.StringFlag{
				Name:  "sortby",
				Usage: T("Column to sort by, default:id, options are: id,username,datacenter,storage_type,capacity_gb,bytes_used,ip_addr,active_transactions,created_by,mount_addr"),
			},
			cli.StringSliceFlag{
				Name:  "column",
				Usage: T("Column to display. Options are: id,username,datacenter,storage_type,capacity_gb,bytes_used,ip_addr,active_transactions,mount_addr,created_by,notes. This option can be specified multiple times"),
			},
			cli.StringSliceFlag{
				Name:   "columns",
				Hidden: true,
			},
			OutputFlag(),
		},
	}
}

func FileVolumeDetailMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_VOLUME_DETAIL_NAME,
		Description: T("Display details for a specified volume"),
		Usage: T(`${COMMAND_NAME} sl file volume-detail VOLUME_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl file volume-detail 12345678 
   This command shows details of volume with ID 12345678.`),
		Flags: []cli.Flag{
			OutputFlag(),
		},
	}
}

func FileVolumeDeplicateMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_VOLUME_DUPLICATE_NAME,
		Description: T("Order a file volume by duplicating an existing volume"),
		Usage: T(`${COMMAND_NAME} sl file volume-duplicate VOLUME_ID [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl file volume-duplicate 12345678 
   This command shows order a new volume by duplicating the volume with ID 12345678.`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "o,origin-snapshot-id",
				Usage: T("ID of an original volume snapshot to use for duplication"),
			},
			cli.IntFlag{
				Name:  "s,duplicate-size",
				Usage: T("Size of duplicate file volume in GB, if no size is specified, the size of the original volume will be used"),
			},
			cli.IntFlag{
				Name:  "i,duplicate-iops",
				Usage: T("Performance Storage IOPS, between 100 and 6000 in multiples of 100, if no IOPS value is specified, the IOPS value of the original volume will be used"),
			},
			cli.Float64Flag{
				Name:  "t,duplicate-tier",
				Usage: T("Endurance Storage Tier, if no tier is specified, the tier of the original volume will be used"),
			},
			cli.IntFlag{
				Name:  "n,duplicate-snapshot-size",
				Usage: T("The size of snapshot space to order for the duplicate, if no snapshot space size is specified, the snapshot space size of the original volume will be used"),
				Value: -1,
			},
			cli.BoolFlag{
				Name:  "d,dependent-duplicate",
				Usage: T("Whether or not this duplicate will be a dependent duplicate of the origin volume."),
			},
			ForceFlag(),
			OutputFlag(),
		},
	}
}

func FileVolumeModifyMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_VOLUME_MODIFY_NAME,
		Description: T("Modify an existing file storage volume"),
		Usage: T(`${COMMAND_NAME} sl file volume-modify VOLUME_ID [OPTIONS]

   EXAMPLE:
	  ${COMMAND_NAME} sl file volume-modify 12345678 --new-size 1000 --new-iops 4000 
	  This command modify a volume 12345678 with size is 1000GB, IOPS is 4000.
	  ${COMMAND_NAME} sl file volume-modify 12345678 --new-size 500 --new-tier 4
	  This command modify a volume 12345678 with size is 500GB, tier level is 4 IOPS per GB.`),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "c,new-size",
				Usage: T("New Size of file volume in GB. ***If no size is given, the original size of volume is used.***\n      Potential Sizes: [20, 40, 80, 100, 250, 500, 1000, 2000, 4000, 8000, 12000]\n      Minimum: [the original size of the volume]"),
			},
			cli.IntFlag{
				Name:  "i,new-iops",
				Usage: T("Performance Storage IOPS, between 100 and 6000 in multiples of 100 [only for performance volumes] ***If no IOPS value is specified, the original IOPS value of the volume will be used.***\n      Requirements: [If original IOPS/GB for the volume is less than 0.3, new IOPS/GB must also be less than 0.3. If original IOPS/GB for the volume is greater than or equal to 0.3, new IOPS/GB for the volume must also be greater than or equal to 0.3.]"),
			},
			cli.Float64Flag{
				Name:  "t, new-tier",
				Usage: T("Endurance Storage Tier (IOPS per GB) [only for endurance volumes] ***If no tier is specified, the original tier of the volume will be used.***\n      Requirements: [If original IOPS/GB for the volume is 0.25, new IOPS/GB for the volume must also be 0.25. If original IOPS/GB for the volume is greater than 0.25, new IOPS/GB for the volume must also be greater than 0.25.]"),
			},
			ForceFlag(),
			OutputFlag(),
		},
	}
}

func FileVolumeOrderMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_VOLUME_ORDER_NAME,
		Description: T("Order a file storage volume"),
		Usage: T(`${COMMAND_NAME} sl file volume-order [OPTIONS]

EXAMPLE:
   ${COMMAND_NAME} sl file volume-order --storage-type performance --size 1000 --iops 4000  -d dal09
   This command orders a performance volume with size is 1000GB, IOPS is 4000, located at dal09.
   ${COMMAND_NAME} sl file volume-order --storage-type endurance --size 500 --tier 4 -d dal09 --snapshot-size 500
   This command orders a endurance volume with size is 500GB, tier level is 4 IOPS per GB,located at dal09, and additional snapshot space size is 500GB.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "t,storage-type",
				Usage: T("Type of storage volume [required], options are: performance,endurance"),
			},
			cli.IntFlag{
				Name:  "s,size",
				Usage: T("Size of storage volume in GB [required]"),
			},
			cli.IntFlag{
				Name:  "i,iops",
				Usage: T("Performance Storage IOPs, between 100 and 6000 in multiples of 100 [required for storage-type performance]"),
			},
			cli.Float64Flag{
				Name:  "e,tier",
				Usage: T("Endurance Storage Tier (IOP per GB) [required for storage-type endurance], options are: 0.25,2,4,10"),
			},
			cli.StringFlag{
				Name:  "d,datacenter",
				Usage: T("Datacenter short name [required]"),
			},
			cli.IntFlag{
				Name:  "n,snapshot-size",
				Usage: T("Optional parameter for ordering snapshot space along with the volume"),
			},
			cli.StringFlag{
				Name:  "b,billing",
				Usage: T("Optional parameter for Billing rate (default to monthly), options are: hourly, monthly"),
			},
			ForceFlag(),
			OutputFlag(),
		},
	}
}

func FileVolumeOptionsMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_VOLUME_OPTIONS_NAME,
		Description: T("List all options for ordering a file storage"),
		Usage: T(`${COMMAND_NAME} sl file volume-options
	
EXAMPLE:
   ${COMMAND_NAME} sl file volume-options
   This command lists all options for creating a file storage volume, including storage type, volume size, IOPS, tier level, datacenter, and snapshot size.`),
	}
}

func FileVolumeLimitsMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_VOLUME_LIMITS_NAME,
		Description: T("Lists the storage limits per datacenter for this account."),
		Usage: T(`${COMMAND_NAME} sl file volume-limits [OPTIONS]

EXAMPLE:
	${COMMAND_NAME} sl file volume-limits
	This command lists the storage limits per datacenter for this account.`),
		Flags: []cli.Flag{
			OutputFlag(),
		},
	}
}

func FileVolumeRefreshMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_VOLUME_REFRESH_NAME,
		Description: T("Refresh a duplicate volume with a snapshot from its parent."),
		Usage: T(`${COMMAND_NAME} sl file volume-refresh VOLUME_ID SNAPSHOT_ID

EXAMPLE:
	${COMMAND_NAME} sl file volume-refresh VOLUME_ID SNAPSHOT_ID
	Refresh a duplicate VOLUME_ID with a snapshot from its parent SNAPSHOT_ID.`),
	}
}

func FileVolumeConvertMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_FILE_VOLUME_CONVERT_NAME,
		Description: T("Convert a dependent duplicate volume to an independent volume."),
		Usage: T(`${COMMAND_NAME} sl file volume-convert VOLUME_ID

EXAMPLE:
	${COMMAND_NAME} sl file volume-convert VOLUME_ID
	Convert a dependent duplicate VOLUME_ID to an independent volume.`),
	}
}


func FileDisasterRecoveryFailoverMetaData() cli.Command {
	return cli.Command{
		Category:    CMD_FILE_NAME,
		Name:        CMD_BLK_DISASTER_FAILOVER_NAME,
		Description: T("Failover an inaccessible volume to its available replicant volume."),
		Usage: T(`${COMMAND_NAME} sl file disaster-recovery-failover VOLUME_ID REPLICA_ID

If a volume (with replication) becomes inaccessible due to a disaster event, this method can be used to immediately
failover to an available replica in another location. This method does not allow for fail back via the API.
To fail back to the original volume after using this method, open a support ticket.
To test failover, use '${COMMAND_NAME} sl file replica-failover' instead.

EXAMPLE:
	${COMMAND_NAME} sl file disaster-recovery-failover 12345678 87654321
	This command performs failover operation for volume with ID 12345678 to replica volume with ID 87654321.`),
	}
}