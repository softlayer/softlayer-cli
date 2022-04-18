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

var _ = Describe("autoscale list", func() {
	var (
		fakeUI               *terminal.FakeUI
		fakeAutoScaleManager *testhelpers.FakeAutoScaleManager
		fakeSecurityManager  *testhelpers.FakeSecurityManager
		cmd                  *autoscale.DetailCommand
		cliCommand           cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeAutoScaleManager = new(testhelpers.FakeAutoScaleManager)
		fakeSecurityManager = new(testhelpers.FakeSecurityManager)
		cmd = autoscale.NewDetailCommand(fakeUI, fakeAutoScaleManager, fakeSecurityManager)
		cliCommand = cli.Command{
			Name:        autoscale.AutoScaleDetailMetaData().Name,
			Description: autoscale.AutoScaleDetailMetaData().Description,
			Usage:       autoscale.AutoScaleDetailMetaData().Usage,
			Flags:       autoscale.AutoScaleDetailMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("autoscale detail", func() {

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

			It("Set invalid output", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeAutoScaleManager.GetScaleGroupReturns(datatypes.Scale_Group{}, errors.New("Failed to get scale group."))
			})
			It("Failed get scale group", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get scale group."))
			})
		})
	})
})
