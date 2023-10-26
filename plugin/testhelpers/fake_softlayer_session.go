package testhelpers

import (
	"encoding/json"
	"errors"
	"fmt"
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
	FileNames   []string
	ApiError    sl.Error
	ErrorMap    map[string]sl.Error
	ApiCallLogs []ApiCallLog
}

type ApiCallLog struct {
	Service string
	Method  string
	Args    []interface{}
	Options *sl.Options
}

func (h *FakeTransportHandler) DoRequest(sess *session.Session, service string, method string, args []interface{}, options *sl.Options, pResult interface{}) error {

	if options == nil {
		options = new(sl.Options)
	}

	identifier := 0
	apiSig := fmt.Sprintf("%s::%s", service, method)

	if options.Id != nil {
		identifier = *options.Id
	}

	// fmt.Printf("%s::%s(id=%d)\n", service, method, identifier)

	h.AddApiLog(service, method, args, options)

	// If we have an error defined for this method, return that.
	if apiError, ok := h.ErrorMap[apiSig]; ok {
		return apiError
	}

	// This is required to prevent pagination requests from going off in an infinite loop
	if options.Offset != nil && *options.Offset > 0 {
		pResult = []byte("[]")
		return nil
	}

	// This fakes getting data from the SL API.
	b, err := readJsonTestFixtures(service, method, h.FileNames, identifier)

	// Incase of file not found, or other JSON errors, this presents the error somewhat nicely to the cli
	if err != nil {
		slError := sl.Error{
			StatusCode: 555,
			Exception:  fmt.Sprintf("%v", err),
			Message:    "Erroring doing Fake Handling",
			Wrapped:    nil,
		}
		return slError
	}
	err = json.Unmarshal(b, pResult)
	if err != nil {
		slError := sl.Error{
			StatusCode: 559,
			Exception:  fmt.Sprintf("%v", err),
			Message:    "Erroring doing json.Unmarshal",
			Wrapped:    nil,
		}
		return slError
	}
	return err
}

// Logs whenever the API is called.
func (h *FakeTransportHandler) AddApiLog(service string, method string, args []interface{}, options *sl.Options) {
	apiLog := ApiCallLog{
		Service: service,
		Method:  method,
		Args:    args,
		Options: options,
	}
	h.ApiCallLogs = append(h.ApiCallLogs, apiLog)
}

// Will return an error when the service+method API is called
func (h *FakeTransportHandler) AddApiError(service string, method string, errorCode int, errorMessage string) {
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

func (h *FakeTransportHandler) ClearApiCallLogs() {
	h.ApiCallLogs = []ApiCallLog{}
}

func (h *FakeTransportHandler) ClearErrors() {
	h.ErrorMap = make(map[string]sl.Error)
}

func NewFakeSoftlayerSession(fileNames []string) *session.Session {

	sess := &session.Session{}
	sess.TransportHandler = NewFakeTransportHandler(fileNames)
	return sess
}

func NewFakeTransportHandler(fileNames []string) session.TransportHandler {
	slError := sl.Error{
		StatusCode: 0,
		Exception:  "",
		Message:    "",
		Wrapped:    nil,
	}
	errorMap := make(map[string]sl.Error)
	apiCallLogs := []ApiCallLog{}
	var transportHandler session.TransportHandler
	transportHandler = &FakeTransportHandler{fileNames, slError, errorMap, apiCallLogs}
	return transportHandler
}

// Casts the session transport handler to the FakeTransportHandler
func GetSessionHandler(sess *session.Session) *FakeTransportHandler {
	transport := sess.TransportHandler.(*FakeTransportHandler)
	return transport
}

// This function tries to find an appropriate JSON file to use as a response object.
// Fixtures are placed in the plugin/testfixtures directory in this patter:
// testfixtures/SoftLayer_Service/method.json : For general use
// testfixtures/SoftLayer_Service/method-id.json : Will be used if the ID in the request matches, otherwise fallback to general method
// testfixtures/SoftLayer_Service/method-specialCase.json : Will be used if this is in the fileNames array
func readJsonTestFixtures(service string, method string, fileNames []string, identifier int) ([]byte, error) {
	wd, _ := os.Getwd()
	var fixture, workingPath string
	scope := ".."

	// The second check is for windows
	if strings.Contains(wd, "plugin/commands") || strings.Contains(wd, "plugin\\commands") {
		scope += "/.."
	}
	// fmt.Printf("WD: %v, Scope: %v", wd, scope)
	baseFixture := filepath.Join(wd, scope, "testfixtures", service+"/"+method+".json")

	// If we specified a file name, check there first
	if len(fileNames) > 0 {
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

	// Check for a matchin SoftLayer_Service/method-1234.json file
	workingPath = fmt.Sprintf("%s/%s-%d.json", service, method, identifier)
	if _, err := os.Stat(filepath.Join(wd, scope, "testfixtures", workingPath)); err == nil {
		fixture = filepath.Join(wd, scope, "testfixtures", workingPath)
		return ioutil.ReadFile(fixture) // #nosec
	}
	// Default to the base fixture `testfixtures/SoftLayer_Service/method.json`
	if _, err := os.Stat(baseFixture); err == nil {
		fixture = filepath.Join(baseFixture)
		return ioutil.ReadFile(fixture) // #nosec
	}

	fileNames = append(fileNames, baseFixture)
	apiCall := fmt.Sprintf("%s::%s(id=%d)", service, method, identifier)
	files := utils.StringSliceToString(fileNames)
	return nil, errors.New("Fixture for " + apiCall + " failed to load, looked in these files: " + files)
}
