package account_test


import (

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"


	"github.com/softlayer/softlayer-go/session"

	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/account"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Account Bandwidth-Pools", func() {
	var (
		fakeUI			*terminal.FakeUI
		// fakeManager		*testhelpers.FakeAccountManager
		cmd				*account.BandwidthPoolsCommand
		cliCommand		cli.Command
		fakeSession   	*session.Session
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		// fakeManager = new(testhelpers.FakeAccountManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		cmd = account.NewBandwidthPoolsCommand(fakeUI, fakeSession)
		cliCommand = cli.Command{
			Name:	account.BandwidthPoolsMetaData().Name,
			Description: account.BandwidthPoolsMetaData().Description,
			Usage:	account.BandwidthPoolsMetaData().Usage,
			Flags:	account.BandwidthPoolsMetaData().Flags,
			Action:	cmd.Run,
		}
	})

	Describe("Bandwidth-Pools Testing", func() {
		Context("Happy Path", func() {
			It("Runs without issue", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).NotTo(HaveOccurred())
				outputs := fakeUI.Outputs()
				Expect(outputs).To(ContainSubstring("3361 GB      7.13 GB         7.70 GB"))
			})
		})
	})
})
