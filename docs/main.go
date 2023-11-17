package main

import (
	"fmt"
	"os"
	"strings"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	// "sort"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	// "github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	sl_plugin "github.ibm.com/SoftLayer/softlayer-cli/plugin"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"github.com/softlayer/softlayer-go/session"
)

var fileName string
var rootCmd = &cobra.Command{
	Use: "doc-gen",
	Short: "Generate the documentation for the sl plugin",
	RunE: func(Cmd *cobra.Command, args []string) error {
		CliDocs()
		return nil
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf(err.Error())
	}
	return
}

// For top level commands, like `sl account` or `sl hardware`
type SlCmdGroup struct {
	Name string
	Commands []SlCmdDoc
	Help string
}

// For specific commands
type SlCmdDoc struct {
	Name string
	Use string
	Flags []SlCmdFlag
	Help string
	LongHelp string
}

// For a commands flags
type SlCmdFlag struct {
	Name string
	Help string
}

// This function builds the documentation for IBMCLOUD docs
func CliDocs() {
	// fmt.Printf("IBMCLOUD SL Command Directory\n")
	SlCommands := sl_plugin.GetTopCobraCommand(nil, nil)
	CmdGroups := []SlCmdGroup{}
	for _, iCmd := range SlCommands.Commands() {
		thisCmdGroup := SlCmdGroup{
			Name: iCmd.Name(),
			Commands: nil,
			Help: iCmd.Short,
		}
		if len(iCmd.Commands()) > 0 {
			thisCmdGroup.Commands = buildSlCmdDoc(iCmd)
		}
		CmdGroups = append(CmdGroups, thisCmdGroup)
	}
	jOut, err := json.Marshal(CmdGroups)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jOut))
}

func buildSlCmdDoc(topCommand *cobra.Command) []SlCmdDoc {
	docs := []SlCmdDoc{}
	for _, iCmd := range topCommand.Commands() {
		thisDoc := SlCmdDoc{
			Name: iCmd.Name(),
			Use: iCmd.Use,
			Flags: nil,
			Help: iCmd.Short,
			LongHelp: iCmd.Long,
		}
		thisDoc.Flags = buildSlCmdFlag(iCmd)

		docs = append(docs, thisDoc)
	}
	return docs
}

func buildSlCmdFlag(topCommand *cobra.Command) []SlCmdFlag {
	flags := []SlCmdFlag{}
	flagSet := topCommand.Flags()
	flagSet.VisitAll(func(pflag *pflag.Flag) {
		thisFlag := SlCmdFlag{
			Name: pflag.Name,
			Help: pflag.Usage,
		}
		flags = append(flags, thisFlag)
	})
	return flags
}

// This function uses the build in Cobra documentation generator, its fine.
func CobraDocs() {
	fmt.Printf("Generating Documentation\n")

	var fakeUI              *terminal.FakeUI
	var fakeSession         *session.Session
	fakeUI = terminal.NewFakeUI()
	fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
	slMeta := sl_plugin.GetTopCobraCommand(fakeUI, fakeSession)

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf(err.Error())
	}
	if !strings.HasSuffix(filepath.ToSlash(cwd), "softlayer-cli/docs") {
		fmt.Printf("%v is the wrong directory, you need to run this command in the softlayer-cli/docs directory.\n", cwd)

		return
	} 
	err = doc.GenMarkdownTree(slMeta, "./")
	if err != nil {
		fmt.Printf(err.Error())
	}
	// err = os.Rename("./sl.md", "./index.md")
	// if err != nil {
	// 	fmt.Errorf(err.Error())
	// }
	// Need to make sure we have an index file
	bytesRead, err := ioutil.ReadFile("./sl.md")
    if err != nil {
        fmt.Printf(err.Error())
    }
    err = ioutil.WriteFile("./index.md", bytesRead, 0755)
    if err != nil {
        fmt.Printf(err.Error())
    }
	fmt.Printf("Jobs done.\n")

}