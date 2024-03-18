package managers

import (
	"encoding/json"
	"fmt"

	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

//counterfeiter:generate -o ../testhelpers/ . MetadataManager
type MetadataManager interface {
	CallAPIService(string, string, sl.Options, string) (string, error)
}

type metadataManager struct {
	Session *session.Session
}

func NewMetadataManager(session *session.Session) *metadataManager {
	return &metadataManager{
		session,
	}
}

func (call metadataManager) CallAPIService(service string, method string, options sl.Options, parameters string) (string, error) {
	call.Session.Endpoint = *sl.String("https://api.service.softlayer.com/rest/v3.1")

	var output string
	var arg []interface{}
	if parameters != "" {
		_parameters := []byte(parameters)
		unmarshalErr := json.Unmarshal(_parameters, &arg)
		if unmarshalErr != nil {
			fmt.Println(unmarshalErr.Error())
		}
	}
	err := call.Session.DoRequest(service, method, arg, &options, &output)
	return output, err
}
