package autoscale_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"

	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/autoscale"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("autoscale list", func() {
	var (
		fakeUI               *terminal.FakeUI
		cliCommand           *autoscale.ListCommand
		fakeSession          *session.Session
		slCommand            *metadata.SoftlayerCommand
		fakeAutoScaleManager *testhelpers.FakeAutoScaleManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeAutoScaleManager = new(testhelpers.FakeAutoScaleManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = autoscale.NewListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.AutoScaleManager = fakeAutoScaleManager
	})

	Describe("autoscale list", func() {

		Context("Return error", func() {
			It("Set invalid output", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeAutoScaleManager.ListScaleGroupsReturns([]datatypes.Scale_Group{}, errors.New("Failed to get scale groups."))
			})
			It("Failed get scale groups", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get scale groups."))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerScaleGroups := []datatypes.Scale_Group{
					datatypes.Scale_Group{
						Id:   sl.Int(111111),
						Name: sl.String("scalegroup1"),
						Status: &datatypes.Scale_Group_Status{
							Name: sl.String("Active"),
						},
						MinimumMemberCount:      sl.Int(1),
						MaximumMemberCount:      sl.Int(5),
						VirtualGuestMemberCount: sl.Uint(1),
					},
				}
				fakeAutoScaleManager.ListScaleGroupsReturns(fakerScaleGroups, nil)
			})
			It("List scale groups", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("111111"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("scalegroup1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Active"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1/5"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1"))
			})
		})

	})
})
