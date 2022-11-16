package autoscale_test

import (
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/autoscale"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("autoscale detail", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *autoscale.DetailCommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = autoscale.NewDetailCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})

	Describe("autoscale detail", func() {

		Context("Return error", func() {
			It("Set command without Id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})

			It("Set command with an invalid Id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Autoscale Group ID'. It must be a positive integer."))
			})

			It("Set invalid output", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Autoscale detail, correct use", func() {
			It("return autoscale detail", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "12222222")
				Expect(err).NotTo(HaveOccurred())
				output := strings.Split(fakeUI.Outputs(), "\n")
				Expect(output[0]).To(ContainSubstring("Name"))
				Expect(output[0]).To(ContainSubstring("Value"))
				Expect(output[1]).To(ContainSubstring("ID"))
				Expect(output[1]).To(ContainSubstring("12222222"))
				Expect(output[7]).To(ContainSubstring("Cooldown"))
				Expect(output[7]).To(ContainSubstring("1800 seconds"))
				Expect(output[8]).To(ContainSubstring("Last Action"))
				Expect(output[8]).To(ContainSubstring("2019-10-02T20:26:17Z"))
				Expect(output[12]).To(ContainSubstring("Virtual Guest Member Template"))
				Expect(output[12]).To(ContainSubstring("Name"))
				Expect(output[12]).To(ContainSubstring("Value"))
				Expect(output[13]).To(ContainSubstring("Hostname"))
				Expect(output[13]).To(ContainSubstring("testing"))
				Expect(output[14]).To(ContainSubstring("Domain"))
				Expect(output[14]).To(ContainSubstring("tech-support.com"))
			})
			It("return autoscale detail in format json", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Name"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Value"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"ID",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"12222222"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Minimum Members",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"2"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Maximum Members",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"6"`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"SSH Key 490279   label   "`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"SAN Disk 0       25   "`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"SAN Disk 2       10   "`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Post Install     https://test.com/   "`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Active Guests",`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
	})
})
