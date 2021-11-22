package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

//Manages SoftLayer Dedicated host.
type DedicatedHostManager interface {
	ListGuests(identifier int, cpu int, domain string, hostname string, memory int, tags []string, mask string) ([]datatypes.Virtual_Guest, error)
}

type dedicatedhostManager struct {
	AccountService       services.Account
	VirtualDedicatedHost services.Virtual_DedicatedHost
}

func NewDedicatedhostManager(session *session.Session) *dedicatedhostManager {
	return &dedicatedhostManager{
		services.GetAccountService(session),
		services.GetVirtualDedicatedHostService(session),
	}
}

//Retrieve a list of all virtual servers on the dedicated host.
//integer identifier: The identifier of a dedicated host.
//integer cpus: filter based on number of CPUS.
//string domain: filter based on domain.
//string hostname: filter based on hostname.
//integer memory: filter based on amount of memory.
//list tags: filter based on list of tags.
func (d dedicatedhostManager) ListGuests(identifier int, cpu int, domain string, hostname string, memory int, tags []string, mask string) ([]datatypes.Virtual_Guest, error) {
	filters := filter.New()
	if cpu != 0 {
		filters = append(filters, filter.Path("guests.maxCpu").Eq(cpu))
	}
	if domain != "" {
		filters = append(filters, utils.QueryFilter(domain, "guests.domain"))
	}
	if hostname != "" {
		filters = append(filters, utils.QueryFilter(hostname, "guests.hostname"))
	}
	if memory != 0 {
		filters = append(filters, filter.Path("guests.maxMemory").Eq(memory))
	}
	if len(tags) > 0 {
		tagInterfaces := make([]interface{}, len(tags))
		for i, v := range tags {
			tagInterfaces[i] = v
		}
		filters = append(filters, filter.Path("guests.tagReferences.tag.name").In(tagInterfaces...))
	}

	guestList, err := d.VirtualDedicatedHost.Id(identifier).Mask(mask).Filter(filters.Build()).GetGuests()
	if err != nil {
		return []datatypes.Virtual_Guest{}, err
	}
	return guestList, nil
}
