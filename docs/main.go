package main

import (
	"fmt"
	"os"
	"strings"
	"io/ioutil"
	"path/filepath"
	// "sort"
	"github.com/spf13/cobra/doc"
	// "github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	sl_plugin "github.ibm.com/SoftLayer/softlayer-cli/plugin"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"github.com/softlayer/softlayer-go/session"
)


func main() {
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