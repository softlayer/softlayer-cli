package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

type EventLogManager interface {
	GetEventLogs(mask string, filter string, limit int) ([]datatypes.Event_Log, error)
	GetEventLogTypes() ([]string, error)
}

type eventLogManager struct {
	EventLogService services.Event_Log
}

func NewEventLogManager(session *session.Session) *eventLogManager {
	return &eventLogManager{
		services.GetEventLogService(session),
	}
}

// Get Event Logs
// mask: object mask
// dateFilter: object filter
// limit: limit of event logs
func (as eventLogManager) GetEventLogs(mask string, filter string, limit int) ([]datatypes.Event_Log, error) {
	return as.EventLogService.Limit(limit).Filter(filter).Mask(mask).GetAllObjects()
}

// Get Event Log Types
func (as eventLogManager) GetEventLogTypes() ([]string, error) {
	return as.EventLogService.GetAllEventObjectNames()
}
