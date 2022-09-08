package user

import (
	crand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"time"

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
		password = string(GeneratePassword(23, 4))
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

// random source leveraging crypto/rand to provide
// true non-determinstic
type cryptoSource struct{}

func (s cryptoSource) Seed(seed int64) {}

func (s cryptoSource) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

func (s cryptoSource) Uint64() (v uint64) {
	err := binary.Read(crand.Reader, binary.BigEndian, &v)
	if err != nil {
		log.Fatal(err)
	}
	return v
}

// GeneratePassword will create a random password
// Returns a 23 character random string
// 0  only number
// 1  lower and upper
// 2   upper
// 3   special
// 4  all
func GeneratePassword(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{{10, 48}, {26, 97}, {26, 65}, {10, 38}}, make([]byte, size)
	isAll := kind > 3 || kind < 0

	// #nosec G404: Use "crypto/rand" as the seed, which should resolve the pseudo "math/rand"
	rnd := rand.New(&cryptoSource{})
	fmt.Println(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if isAll { // random ikind
			ikind = rnd.Intn(4)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rnd.Intn(scope))
	}
	return result
}

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
