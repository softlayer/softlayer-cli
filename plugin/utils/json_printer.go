package utils

import (
	"encoding/json"
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
)

func PrintPrettyJSON(ui terminal.UI, data interface{}) error {
	jsonBytes, err := prettyJSON(data)
	if err != nil {
		return err
	}

	fmt.Fprintf(ui.Writer(), "%s\n", jsonBytes)
	return nil

}

func prettyJSON(data interface{}) (string, error) {
	jsonBytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
