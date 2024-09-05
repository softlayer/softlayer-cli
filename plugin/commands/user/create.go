package user

import (
	"encoding/json"
	"fmt"
	"reflect"

	gopass "github.com/sethvargo/go-password/password"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/spf13/cobra"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CreateCommand struct {
	*metadata.SoftlayerCommand
	UserManager managers.UserManager
	Command     *cobra.Command
	Email       string
	Password    string
	FromUser    int
	Template    string
	VpnPassword string
	ForceFlag   bool
}

func NewCreateCommand(sl *metadata.SoftlayerCommand) (cmd *CreateCommand) {
	thisCmd := &CreateCommand{
		SoftlayerCommand: sl,
		UserManager:      managers.NewUserManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "create " + T("USERNAME"),
		Short: T("Creates a user"),
		Long: T(`${COMMAND_NAME} sl user create USERNAME [OPTIONS] 

EXAMPLE: 	
    ${COMMAND_NAME} sl user create my@email.com --email my@email.com --password generate --template '{"firstName": "Test", "lastName": "Testerson"}'
    Remember to set the permissions and access for this new user.`),
		Args: metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().StringVar(&thisCmd.Email, "email", "", T("Email address for this user. Required for creation"))
	cobraCmd.Flags().StringVar(&thisCmd.Password, "password", "", T("Password to set for this user. If no password is provided, the user is sent an email to generate one, which expires in 24 hours. Specify the '-p generate' option to generate a password for you. Passwords require 8+ characters, uppercase and lowercase, a number and a symbol"))
	cobraCmd.Flags().IntVar(&thisCmd.FromUser, "from-user", 0, T("Base user to use as a template for creating this user. The default is to use the user that is running this command. Information provided in --template supersedes this template"))
	cobraCmd.Flags().StringVar(&thisCmd.Template, "template", " ", T("A json string describing https://softlayer.github.io/reference/datatypes/SoftLayer_User_Customer/"))
	cobraCmd.Flags().StringVar(&thisCmd.VpnPassword, "vpn-password", "", T("VPN password to set for this user."))
	cobraCmd.Flags().BoolVarP(&thisCmd.ForceFlag, "force", "f", false, T("Force operation without confirmation"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CreateCommand) Run(args []string) error {
	userName := args[0]

	fromUserId := cmd.FromUser

	var userTemplate datatypes.User_Customer
	var err error
	if fromUserId == 0 {
		userTemplate, err = cmd.UserManager.GetCurrentUser()
	} else {
		userTemplate, err = cmd.UserManager.GetUser(fromUserId, "mask[id,firstName,lastName,email,companyName,address1,city,country,postalCode,state,userStatusId,timezoneId]")
	}

	if cmd.Template != " " {
		var templateStruct datatypes.User_Customer
		template := cmd.Template
		err := json.Unmarshal([]byte(template), &templateStruct)
		if err != nil {
			return errors.NewInvalidUsageError(fmt.Sprintf(T("Unable to unmarshal template json: %s\n"), err.Error()))
		}
		StructAssignment(&userTemplate, &templateStruct)
	}

	userTemplate.Username = &userName

	password := cmd.Password
	if password == "generate" {
		password = gopass.MustGenerate(18, 4, 4, false, false)
	}

	vpnPassword := cmd.VpnPassword

	emailAddress := cmd.Email

	userTemplate.Email = &emailAddress

	if !cmd.ForceFlag {
		confirm, err := cmd.UI.Confirm(T("You are about to create the following user: {{.UserName}}. Do you wish to continue?", map[string]interface{}{"UserName": userName}))
		if err != nil {
			return err
		}
		if !confirm {
			cmd.UI.Print(T("Aborted."))
			return nil
		}
	}

	result, err := cmd.UserManager.CreateUser(userTemplate, password, vpnPassword)
	switch err.(type) {
	case sl.Error:
		if err.(sl.Error).Exception == "SoftLayer_Exception_User_Customer_DelegateIamIdInvitationToPaas" {
			cmd.UI.Print(err.(sl.Error).Message)
		} else {
			return errors.NewAPIError(T("Failed to add user.\n"), err.Error(), 2)
		}

	case nil:
		printUser(result, password, cmd.UI)

	default:
		return errors.NewAPIError(T("Failed to add user.\n"), err.Error(), 2)
	}
	return nil
}

func printUser(user datatypes.User_Customer, password string, ui terminal.UI) {
	table := ui.Table([]string{T("name"), T("value")})
	table.Add(T("Username"), utils.FormatStringPointer(user.Username))
	table.Add(T("Email"), utils.FormatStringPointer(user.Email))
	table.Add(T("Password"), password)

	table.Print()
}

// Values of B get copied into A
// A <--- B
func StructAssignment(A, B interface{}) { //a =b
	av := reflect.ValueOf(A).Elem()
	at := av.Type()

	bv := reflect.ValueOf(B).Elem()
	bt := bv.Type()

	for k := 0; k < av.NumField(); k++ {
		for j := 0; j < bv.NumField(); j++ {
			if at.Field(k).Name == bt.Field(j).Name && bv.Field(k).Kind() != reflect.Struct && !bv.Field(k).IsNil() {
				tmp := bv.Field(j).Elem()
				av.Field(k).Set(tmp.Addr())
			}
		}
	}
}
