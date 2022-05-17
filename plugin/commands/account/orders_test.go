package account_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/session"

	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/account"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Account list Orders", func() {
	var (
		fakeUI      *terminal.FakeUI
		cmd         *account.OrdersCommand
		cliCommand  cli.Command
		fakeSession *session.Session
		fakeAccountManager managers.AccountManager
	)
	BeforeEach(func() {
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeAccountManager = managers.NewAccountManager(fakeSession)
		fakeUI = terminal.NewFakeUI()
		cmd = account.NewOrdersCommand(fakeUI, fakeAccountManager)
		cliCommand = cli.Command{
			Name:        account.OrdersMetaData().Name,
			Description: account.OrdersMetaData().Description,
			Usage:       account.OrdersMetaData().Usage,
			Flags:       account.OrdersMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Account orders", func() {
		Context("Account orders, Invalid Usage", func() {
			It("Set command with an invalid output option", func() {
				err := testhelpers.RunCommand(cliCommand, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
			It("Set command with an invalid date option", func() {
				err := testhelpers.RunCommand(cliCommand, "--limit", "abcd")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`invalid value "abcd" for flag -limit: parse error`))
			})
		})

		Context("Account orders, correct use", func() {
			It("return account orders", func() {
				err := testhelpers.RunCommand(cliCommand, "--limit", "10")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Orders"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Id          State      User                       Date                   Amount     Item"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456789   APPROVED   test.test@ibm.com          2022-04-26T19:50:06Z   0.000000   1 x 2.0 GHz or higher Core,1 x 2.0 GHz or higher Core,1 GB,Reboot / Remote Conso..."))
				Expect(fakeUI.Outputs()).To(ContainSubstring("91954410    APPROVED   123456_test.test@ibm.com   2022-04-26T19:39:17Z   0.000000   1 x 2.0 GHz or higher Core,1 x 2.0 GHz or higher Core,1 GB,Reboot / Remote Conso..."))
			})
			It("return account orders in format json", func() {
				err := testhelpers.RunCommand(cliCommand, "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Orders":`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Id": "123456789","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"State": "APPROVED","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"User": "test.test@ibm.com","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Amount": "0.000000","`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`"Item": "1 x 2.0 GHz or higher Core,1 x 2.0 GHz or higher Core,1 GB,Reboot / Remote Conso...""`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
	})
})
