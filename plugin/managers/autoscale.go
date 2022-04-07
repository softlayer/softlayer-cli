package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

type AutoScaleManager interface {
	GetLogsScaleGroup(id int, mask string, dateFilter string) ([]datatypes.Scale_Group_Log, error)
}

type autoScaleManager struct {
	AutoScaleService services.Scale_Group
	AccountService   services.Account
}

func NewAutoScaleManager(session *session.Session) *autoScaleManager {
	return &autoScaleManager{
		services.GetScaleGroupService(session),
		services.GetAccountService(session),
	}
}

//Get logs about specific autoscale group
//id: Auto Sacale Group Id
//mask: object mask
//dateFilter: Earliest date to retrieve logs for [YYYY-MM-DD]
func (as autoScaleManager) GetLogsScaleGroup(id int, mask string, dateFilter string) ([]datatypes.Scale_Group_Log, error) {
	if dateFilter != "" {
		filter := filter.New(filter.Path("logs.createDate").DateAfter(dateFilter))
		return as.AutoScaleService.Filter(filter.Build()).Id(id).Mask(mask).GetLogs()
	}
	return as.AutoScaleService.Id(id).Mask(mask).GetLogs()
}
