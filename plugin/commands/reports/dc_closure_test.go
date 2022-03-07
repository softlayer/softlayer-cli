package reports_test


import (

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"


	"github.com/softlayer/softlayer-go/session"

	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/reports"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Reports Datacenter-Closures", func() {
	var (
		fakeUI			*terminal.FakeUI
		cmd				*reports.DCClosuresCommand
		cliCommand		cli.Command
		fakeSession   	*session.Session
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		cmd = reports.NewDCClosuresCommand(fakeUI, fakeSession)
		cliCommand = cli.Command{
			Name:	reports.DCClosuresMetaData().Name,
			Description: reports.DCClosuresMetaData().Description,
			Usage:	reports.DCClosuresMetaData().Usage,
			Flags:	reports.DCClosuresMetaData().Flags,
			Action:	cmd.Run,
		}
	})

	Describe("Datacenter-Closures Testing", func() {
		Context("Happy Path", func() {
			It("Runs without issue", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).NotTo(HaveOccurred())
				outputs := fakeUI.Outputs()
				Expect(outputs).To(ContainSubstring("Hello World"))
			})
			It("Outputs JSON", func() {
				err := testhelpers.RunCommand(cliCommand, "--output=JSON")
				Expect(err).NotTo(HaveOccurred())
				outputs := fakeUI.Outputs()
				Expect(outputs).To(ContainSubstring("\"amountIn\": 7.54252,"))
	
			})
		})
	})
})
