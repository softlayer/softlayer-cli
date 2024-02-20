package account_test

import (
	"errors"
	"time"

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

var _ = Describe("account hook-create", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *account.HookCreateCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
		fakeAccountManager *testhelpers.FakeAccountManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeAccountManager = new(testhelpers.FakeAccountManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = account.NewHookCreateCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.AccountManager = fakeAccountManager
	})

	Describe("account hook-create", func() {

		Context("Return error", func() {

			It("Set invalid output", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--name=myhook", "--uri=http://myuritest.com", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeAccountManager.CreateProvisioningScriptReturns(datatypes.Provisioning_Hook{}, errors.New("Failed to create Provisioning Hook"))
			})
			It("Failed to create Provisioning Hook", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--name=myhook", "--uri=http://myuritest.com")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to create Provisioning Hook"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2017-11-08T00:00:00Z")
				fakerHook := datatypes.Provisioning_Hook{
					Id:         sl.Int(123456),
					Name:       sl.String("My Hook"),
					CreateDate: sl.Time(created),
					Uri:        sl.String("http://myuritest.com"),
				}
				fakeAccountManager.CreateProvisioningScriptReturns(fakerHook, nil)
			})
			It("Return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--name=myhook", "--uri=http://myuritest.com")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("My Hook"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2017-11-08T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("http://myuritest.com"))
			})
		})
	})
})
