package callapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"

	"github.com/softlayer/softlayer-go/sl"
	// "github.com/softlayer/softlayer-go/session"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type CallAPICommand struct {
	UI              terminal.UI
	CallAPIManager  managers.CallAPIManager
	Init		    int
	Mask			string
	Parameters		string
	Limit			int
	Offset			int
	Filter			string
}


// func GetCommandActionBindings(ui terminal.UI, session *session.Session) map[string]*cobra.Command {
// 	CommandActionBindings := map[string]*cobra.Command {
// 		"call-api": NewCallAPICommand(ui, session),
// 	}
// 	return CommandActionBindings
// }


func NewCallAPICommand(sl *metadata.SoftlayerCommand) *cobra.Command {
	callAPIManager := managers.NewCallAPIManager(sl.Session)
	thisCmd := &CallAPICommand{
		UI:             sl.UI,
		CallAPIManager: callAPIManager,
	}

	cobraCmd := &cobra.Command{
		Use: "call-api",
		Short: T("Call arbitrary API endpoints"),
		Long: T(`${COMMAND_NAME} sl call-api SERVICE METHOD [OPTIONS]

EXAMPLE: 
	${COMMAND_NAME} sl call-api SoftLayer_Network_Storage editObject --init 57328245 --parameters '[{"notes":"Testing."}]'
	This command edit a volume notes.
	
	${COMMAND_NAME} sl call-api SoftLayer_User_Customer getObject --init 7051629 --mask "id,firstName,lastName"
	This command show a user detail.
	
	${COMMAND_NAME} sl call-api SoftLayer_Account getVirtualGuests --filter '{"virtualGuests":{"hostname":{"operation":"cli-test"}}}'
	This command list virtual guests.`),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}
	cobraCmd.Flags().IntVar(&thisCmd.Init, "init", 0, T("Init parameter"))
	cobraCmd.Flags().StringVar(&thisCmd.Mask, "mask", "", T("Object mask: use to limit fields returned"))
	cobraCmd.Flags().StringVar(&thisCmd.Parameters, "parameters", "", T("Append parameters to web call"))
	cobraCmd.Flags().IntVar(&thisCmd.Limit, "limit", 0, T("Result limit"))
	cobraCmd.Flags().IntVar(&thisCmd.Offset, "offset", 0, T("Result offset"))
	cobraCmd.Flags().StringVar(&thisCmd.Filter, "filter", "", T("Object filters"))

	return cobraCmd
}

func (cmd *CallAPICommand) Run(args []string)  error {
    var err error
    var output []byte
    var out bytes.Buffer
    var options sl.Options

    if cmd.Init != 0 {
    	options.Id = &cmd.Init
    }
    if !strings.HasPrefix(cmd.Mask, "mask[") && (strings.Contains(cmd.Mask, "[") || strings.Contains(cmd.Mask, ",")) {
        cmd.Mask = fmt.Sprintf("mask[%s]", cmd.Mask)
    }
    options.Mask = cmd.Mask

    if cmd.Offset != 0 {
    	options.Offset = &cmd.Offset
    }
    if cmd.Limit != 0 {
    	options.Limit = &cmd.Limit
    }
    if cmd.Filter != "" { 
    	options.Filter = cmd.Filter
    }

    output, err = cmd.CallAPIManager.CallAPI(args[0], args[1], options, cmd.Parameters)
    if err != nil {
    	return err
    }
    err = json.Indent(&out, output, "", "\t")
    if err != nil {
        _, err := cmd.UI.Writer().Write(output)
        if err != nil {
            return err
        }
    } else {
        _, err := out.WriteTo(cmd.UI.Writer())
        if err != nil {
            return err
        }
    }
	return nil
}
