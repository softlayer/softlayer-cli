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

var _ = Describe("autoscale scale", func() {
	var (
		fakeUI               *terminal.FakeUI
		fakeAutoScaleManager *testhelpers.FakeAutoScaleManager
		cmd                  *autoscale.ScaleCommand
		cliCommand           cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeAutoScaleManager = new(testhelpers.FakeAutoScaleManager)
		cmd = autoscale.NewScaleCommand(fakeUI, fakeAutoScaleManager)
		cliCommand = cli.Command{
			Name:        autoscale.AutoScaleScaleMetaData().Name,
			Description: autoscale.AutoScaleScaleMetaData().Description,
			Usage:       autoscale.AutoScaleScaleMetaData().Usage,
			Flags:       autoscale.AutoScaleScaleMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("autoscale scale", func() {

		Context("Return error", func() {
			It("Set command without Id", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})

			It("Set command with an invalid Id", func() {
				err := testhelpers.RunCommand(cliCommand, "abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Autoscale Group ID'. It must be a positive integer."))
			})

			It("Set command without --amount", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '--amount' is required"))
			})

			It("Set command without --by or --to", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--amount=1")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --to or --by is required"))
			})

			It("Set command with --by and --to", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--amount=1", "--by", "--to")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '[--to]', '[--by]' are exclusive"))
			})

			It("Set command with --up and --down", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--amount=1", "--by", "--up", "--down")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '[--up]', '[--down]' are exclusive"))
			})
		})

		Context("Scale method returns error", func() {
			BeforeEach(func() {
				fakeAutoScaleManager.ScaleReturns([]datatypes.Scale_Member{}, errors.New("Failed to scale Auto Scale Group."))
			})
			It("Failed get scale group", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--amount=1", "--by", "--down")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to scale Auto Scale Group."))
			})
		})

		Context("GetVirtualGuestMembers method returns error", func() {
			BeforeEach(func() {
				fakeAutoScaleManager.GetVirtualGuestMembersReturns([]datatypes.Scale_Member_Virtual_Guest{}, errors.New("Failed to get virtual guest members."))
			})
			It("Failed get scale group", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--amount=1", "--to")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get virtual guest members."))
			})
		})

		Context("ScaleTo method returns error", func() {
			BeforeEach(func() {
				fakeAutoScaleManager.ScaleToReturns([]datatypes.Scale_Member{}, errors.New("Failed to scale Auto Scale Group."))
			})
			It("Failed get scale group", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--amount=1", "--to")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to scale Auto Scale Group."))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerNewMembers := []datatypes.Scale_Member{
					datatypes.Scale_Member{
						Id: sl.Int(111111),
					},
				}
				fakeAutoScaleManager.ScaleReturns(fakerNewMembers, nil)
			})

			It("Set command with --by and --down", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--amount=1", "--by", "--down")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Auto Scale Group was scaled successfully"))
			})
		})
	})
})
