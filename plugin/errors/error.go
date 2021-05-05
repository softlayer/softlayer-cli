package errors

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

const (
	SL_EXP_OBJ_NOT_FOUND = "SoftLayer_Exception_ObjectNotFound"
)

func Error_Not_Login(context plugin.PluginContext) error {
	return errors.New(T("Please run command '{{.CommandName}} login' to login to IBM Cloud.",
		map[string]interface{}{"CommandName": context.CLIName()}))
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
	message := i18n.T("Invalid input for '{{.Name}}'. It must be a positive integer.", map[string]interface{}{"Name": err.InputName})
	return message
}
