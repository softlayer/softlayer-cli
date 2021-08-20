package virtual_test

import (

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS bandwidth", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeVSManager *testhelpers.FakeVirtualServerManager
		cmd           *virtual.BandwidthCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		cmd = virtual.NewBandwidthCommand(fakeUI, fakeVSManager)
		cliCommand = cli.Command{
			Name:        metadata.VSBandwidthMetaData().Name,
			Description: metadata.VSBandwidthMetaData().Description,
			Usage:       metadata.VSBandwidthMetaData().Usage,
			Flags:       metadata.VSBandwidthMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("VS bandwidth", func() {
		Context("Argument Checking", func() {
			It("Error on missing ID", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})
		Context("DateTime parsing checks", func() {
			It("YYYY-MM-DD Parsing works properly", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "-s", "2021-08-01", "-e", "2021-08-10")
				Expect(err).NotTo(HaveOccurred())

			})
		})
		
	})
})
