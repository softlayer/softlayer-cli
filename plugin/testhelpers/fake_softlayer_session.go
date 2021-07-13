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
	ErrorMap map[string]sl.Error

}

func (h FakeTransportHandler) DoRequest(sess *session.Session, service string, method string, args []interface{}, options *sl.Options, pResult interface{}) error {
	// fmt.Println("\nservice:\t", service)
	// fmt.Println("method:\t", method)
	// fmt.Println("filenames:\t", h.FileNames)
	// for x, arg := range args {
	// 	fmt.Printf("args %v:\t %v", x, arg)
	// }

	identifier := 0
	if options.Id != nil {
		// fmt.Println("options-id:\t", *options.Id)
		identifier = *options.Id
	}
	apiSig := fmt.Sprintf("%s::%s", service, method)
	fmt.Printf("%s::%s(id=%d)\n", service, method, identifier)

	// if options.Mask != "" {
	// 	fmt.Println("options-mask:\t", options.Mask)
	// }
	// if options.Filter != "" {
	// 	fmt.Println("options-filter:\t", options.Filter)
	// }

	if apiError, ok := h.ErrorMap[apiSig]; ok {
		fmt.Printf("Found an error for %s -> %s\n", apiSig, apiError)
		return h.ErrorMap[apiSig]
	}
	if h.ApiError.StatusCode > 0 {
		return h.ApiError
	}
	// This is required to prevent pagination requests from going off in an infinite loop
	if options.Offset != nil  && *options.Offset > 0 {
		pResult = []byte("[]")
		return nil
	}
	b, err := readJsonTestFixtures(service, method, h.FileNames, identifier)
	if err != nil {

		slError := sl.Error{
			StatusCode: 555,
			Exception:  fmt.Sprintf("%v",err),
			Message:    "Erroring doing Fake Handling",
			Wrapped:    nil,
		}
		return slError
	}
	err = json.Unmarshal(b, pResult)
	return err
}

func (h *FakeTransportHandler) AddApiError(service string,  method string, errorCode int, errorMessage string) {
	if h.ErrorMap == nil {
		h.ErrorMap = make(map[string]sl.Error)
	}
	apiSig := service + "::" + method
	slError := sl.Error{
		StatusCode: errorCode,
		Exception:  errorMessage,
		Message:    errorMessage,
		Wrapped:    nil,
	}
	h.ErrorMap[apiSig] = slError
}

func (h FakeTransportHandler) ClearErrors() {
	h.ErrorMap = make(map[string]sl.Error)
}

func NewFakeSoftlayerSession(fileNames []string) *session.Session {
	slError := sl.Error{
		StatusCode: 0,
		Exception:  "",
		Message:    "",
		Wrapped:    nil,
	}
	errorMap := make(map[string]sl.Error)
	return &session.Session{
		TransportHandler: FakeTransportHandler{fileNames, slError, errorMap},
	}
}

// Use this constructor to force DoRequests to return a SL error
func NewFakeSoftlayerSessionErrors(errorCode int, message string) *session.Session {
	slError := sl.Error{
		StatusCode: errorCode,
		Exception:  message,
		Message:    message,
		Wrapped:    nil,
	}
	errorMap := make(map[string]sl.Error)
	return &session.Session{
		TransportHandler: FakeTransportHandler{nil, slError, errorMap},
	}
}


// This function tries to find an appropriate JSON file to use as a response object.
// Fixtures are placed in the plugin/testfixtures directory in this patter:
// testfixtures/SoftLayer_Service/method.json : For general use
// testfixtures/SoftLayer_Service/method_id.json : Will be used if the ID in the request matches, otherwise fallback to general method
// testfixtures/SoftLayer_Service/method_specialCase.json : Will be used if this is in the fileNames array
func readJsonTestFixtures(service string, method string, fileNames []string, identifier int) ([]byte, error) {
	wd, _ := os.Getwd()
	var fixture, workingPath string
	scope := ".."
	baseFixture := filepath.Join(wd, scope, "testfixtures", service+"/"+method+".json")

	if len(fileNames) == 0 {
		// Check to see if we have a fixture that matches the ID
		// actual path should be of the format testfixtures/SoftLayer_Service/method-123.json
		workingPath =  fmt.Sprintf("%s/%s-%d.json", service, method, identifier)
		if _, err := os.Stat(filepath.Join(wd, scope, "testfixtures", workingPath)); err == nil {
			fixture = filepath.Join(wd, scope, "testfixtures", workingPath)
			return ioutil.ReadFile(fixture) // #nosec
		}
	} else {
		if strings.Contains(wd, "plugin/commands") {
			scope += "/.."
		}

		//find the file name that matches the service and method name
		for _, filename := range fileNames {
			//fmt.Println("check file:" + filename)
			// If the file exists as is, just load and return it.
			// actual path should be of the format testfixtures/SoftLayer_Service/method-string.json
			workingPath = service + "/" + filename + ".json"
			if _, err := os.Stat(filepath.Join(wd, scope, "testfixtures", workingPath)); err == nil {
				fixture = filepath.Join(wd, scope, "testfixtures", workingPath)
				return ioutil.ReadFile(fixture) // #nosec
			}
		}
	}

	// Default to the base fixture `testfixtures/SoftLayer_Service/method.json`
	if _, err := os.Stat(baseFixture); err == nil {
		fixture = filepath.Join(baseFixture)
		return ioutil.ReadFile(fixture) // #nosec
	}

	fileNames = append(fileNames, baseFixture)
	apiCall := fmt.Sprintf("%s::%s(id=%d)", service, method, identifier)
	files := utils.StringSliceToString(fileNames)
	return nil, errors.New("Fixture for " + apiCall + " failed to load, looked in these files: " +  files)
}
