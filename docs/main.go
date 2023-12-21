package main

import (
	"fmt"
	"os"
	"strings"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"text/template"
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
	checkError(err)
	return
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// For top level commands, like `sl account` or `sl hardware`
type SlCmdGroup struct {
	Name string
	CommandShortLink string
	Commands []SlCmdDoc
	Help string
}

// For specific commands
type SlCmdDoc struct {
	Name string
	CommandShortLink string
	Use string
	Flags []SlCmdFlag
	Help string
	LongHelp string
	Backtick string
	CommandPath string
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
		shortName := strings.ReplaceAll(iCmd.Name(), " ", "_")
		shortName = strings.ReplaceAll(iCmd.Name(), "-", "_")
		thisCmdGroup := SlCmdGroup{
			Name: iCmd.Name(),
			CommandShortLink: fmt.Sprintf("sl_%v", shortName),
			Commands: nil,
			Help: iCmd.Short,

		}
		if len(iCmd.Commands()) > 0 {
			thisCmdGroup.Commands = buildSlCmdDoc(iCmd)
		}
		PrintMakrdown(thisCmdGroup)
		CmdGroups = append(CmdGroups, thisCmdGroup)
	}
	jOut, err := json.Marshal(CmdGroups)
	os.WriteFile("sl.json", jOut, 0755)
	checkError(err)
	// fmt.Println(string(jOut))
}

// Generates the Markdown
func PrintMakrdown(cmd SlCmdGroup) {

	var cmdTemplate = `
# ibmcloud sl {{.Name}}
{: #{{.CommandShortLink}}}

{{.Help}}

{{range .Commands}}
## ibmcloud {{.CommandPath}}
{: #{{.CommandShortLink}}}

{{.Help}}

{{.LongHelp}}

{{.Backtick}}bash
ibmcloud {{.Use}}
{{.Backtick}}
{: codeblock}

{{if .Flags}}
**Flags**:
{{range .Flags}}
	--{{.Name}} {{.Help}}
{{end}}
{{end}}
{{end}}

`
	mdTemplate, err := template.New("cmd template").Parse(cmdTemplate)
	checkError(err)
	filename := fmt.Sprintf("%v.md", cmd.CommandShortLink)
	outfile, err := os.Create(filename)
	defer outfile.Close()
	err = mdTemplate.Execute(outfile, cmd)
	checkError(err)

}


func buildSlCmdDoc(topCommand *cobra.Command) []SlCmdDoc {
	docs := []SlCmdDoc{}
	for _, iCmd := range topCommand.Commands() {
		shortName := fmt.Sprintf("sl_%s_%s", topCommand.Name(), iCmd.Name())
		shortName = strings.ReplaceAll(shortName, " ", "_")
		shortName = strings.ReplaceAll(shortName, "-", "_")

		thisDoc := SlCmdDoc{
			Name: iCmd.Name(),
			CommandShortLink: shortName,
			CommandPath: iCmd.CommandPath(),
			Use: iCmd.UseLine(),
			Flags: nil,
			Help: iCmd.Short,
			LongHelp: strings.ReplaceAll(iCmd.Long, "${COMMAND_NAME}", "ibmcloud"),
			Backtick:  "```",
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
	checkError(err)
	if !strings.HasSuffix(filepath.ToSlash(cwd), "softlayer-cli/docs") {
		fmt.Printf("%v is the wrong directory, you need to run this command in the softlayer-cli/docs directory.\n", cwd)

		return
	} 
	err = doc.GenMarkdownTree(slMeta, "./")
	checkError(err)
	// err = os.Rename("./sl.md", "./index.md")
	// if err != nil {
	// 	fmt.Errorf(err.Error())
	// }
	// Need to make sure we have an index file
	bytesRead, err := ioutil.ReadFile("./sl.md")
    checkError(err)
    err = ioutil.WriteFile("./index.md", bytesRead, 0755)
    checkError(err)
	fmt.Printf("Jobs done.\n")

}