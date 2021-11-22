package virtual_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"strings"
)

var _ = Describe("VS capacity create", func() {
	var (
		fakeUI           *terminal.FakeUI
		fakeVSManager    *testhelpers.FakeVirtualServerManager
		cmd              *virtual.CapacityCreateCommand
		cliCommand       cli.Command
		context          plugin.PluginContext
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		context = plugin.InitPluginContext("softlayer")
		cmd = virtual.NewCapacityCreateCommand(fakeUI, fakeVSManager, context)
		cliCommand = cli.Command{
			Name:        metadata.VSCapacityCreateMetaData().Name,
			Description: metadata.VSCapacityCreateMetaData().Description,
			Usage:       metadata.VSCapacityCreateMetaData().Usage,
			Flags:       metadata.VSCapacityCreateMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("VS capacity create", func() {
		Context("VS create with incorrect parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--flavor", "C1_1X1X100", "-c", "1")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "flag provided but not defined: -c")).To(BeTrue())
			})
		})
	})
})