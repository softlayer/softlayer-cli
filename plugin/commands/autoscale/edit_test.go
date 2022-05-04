package autoscale_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/autoscale"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("autoscale edit", func() {
	var (
		fakeUI               *terminal.FakeUI
		fakeAutoScaleManager *testhelpers.FakeAutoScaleManager
		cmd                  *autoscale.EditCommand
		cliCommand           cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeAutoScaleManager = new(testhelpers.FakeAutoScaleManager)
		cmd = autoscale.NewEditCommand(fakeUI, fakeAutoScaleManager)
		cliCommand = cli.Command{
			Name:        autoscale.AutoScaleEditMetaData().Name,
			Description: autoscale.AutoScaleEditMetaData().Description,
			Usage:       autoscale.AutoScaleEditMetaData().Usage,
			Flags:       autoscale.AutoScaleEditMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("autoscale edit", func() {

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

			It("Set command without options", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Please pass at least one of the flags."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeAutoScaleManager.GetScaleGroupReturns(datatypes.Scale_Group{}, errors.New("Failed to get AutoScale group."))
			})
			It("Failed get scale group", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--cpu=2")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get AutoScale group."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeAutoScaleManager.EditScaleGroupReturns(false, errors.New("Failed to update Auto Scale Group."))
			})

			It("Failed update scale group name", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--name=scale2")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to update Auto Scale Group."))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakeAutoScaleManager.EditScaleGroupReturns(true, nil)
			})

			It("Set command to edit name, minimum member count and maximum member count", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--name=newname", "--min=1", "--max=1")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerScaleGroup := datatypes.Scale_Group{
					VirtualGuestMemberTemplate: &datatypes.Virtual_Guest{},
				}
				fakeAutoScaleManager.GetScaleGroupReturns(fakerScaleGroup, nil)
				fakeAutoScaleManager.EditScaleGroupReturns(true, nil)
			})

			It("Set command to edit cpu´s, memory and user data", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--cpu=1", "--memory=1024", "--userdata=CENTOS")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
			})
		})
	})
})
