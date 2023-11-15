package main

import (
	"fmt"
	"os"
	"strings"
	"sort"
	"io/ioutil"
	"path/filepath"
	// "sort"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/cobra"
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

// This function builds the documentation for IBMCLOUD docs
func CliDocs() {
	fmt.Printf("IBMCLOUD SL Command Directory\n")

	slPlugin := new(sl_plugin.SoftlayerPlugin)
	slMeta := slPlugin.GetMetadata()
	sort.Slice(slMeta.Commands, func(i, j int) bool {
		one := fmt.Sprintf("%s %s", slMeta.Commands[i].Namespace, slMeta.Commands[i].Name)
		two := fmt.Sprintf("%s %s", slMeta.Commands[j].Namespace, slMeta.Commands[j].Name)
		return one < two
	})
	fmt.Printf("==============================================================\n")
	fileName := ""
	fileContent := ""
	// TODO: call-api, version and metadata need a special case or something for filename...
	for _, slCmd := range slMeta.Commands {
		thisFileName := fmt.Sprintf("cli_%s.md", slCmd.Namespace)
		thisFileName = strings.ReplaceAll(thisFileName, " ", "_")
		if thisFileName != fileName {
			fileName = thisFileName
			fmt.Printf("NameSpace: %s  Name: %s FIleName: %s\n", slCmd.Namespace, slCmd.Name, fileName)
			if fileContent != "" {
				fmt.Printf("Here is where I would write out to a file...\n")
				fileContent = ""
			}
		} 
		
		sort.Slice(slCmd.Flags, func(i, j int) bool {
			return slCmd.Flags[i].Name < slCmd.Flags[j].Name
		})
		// for _, slCmdFlag := range slCmd.Flags {
			// fmt.Printf("\tFlag: %s: %s\n", slCmdFlag.Name, slCmdFlag.Description)

		// }
		// fmt.Printf("\t--------------------------------\n")
		// fmt.Printf("\tDescription: %s\n", slCmd.Description)
		// fmt.Printf("\t--------------------------------\n")
		// fmt.Printf("\tUsage: %s\n", slCmd.Usage)
		// fmt.Printf("==============================================================\n")
	}
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