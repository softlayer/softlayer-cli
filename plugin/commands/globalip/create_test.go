package globalip_test

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
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/globalip"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("GlobalIP create", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeNetworkManager *testhelpers.FakeNetworkManager
		cmd                *globalip.CreateCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cmd = globalip.NewCreateCommand(fakeUI, fakeNetworkManager)
		cliCommand = cli.Command{
			Name:        metadata.GlobalIpCreateMetaData().Name,
			Description: metadata.GlobalIpCreateMetaData().Description,
			Usage:       metadata.GlobalIpCreateMetaData().Usage,
			Flags:       metadata.GlobalIpCreateMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("GlobalIP create", func() {
		Context("GlobalIP create without -f", func() {
			It("return no error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCommand(cliCommand, "")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"This action will incur charges on your account. Continue?"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Aborted."}))
			})
		})

		Context("GlobalIP create with -test", func() {
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--test")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The order is correct."}))
			})
		})

		Context("GlobalIP create with -test but server fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.AddGlobalIPReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--test")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to add global IP.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("GlobalIP create ", func() {
			BeforeEach(func() {
				fakeNetworkManager.AddGlobalIPReturns(datatypes.Container_Product_Order_Receipt{
					OrderId: sl.Int(12345678),
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Order 12345678 was placed."}))
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--v6", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Order 12345678 was placed."}))
			})
		})
	})
})
