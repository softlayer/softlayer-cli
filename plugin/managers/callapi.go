package managers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

type CallAPIManager interface {
	CallAPI(string, string, sl.Options, string) ([]byte, error)
}

type callAPIManager struct {
	Session *session.Session
}

func NewCallAPIManager(session *session.Session) *callAPIManager {
	return &callAPIManager{
		session,
	}
}

func (call callAPIManager) CallAPI(service string, method string, options sl.Options, parameters string) ([]byte, error) {

	if !strings.HasPrefix(service, "SoftLayer") {
		service = fmt.Sprintf("SoftLayer_%s", service)
	}

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

	return []byte(output), err
}
