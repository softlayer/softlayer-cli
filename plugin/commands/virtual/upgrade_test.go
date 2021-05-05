package virtual_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS upgrade", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeVSManager *testhelpers.FakeVirtualServerManager
		cmd           *virtual.UpgradeCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		cmd = virtual.NewUpgradeCommand(fakeUI, fakeVSManager)
		cliCommand = cli.Command{
			Name:        metadata.VSUpgradeMataData().Name,
			Description: metadata.VSUpgradeMataData().Description,
			Usage:       metadata.VSUpgradeMataData().Usage,
			Flags:       metadata.VSUpgradeMataData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("VS upgrade", func() {
		Context("VS upgrade without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})
		Context("VS upgrade with wrong vs ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Virtual server ID'. It must be a positive integer.")).To(BeTrue())
			})
		})
		Context("VS upgrade with wrong parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--private")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: Must specify [--cpu] when using [--private].")).To(BeTrue())
			})
		})
		Context("VS upgrade with wrong parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: Must provide [--cpu], [--memory], [--network] or [--flavor] to upgrade.")).To(BeTrue())
			})
		})
		Context("VS upgrade without -f", func() {
			It("return no error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCommand(cliCommand, "1234", "--cpu", "8")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"This action will incur charges on your account. Continue?"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Aborted"}))
			})
		})
		Context("VS upgrade with server fails", func() {
			BeforeEach(func() {
				fakeVSManager.UpgradeInstanceReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--cpu", "8", "-f")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to upgrade virtual server instance: 1234.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
		Context("VS upgrade", func() {
			BeforeEach(func() {
				fakeVSManager.UpgradeInstanceReturns(datatypes.Container_Product_Order_Receipt{
					OrderId: sl.Int(12345678),
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--cpu", "8", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Order 12345678 to upgrade virtual server instance: 1234 was placed."}))
			})
		})
	})
})
