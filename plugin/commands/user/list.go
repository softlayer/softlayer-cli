package user

import (
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

const (
	NO_ZERO_VALUE = "yes"
)

type ListCommand struct {
	*metadata.SoftlayerCommand
	UserManager managers.UserManager
	Command     *cobra.Command
	Column      []string
}

var maskMap = map[string]string{
	"id":                "id",
	"username":          "username",
	"email":             "email",
	"displayName":       "displayName",
	"status":            "userStatus.name",
	"hardwareCount":     "hardwareCount",
	"virtualGuestCount": "virtualGuestCount",
	"2FA":               "externalBindingCount",
	"classicAPIKey":     "apiAuthenticationKeyCount",
	"vpn":               "sslVpnAllowedFlag",
}

func NewListCommand(sl *metadata.SoftlayerCommand) (cmd *ListCommand) {
	thisCmd := &ListCommand{
		SoftlayerCommand: sl,
		UserManager:      managers.NewUserManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "list",
		Short: T("List Users"),
		Args:  metadata.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	columns := []string{}
	for key, _ := range maskMap {
		columns = append(columns, key)
	}
	sort.Strings(columns)
	subs := map[string]interface{}{"Columns": strings.Join(columns, ", ")}
	cobraCmd.Flags().StringSliceVar(&thisCmd.Column, "column", []string{},
		T("Column to display. options are: {{.Columns}}. This option can be specified multiple times", subs))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *ListCommand) Run(args []string) error {

	columns := cmd.Column

	defaultColumns := []string{"id", "username", "email", "displayName", "2FA", "classicAPIKey", "vpn"}
	optionalColumns := []string{"status", "hardwareCount", "virtualGuestCount"}

	showColumns, err := utils.ValidateColumns2("", columns, defaultColumns, optionalColumns, []string{})
	if err != nil {
		return err
	}

	mask := utils.GetMask(maskMap, showColumns, "")

	outputFormat := cmd.GetOutputFlag()

	users, err := cmd.UserManager.ListUsers(mask)
	if err != nil {
		return errors.NewAPIError(T("Failed to list users.\n"), err.Error(), 2)
	}

	if outputFormat == "JSON" {
		return utils.PrintPrettyJSON(cmd.UI, users)
	}

	table := cmd.UI.Table(showColumns)

	for _, user := range users {
		values := make(map[string]string)
		values["id"] = utils.FormatIntPointer(user.Id)
		values["username"] = utils.FormatStringPointer(user.Username)
		values["email"] = utils.FormatStringPointer(user.Email)
		values["displayName"] = utils.FormatStringPointer(user.DisplayName)

		values["2FA"] = utils.ReplaceUIntPointerValue(user.ExternalBindingCount, NO_ZERO_VALUE)
		values["classicAPIKey"] = utils.ReplaceUIntPointerValue(user.ApiAuthenticationKeyCount, NO_ZERO_VALUE)
		values["vpn"] = utils.FormatBoolPointerToYN(user.SslVpnAllowedFlag)

		if user.UserStatus != nil {
			values["status"] = utils.FormatStringPointer(user.UserStatus.Name)
		} else {
			values["status"] = utils.EMPTY_VALUE
		}

		values["hardwareCount"] = utils.FormatUIntPointer(user.HardwareCount)
		values["virtualGuestCount"] = utils.FormatUIntPointer(user.VirtualGuestCount)

		row := make([]string, len(showColumns))
		for i, col := range showColumns {
			row[i] = values[col]
		}
		table.Add(row...)
	}
	table.Print()
	return nil
}
