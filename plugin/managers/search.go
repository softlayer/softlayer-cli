package managers

import (
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

//counterfeiter:generate -o ../testhelpers/ . SearchManager
type SearchManager interface {
	AdvancedSearch(mask string, params string) ([]datatypes.Container_Search_Result, error)
	GetTypes() ([]datatypes.Container_Search_ObjectType, error)
}

type searchManager struct {
	SearchService services.Search
	Session       *session.Session
}

func NewSearchManager(session *session.Session) *searchManager {
	return &searchManager{
		services.GetSearchService(session),
		session,
	}
}

/*
https://sldn.softlayer.com/reference/services/SoftLayer_Search/advancedSearch/
*/
func (s searchManager) AdvancedSearch(mask string, params string) ([]datatypes.Container_Search_Result, error) {

	return s.SearchService.Mask(mask).AdvancedSearch(&params)
}


func (s searchManager) GetTypes() ([]datatypes.Container_Search_ObjectType, error) {
	return s.SearchService.GetObjectTypes()
}