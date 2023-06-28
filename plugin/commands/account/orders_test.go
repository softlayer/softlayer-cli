package account_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/account"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Account list Orders", func() {
	var (
		fakeUI      *terminal.FakeUI
		cliCommand  *account.OrdersCommand
		fakeSession *session.Session
		slCommand   *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = account.NewOrdersCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
	})

	Describe("Account orders", func() {
		Context("Account orders, Invalid Usage", func() {
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
			It("Set command with an invalid limit option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--limit", "abcd")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`invalid argument "abcd"`))
			})
		})

		Context("Account orders, correct use", func() {
			It("return account orders", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--limit", "10")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Id          State      User                       Date                   Amount     Item"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456789   APPROVED   test.test@ibm.com          2022-04-26T19:50:06Z   0.000000   1 x 2.0 GHz or higher Core,1 x 2.0 GHz or higher Core,1 GB,Reboot / Remote Conso..."))
				Expect(fakeUI.Outputs()).To(ContainSubstring("91954410    APPROVED   123456_test.test@ibm.com   2022-04-26T19:39:17Z   0.000000   1 x 2.0 GHz or higher Core,1 x 2.0 GHz or higher Core,1 GB,Reboot / Remote Conso..."))
			})
			It("return account orders in format json", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Orders":`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Id": "123456789","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"State": "APPROVED","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"User": "test.test@ibm.com","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Amount": "0.000000","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Item": "1 x 2.0 GHz or higher Core,1 x 2.0 GHz or higher Core,1 GB,Reboot / Remote Conso...""`))
			})
		})

		Context("Account orders, correct use with upgrades", func() {
			It("return account orders and upgrades", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--limit", "10", "--upgrades")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456789"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("APPROVED"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("test.test@ibm.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2022-04-26T19:50:06Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("0.000000"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1 x 2.0 GHz or higher Core,1 x 2.0 GHz or higher Core,1 GB,Reboot / Remote Conso..."))
				Expect(fakeUI.Outputs()).To(ContainSubstring("91954410"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456_test.test@ibm.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2022-04-26T19:39:17Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("3237486"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("3"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("152510472"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("3229998"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("152452042"))
			})
		})
	})
})
