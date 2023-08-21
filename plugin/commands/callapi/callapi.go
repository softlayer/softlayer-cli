package callapi

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/softlayer/softlayer-go/sl"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

type CallAPICommand struct {
	*metadata.SoftlayerCommand
	CallAPIManager managers.CallAPIManager
	Command        *cobra.Command
	Init           int
	Mask           string
	Parameters     string
	Limit          int
	Offset         int
	Filter         string
}

func NewCallAPICommand(sl *metadata.SoftlayerCommand) *CallAPICommand {
	callAPIManager := managers.NewCallAPIManager(sl.Session)
	thisCmd := &CallAPICommand{
		SoftlayerCommand: sl,
		CallAPIManager:   callAPIManager,
	}

	cobraCmd := &cobra.Command{
		Use:   "call-api",
		Short: T("Call arbitrary API endpoints"),
		Long: T(`${COMMAND_NAME} sl call-api SERVICE METHOD [OPTIONS]

EXAMPLE: 
	${COMMAND_NAME} sl call-api SoftLayer_Network_Storage editObject --init 57328245 --parameters '[{"notes":"Testing."}]'
	This command edit a volume notes.
	
	${COMMAND_NAME} sl call-api SoftLayer_User_Customer getObject --init 7051629 --mask "id,firstName,lastName"
	This command show a user detail.
	
	${COMMAND_NAME} sl call-api SoftLayer_Account getVirtualGuests --filter '{"virtualGuests":{"hostname":{"operation":"cli-test"}}}'
	This command list virtual guests.`),
		Args: metadata.ExactArgs(2),
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

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *CallAPICommand) Run(args []string) error {
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

	if len(output) == 0 {
		cmd.UI.Print("Null")
		return nil
	}

	// 1
	// var jsonMap []map[string]interface{}
	// err = json.Unmarshal([]byte(output), &jsonMap)
	// if err != nil {
	// 	fmt.Println("Error parsing JSON data:", err)
	// }
	// fmt.Println("JSON data:", jsonMap)

	// keys := make([]int, 0, len(jsonMap))

	// for _, dataMap := range jsonMap {
	// 	for _, value := range dataMap {
	// 		keys = append(keys, value)
	// 	}
	// }

	// fmt.Println(jsonMap)
	// fmt.Println("Before Keys: \n", keys)

	// sort.Sort(jsonMap)
	// sort.SliceStable(keys, func(i, j int) bool {
	// 	return keys[i] < keys[j]
	// })

	// fmt.Println("After Keys: \n", keys)

	//2
	// data := DataSlice{{"key1":4567,"key2":"def"},{"key1":1234,"key2":"abc"},}
	// type DataMap map[string]interface{}
	// type DataSlice []DataMap

	// func (s DataSlice) Len() int{
	// 	return len(s)
	// }

	// func (s DataSlice) Less(i, j int) bool {
	// 	val1, ok1 := s[i]["key1"].(int)
	// 	val2, ok2 := s[j]["key1"].(int)

	// 	if !ok || !ok2 {
	// 		return false
	// 	}
	// 	return val1 < val2
	// }

	// func (s DataSlice) Swap(i, j int){
	// 	s[i], s[j] = s[j], s[i]
	// }

	// fmt.Println("\nBefore Sorting: ")

	// sort.Sort(data)

	// fmt.Println("\nAfter Sorting: ")

	// err = json.Indent(&out, output, "", "\t")
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
