package user

import (
	crand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"log"
	"math/rand"
	"reflect"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type CreateCommand struct {
	UI          terminal.UI
	UserManager managers.UserManager
}

func NewCreateCommand(ui terminal.UI, userManager managers.UserManager) (cmd *CreateCommand) {
	return &CreateCommand{
		UI:          ui,
		UserManager: userManager,
	}
}

func (cmd *CreateCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}
	userName := c.Args()[0]

	fromUserId := c.Int("from-user")

	var userTemplate datatypes.User_Customer
	var err error
	if !c.IsSet("from-user") {
		userTemplate, err = cmd.UserManager.GetCurrentUser()
	} else {
		userTemplate, err = cmd.UserManager.GetUser(fromUserId, "mask[id,firstName,lastName,email,companyName,address1,city,country,postalCode,state,userStatusId,timezoneId]")
	}

	if c.IsSet("template") {
		var templateStruct datatypes.User_Customer
		template := c.String("template")
		err := json.Unmarshal([]byte(template), &templateStruct)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf(T("Unable to unmarshal template json: %s\n"), err.Error()), 1)
		}
		StructAssignment(&userTemplate, &templateStruct)
	}

	userTemplate.Username = &userName

	password := c.String("password")
	if password == "generate" {
		password = string(GeneratePassword(23, 4))
	}

	vpnPassword := c.String("vpn-password")

	emailAddress := c.String("email")

	userTemplate.Email = &emailAddress

	if !c.IsSet("f") && !c.IsSet("force") {
		confirm, err := cmd.UI.Confirm(T("You are about to create the following user: {{.UserName}}. Do you wish to continue?", map[string]interface{}{"UserName": userName}))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
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
			return cli.NewExitError(T("Failed to add user.\n")+err.Error(), 2)
		}

	case nil:
		printUser(result, password, cmd.UI)

	default:
		return cli.NewExitError(T("Failed to add user.\n")+err.Error(), 2)
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

func UserCreateMetaData() cli.Command {
	return cli.Command{
		Category:    "user",
		Name:        "create",
		Description: T("Creates a user"),
		Usage: T(`${COMMAND_NAME} sl user create USERNAME [OPTIONS] 

EXAMPLE: 	
    ${COMMAND_NAME} sl user create my@email.com --email my@email.com --password generate --template '{"firstName": "Test", "lastName": "Testerson"}'
    Remember to set the permissions and access for this new user.`),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "email",
				Usage: T("Email address for this user. Required for creation"),
			},
			cli.StringFlag{
				Name:  "password",
				Usage: T("Password to set for this user. If no password is provided, the user is sent an email to generate one, which expires in 24 hours. Specify the '-p generate' option to generate a password for you. Passwords require 8+ characters, uppercase and lowercase, a number and a symbol"),
			},
			cli.IntFlag{
				Name:  "from-user",
				Usage: T("Base user to use as a template for creating this user. The default is to use the user that is running this command. Information provided in --template supersedes this template"),
			},
			cli.StringFlag{
				Name:  "template",
				Usage: T("A json string describing https://softlayer.github.io/reference/datatypes/SoftLayer_User_Customer/"),
			},
			cli.StringFlag{
				Name:  "vpn-password",
				Usage: T("VPN password to set for this user."),
			},
			metadata.ForceFlag(),
		},
	}
}
