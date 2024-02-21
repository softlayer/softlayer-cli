package account_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/account"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("account hooks", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *account.HooksCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
		fakeAccountManager *testhelpers.FakeAccountManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeAccountManager = new(testhelpers.FakeAccountManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = account.NewHooksCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.AccountManager = fakeAccountManager
	})

	Describe("account hooks", func() {

		Context("Return error", func() {

			It("Set invalid output", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeAccountManager.GetPostProvisioningHooksReturns([]datatypes.Provisioning_Hook{}, errors.New("Failed to get Provisioning Hooks."))
			})
			It("Failed get Failed to get Provisioning Hooks", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get Provisioning Hooks."))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerHooks := []datatypes.Provisioning_Hook{
					datatypes.Provisioning_Hook{
						Id:   sl.Int(123456),
						Name: sl.String("My Hook"),
						Uri:  sl.String("http://myuritest.com"),
					},
				}
				fakeAccountManager.GetPostProvisioningHooksReturns(fakerHooks, nil)
			})
			It("Return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("My Hook"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("http://myuritest.com"))
			})
		})
	})
})
