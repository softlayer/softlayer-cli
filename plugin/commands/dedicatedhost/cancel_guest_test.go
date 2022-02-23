package dedicatedhost_test

import (
	"errors"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dedicatedhost"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Dedicated host create", func() {
	var (
		fakeUI                   *terminal.FakeUI
		FakeDedicatedhostManager *testhelpers.FakeDedicatedhostManager
		cmd                      *dedicatedhost.CancelCommand
		cliCommand               cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeDedicatedhostManager = new(testhelpers.FakeDedicatedhostManager)
		cmd = dedicatedhost.NewCancelCommand(fakeUI, FakeDedicatedhostManager)
		cliCommand = cli.Command{
			Name:        metadata.DedicatedhostCancelGuestsMetaData().Name,
			Description: metadata.DedicatedhostCancelGuestsMetaData().Description,
			Usage:       metadata.DedicatedhostCancelGuestsMetaData().Usage,
			Flags:       metadata.DedicatedhostCancelGuestsMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Dedicatedhost VS cancel", func() {
		Context("Dedicatedhost VS cancel without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})
		Context("Dedicatedhost VS cancel with wrong VS ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Host ID'. It must be a positive integer."))

			})
		})

		Context("Dedicatedhost VS cancel with correct Host ID but not continue", func() {
			It("return no error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"This will cancel all virtual server instances in the dedicatedhost: 1234 and cannot be undone. Continue?"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Aborted."}))
			})
		})

		Context("Dedicatedhost VS cancel with server fails", func() {
			BeforeEach(func() {
				FakeDedicatedhostManager.CancelGuestsReturns(nil, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to cancel all guests in the dedicatedhost: 1234.\n"))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("Dedicatedhost VS cancel vs successfully", func() {
			BeforeEach(func() {
				FakeDedicatedhostManager.CancelGuestsReturns([]managers.StatusInfo{
					{
						Id:     *sl.Int(1234567),
						Fqdn:   *sl.String("test.softlayer"),
						Status: *sl.String("Cancelled"),
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1234567"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("test.softlayer"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Cancelled"))
			})
		})

	})
})
