package vlan_test

import (
	"errors"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/vlan"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VLAN create", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *vlan.CreateCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
		fakeNetworkManager *testhelpers.FakeNetworkManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = vlan.NewCreateCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cliCommand.NetworkManager = fakeNetworkManager
	})

	Describe("VLAN create", func() {
		Context("VLAN create with -r and -d", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-r", "router123", "-d", "dal09")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: [-r|--router] is not allowed with [-d|--datacenter] or [-t|--vlan-type]."))
				Expect(err.Error()).To(ContainSubstrings([]string{"sl vlan options' to check available options."}))
			})
		})

		Context("VLAN create with -r and -t", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-r", "router123", "-t", "public")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: [-r|--router] is not allowed with [-d|--datacenter] or [-t|--vlan-type]."))
				Expect(err.Error()).To(ContainSubstrings([]string{"sl vlan options' to check available options."}))
			})
		})

		Context("VLAN create with -t but no -d", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-t", "public")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: [-d|--datacenter] and [-t|--vlan-type] are required."))
				Expect(err.Error()).To(ContainSubstrings([]string{"sl vlan options' to check available options."}))
			})
		})

		Context("VLAN create with -d but no -t", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-d", "dal10")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: [-d|--datacenter] and [-t|--vlan-type] are required."))
				Expect(err.Error()).To(ContainSubstrings([]string{"sl vlan options' to check available options."}))
			})
		})

		Context("VLAN create with -d and -t but not continue", func() {
			It("return no error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-t", "public", "-d", "dal10")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"This action will incur charges on your account. Continue?"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Aborted."}))
			})
		})

		Context("VLAN create with correct parameters but server API call fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.AddVlanReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-t", "public", "-d", "dal10", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to add VLAN."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("VLAN create with correct parameters", func() {
			BeforeEach(func() {
				fakeNetworkManager.AddVlanReturns(datatypes.Container_Product_Order_Receipt{OrderId: sl.Int(12345678)}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-t", "public", "-d", "dal10", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The order 12345678 was placed."}))
			})
		})
	})
})
