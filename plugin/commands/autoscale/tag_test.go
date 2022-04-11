package autoscale_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/autoscale"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("autoscale logs", func() {
	var (
		fakeUI                   *terminal.FakeUI
		fakeAutoScaleManager     *testhelpers.FakeAutoScaleManager
		fakeVirtualServerManager *testhelpers.FakeVirtualServerManager
		cmd                      *autoscale.TagCommand
		cliCommand               cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeAutoScaleManager = new(testhelpers.FakeAutoScaleManager)
		fakeVirtualServerManager = new(testhelpers.FakeVirtualServerManager)
		cmd = autoscale.NewTagCommand(fakeUI, fakeAutoScaleManager, fakeVirtualServerManager)
		cliCommand = cli.Command{
			Name:        autoscale.AutoScaleTagMetaData().Name,
			Description: autoscale.AutoScaleTagMetaData().Description,
			Usage:       autoscale.AutoScaleTagMetaData().Usage,
			Flags:       autoscale.AutoScaleTagMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("autoscale tag", func() {

		Context("Return error", func() {
			It("Set command without Id", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one identifier."))
			})

			It("Set command with an invalid Id", func() {
				err := testhelpers.RunCommand(cliCommand, "abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Autoscale group ID should be a number."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeAutoScaleManager.GetVirtualGuestMembersReturns([]datatypes.Scale_Member_Virtual_Guest{}, errors.New("Failed to get virtual guest members."))
			})
			It("Failed get scale group virtual guests", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get virtual guest members."))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerMembers := []datatypes.Scale_Member_Virtual_Guest{
					datatypes.Scale_Member_Virtual_Guest{
						VirtualGuest: &datatypes.Virtual_Guest{
							Id:       sl.Int(111111),
							Hostname: sl.String("myhostname"),
						},
					},
				}
				fakeAutoScaleManager.GetVirtualGuestMembersReturns(fakerMembers, nil)
				fakeVirtualServerManager.SetTagsReturns(nil)
			})

			It("Set Tags", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--tags=mytag1")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("New Tags"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Setting tags"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Done"))
			})

			It("Clear Tags", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("All tags will be removed"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Setting tags"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Done"))
			})
		})
	})
})
