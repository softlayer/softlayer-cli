package testhelpers

import (
	"fmt"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"

	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

type FakeTransportHandler_True struct {
}

func NewFakeSoftlayerSession_True() *session.Session {
	return &session.Session{
		TransportHandler: FakeTransportHandler_True{},
	}
}

func (h FakeTransportHandler_True) DoRequest(sess *session.Session, service string, method string, args []interface{}, options *sl.Options, pResult interface{}) error {
	*pResult.(*bool) = true
	return nil
}

type FakeTransportHandler struct {
	FileNames []string
	ApiError  sl.Error
}

func (h FakeTransportHandler) DoRequest(sess *session.Session, service string, method string, args []interface{}, options *sl.Options, pResult interface{}) error {
	// fmt.Println("service:\t", service)
	// fmt.Println("method:\t", method)
	// fmt.Println("filenames:\t", h.FileNames)
	// for _, arg := range args {
	// 	fmt.Println("args:\t", arg)
	// }
	// if options.Id != nil {
	// 	fmt.Println("options-id:\t", *options.Id)
	// }
	// if options.Mask != "" {
	// 	fmt.Println("options-mask:\t", options.Mask)
	// }
	// if options.Filter != "" {
	// 	fmt.Println("options-filter:\t", options.Filter)
	// }
	if h.ApiError.StatusCode > 0 {
		return h.ApiError
	}
	b, err := readJsonTestFixtures(service, method, h.FileNames)
	if err != nil {

		slError := sl.Error{
			StatusCode: 500,
			Exception:  fmt.Sprintf("%v",err),
			Message:    "Erroring doing Fake Handling",
			Wrapped:    nil,
		}
		return slError
	}
	err = json.Unmarshal(b, pResult)
	//fmt.Println(pResult)
	return err
}

func NewFakeSoftlayerSession(fileNames []string) *session.Session {
	slError := sl.Error{
		StatusCode: 0,
		Exception:  "",
		Message:    "",
		Wrapped:    nil,
	}
	return &session.Session{
		TransportHandler: FakeTransportHandler{fileNames, slError},
	}
}

// Use this constructor to force DoRequests to return a SL error
func NewFakeSoftlayerSessionErrors(fileNames []string, slError sl.Error) *session.Session {
	return &session.Session{
		TransportHandler: FakeTransportHandler{fileNames, slError},
	}
}


// This function tries to find an appropriate JSON file to use as a response object.
// Fixtures are placed in the plugin/testfixtures directory in this patter:
// testfixtures/SoftLayer_Service/method.json : For general use
// testfixtures/SoftLayer_Service/method_id.json : Will be used if the ID in the request matches, otherwise fallback to general method
// testfixtures/SoftLayer_Service/method_specialCase.json : Will be used if this is in the fileNames array
func readJsonTestFixtures(service string, method string, fileNames []string) ([]byte, error) {
	wd, _ := os.Getwd()
	var fixture string
	scope := ".."

	baseFixture := filepath.Join(wd, scope, "testfixtures", service+"/"+method+".json")
	// fmt.Printf("baseFixture: %v \n", baseFixture)
	if len(fileNames) == 0 {
		fixture = baseFixture
	} else {
		if strings.Contains(wd, "plugin/commands") {
			scope += "/.."
		}
		//find the file name that matches the service and method name
		for _, filename := range fileNames {
			//fmt.Println("check file:" + filename)
			//If the file exists as is, just load and return it.
			if _, err := os.Stat(filepath.Join(wd, scope, "testfixtures", "services", filename)); err == nil {
				fixture = filepath.Join(wd, scope, "testfixtures", "services", filename)
				return ioutil.ReadFile(fixture) // #nosec
			}

			nameSegments := strings.Split(filename, "_")
			if nameSegments[0] == "SoftLayer" && (nameSegments[1] == "Account" || nameSegments[1] == "Ticket") {
				if len(nameSegments) == 3 {
					if service == nameSegments[0]+"_"+nameSegments[1] && method == nameSegments[2] {
						fixture = filepath.Join(wd, scope, "testfixtures", "services", service+"_"+method+".json")
						break
					}
				} else if len(nameSegments) == 4 {
					if service == nameSegments[0]+"_"+nameSegments[1] && method == nameSegments[2] {
						fixture = filepath.Join(wd, scope, "testfixtures", "services", service+"_"+method+"_"+nameSegments[3]+".json")
						break
					} else if service == nameSegments[0]+"_"+nameSegments[1]+"_"+nameSegments[2] && method == nameSegments[3] {
						fixture = filepath.Join(wd, scope, "testfixtures", "services", service+"_"+method+".json")
					}
				}
			} else if nameSegments[0] == "SoftLayer" && nameSegments[1] != "Account" {
				if len(nameSegments) == 4 {
					if service == nameSegments[0]+"_"+nameSegments[1]+"_"+nameSegments[2] && method == nameSegments[3] {
						fixture = filepath.Join(wd, scope, "testfixtures", "services", service+"_"+method+".json")
						break
					}
				} else if len(nameSegments) == 5 {
					if service == nameSegments[0]+"_"+nameSegments[1]+"_"+nameSegments[2] && method == nameSegments[3] {
						fixture = filepath.Join(wd, scope, "testfixtures", "services", service+"_"+method+"_"+nameSegments[4]+".json")
						break
					} else if service == nameSegments[0]+"_"+nameSegments[1]+"_"+nameSegments[2]+"_"+nameSegments[3] && method == nameSegments[4] {
						fixture = filepath.Join(wd, scope, "testfixtures", "services", service+"_"+method+".json")
						break
					}
				} else if len(nameSegments) == 6 {
					if service == nameSegments[0]+"_"+nameSegments[1]+"_"+nameSegments[2]+"_"+nameSegments[3]+"_"+nameSegments[4] && method == nameSegments[5] {
						fixture = filepath.Join(wd, scope, "testfixtures", "services", service+"_"+method+".json")
						break
					}
				}
			}
		}
	}
	if fixture != "" {
		//fmt.Println("read file:" + fixture)
		return ioutil.ReadFile(fixture) // #nosec
	}
	files := utils.StringSliceToString(fileNames)
	return nil, errors.New("failed to find test fixture file:serivce=" + service + ",method=" + method + ",files:" + files)
}
