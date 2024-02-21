package errors

import (
	"errors"
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

const (
	SL_EXP_OBJ_NOT_FOUND = "SoftLayer_Exception_ObjectNotFound"
)

func Error_Not_Login(context plugin.PluginContext) error {
	return errors.New(T("Please run command '{{.CommandName}} login' to login to IBM Cloud.",
		map[string]interface{}{"CommandName": context.CLIName()}))
}

func New(errorString string) error {
	return errors.New(errorString)
}

func CollapseErrors(multiErrors []error) error {
	errorString := ""
	for _, theError := range multiErrors {
		errorString = fmt.Sprintf("%v\n%v", errorString, theError.Error())
	}
	return errors.New(errorString)
}

// InvalidSoftlayerIdInputError represents a error about an invalid id input for softlayer
type InvalidSoftlayerIdInputError struct {
	InputName string
}

func NewInvalidSoftlayerIdInputError(inputName string) *InvalidSoftlayerIdInputError {
	return &InvalidSoftlayerIdInputError{
		InputName: inputName,
	}
}
func (err *InvalidSoftlayerIdInputError) Error() string {
	message := T("Invalid input for '{{.Name}}'. It must be a positive integer.", map[string]interface{}{"Name": err.InputName})
	return message
}

type APIError struct {
	CliMessage string
	APIMessage string
	ErrorCode  int
}

func NewAPIError(cliMessage string, apiMessage string, errorCode int) *APIError {
	return &APIError{
		CliMessage: cliMessage,
		APIMessage: apiMessage,
		ErrorCode:  errorCode,
	}
}

func (err *APIError) Error() string {
	message := err.CliMessage + "\n" + err.APIMessage
	return message
}
