package main

import (
	"fmt"
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

	err := doc.GenMarkdownTree(slMeta, "./markdown")
	if err != nil {
		fmt.Errorf(err.Error())
	}
	fmt.Printf("Jobs done.\n")

}