package dedicatedhost

import (
	"sort"
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/spf13/cobra"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type ListCommand struct {
	*metadata.SoftlayerCommand
	DedicatedHostManager managers.DedicatedHostManager
	Command              *cobra.Command
	Name                 string
	Datacenter           string
	Owner                string
	Order                int
	SortBy               string
}

type tableRow struct {
	Id         string
	Name       string
	Datacenter string
	Router     string
	Cpu        string
	Memory     string
	Disk       string
	Guests     string
}

func NewListCommand(sl *metadata.SoftlayerCommand) (cmd *ListCommand) {
	thisCmd := &ListCommand{
		SoftlayerCommand:     sl,
		DedicatedHostManager: managers.NewDedicatedhostManager(sl.Session),
	}
	cobraCmd := &cobra.Command{
		Use:   "list",
		Short: T("List dedicated hosts on your account"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	thisCmd.Command = cobraCmd
	cobraCmd.Flags().StringVarP(&thisCmd.Name, "name", "n", "", T("Filter by name of the dedicated host"))
	cobraCmd.Flags().StringVarP(&thisCmd.Datacenter, "datacenter", "d", "", T("Filter by datacenter of the dedicated host"))
	cobraCmd.Flags().StringVar(&thisCmd.Owner, "owner", "", T("Filter by owner of the dedicated host"))
	cobraCmd.Flags().IntVar(&thisCmd.Order, "order", 0, T("Filter by ID of the order which purchased this dedicated host"))
	cobraCmd.Flags().StringVar(&thisCmd.SortBy, "sortby", "", T("Column to sort by (Id, Name, Datacenter, Router, Cpu, Memory, Disk, Guests)[default: Id]"))
	return thisCmd
}

func (cmd *ListCommand) Run(args []string) error {
	sortBy := strings.ToLower(cmd.SortBy)
	if sortBy == "" {
		sortBy = "Id"
	} else {
		sortByOptions := []string{"id", "name", "datacenter", "router", "cpu", "memory", "disk", "guests"}
		if !utils.WordInList(sortByOptions, sortBy) {
			return errors.NewInvalidUsageError(T("Invalid --sortBy option."))
		}
	}

	outputFormat := cmd.GetOutputFlag()

	hosts, err := cmd.DedicatedHostManager.ListDedicatedHost(cmd.Name, cmd.Datacenter, cmd.Owner, cmd.Order)
	if err != nil {
		return errors.NewAPIError(T("Failed to list dedicated hosts on your account.\n"), err.Error(), 2)
	}

	if len(hosts) == 0 {
		cmd.UI.Print(T("No dedicated hosts are found."))
		return nil
	}

	// get array with rows
	tableRows := getTableRows(hosts)
	//sort host array
	switch sortBy {
	case "name":
		sort.Sort(ByName(tableRows))
	case "datacenter":
		sort.Sort(ByDatacenter(tableRows))
	case "router":
		sort.Sort(ByRouter(tableRows))
	case "cpu":
		sort.Sort(ByCpu(tableRows))
	case "memory":
		sort.Sort(ByMemory(tableRows))
	case "disk":
		sort.Sort(ByDisk(tableRows))
	case "guests":
		sort.Sort(ByGuests(tableRows))
	default:
		sort.Sort(ById(tableRows))
	}

	table := cmd.UI.Table([]string{T("Id"), T("Name"), T("Datacenter"), T("Router"), T("Cpu (allocated/total)"), T("Memory (allocated/total)"), T("Disk (allocated/total)"), T("Guests")})
	for _, row := range tableRows {
		table.Add(
			row.Id,
			row.Name,
			row.Datacenter,
			row.Router,
			row.Cpu,
			row.Memory,
			row.Disk,
			row.Guests,
		)
	}

	utils.PrintTable(cmd.UI, table, outputFormat)
	return nil
}

func getTableRows(hosts []datatypes.Virtual_DedicatedHost) []tableRow {
	tableRows := []tableRow{}
	row := tableRow{}
	for _, host := range hosts {
		row = tableRow{
			Id:         utils.FormatIntPointer(host.Id),
			Name:       utils.FormatStringPointer(host.Name),
			Datacenter: utils.FormatStringPointer(host.Datacenter.Name),
			Router:     utils.FormatStringPointer(host.BackendRouter.Hostname),
			Cpu:        utils.FormatIntPointer(host.AllocationStatus.CpuAllocated) + "/" + utils.FormatIntPointer(host.AllocationStatus.CpuCount),
			Memory:     utils.FormatIntPointer(host.AllocationStatus.MemoryAllocated) + "/" + utils.FormatIntPointer(host.AllocationStatus.MemoryCapacity),
			Disk:       utils.FormatIntPointer(host.AllocationStatus.DiskAllocated) + "/" + utils.FormatIntPointer(host.AllocationStatus.DiskCapacity),
			Guests:     utils.FormatUIntPointer(host.GuestCount),
		}
		tableRows = append(tableRows, row)
	}
	return tableRows
}

// interface to sort by Name
type ByName []tableRow

func (a ByName) Len() int           { return len(a) }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// interface to sort by Datacenter
type ByDatacenter []tableRow

func (a ByDatacenter) Len() int           { return len(a) }
func (a ByDatacenter) Less(i, j int) bool { return a[i].Datacenter < a[j].Datacenter }
func (a ByDatacenter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// interface to sort by Router
type ByRouter []tableRow

func (a ByRouter) Len() int           { return len(a) }
func (a ByRouter) Less(i, j int) bool { return a[i].Router < a[j].Router }
func (a ByRouter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// interface to sort by Cpu
type ByCpu []tableRow

func (a ByCpu) Len() int           { return len(a) }
func (a ByCpu) Less(i, j int) bool { return a[i].Cpu < a[j].Cpu }
func (a ByCpu) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// interface to sort by Memory
type ByMemory []tableRow

func (a ByMemory) Len() int           { return len(a) }
func (a ByMemory) Less(i, j int) bool { return a[i].Memory < a[j].Memory }
func (a ByMemory) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// interface to sort by Disk
type ByDisk []tableRow

func (a ByDisk) Len() int           { return len(a) }
func (a ByDisk) Less(i, j int) bool { return a[i].Disk < a[j].Disk }
func (a ByDisk) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// interface to sort by Guests
type ByGuests []tableRow

func (a ByGuests) Len() int           { return len(a) }
func (a ByGuests) Less(i, j int) bool { return a[i].Guests < a[j].Guests }
func (a ByGuests) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// interface to sort by Id
type ById []tableRow

func (a ById) Len() int           { return len(a) }
func (a ById) Less(i, j int) bool { return a[i].Id < a[j].Id }
func (a ById) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
