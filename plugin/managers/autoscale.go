package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

type AutoScaleManager interface {
	GetLogsScaleGroup(id int, mask string, filter string) ([]datatypes.Scale_Group_Log, error)
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
//dateFilter: object filter
func (as autoScaleManager) GetLogsScaleGroup(id int, mask string, filter string) ([]datatypes.Scale_Group_Log, error) {
	return as.AutoScaleService.Filter(filter).Id(id).Mask(mask).GetLogs()
}
