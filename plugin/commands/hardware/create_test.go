package hardware_test

import (
	"errors"
	"os"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("hardware create", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeHardwareManager *testhelpers.FakeHardwareServerManager
		cliCommand          *hardware.CreateCommand
		fakeSession         *session.Session
		slCommand           *metadata.SoftlayerCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeHardwareManager = new(testhelpers.FakeHardwareServerManager)
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = hardware.NewCreateCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.HardwareManager = fakeHardwareManager
	})

	Describe("hardware create", func() {
		Context("hardware create with non-exist template file", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-m", "/tmp/template.txt")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Template file: /tmp/template.txt does not exist."))
			})
		})
		Context("hardware create with no size", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-s|--size' is required"))
			})
		})
		Context("hardware create with no hostname", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-s", "S1270_32GB_2X960GBSSD_NORAID")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-H|--hostname' is required"))
			})
		})
		Context("hardware create with no domain", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-D|--domain' is required"))
			})
		})
		Context("hardware create with no osName", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-o|--os' is required"))
			})
		})
		Context("hardware create with no datacenter", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-d|--datacenter' is required"))
			})
		})
		Context("hardware create with no port speed", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-p|--port-speed' is required"))
			})
		})
		Context("hardware create with wrong billing", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: -b|--billing has to be either hourly or monthly."))
			})
		})
		Context("hardware create with get package fails", func() {
			BeforeEach(func() {
				fakeHardwareManager.GetPackageReturns(datatypes.Product_Package{}, errors.New("Internal server error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "hourly", "-t")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get product package for hardware server."))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
		})
		Context("hardware create with verify order fails", func() {
			BeforeEach(func() {
				fakeHardwareManager.VerifyOrderReturns(datatypes.Container_Product_Order{}, errors.New("Internal server error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "hourly", "-t")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to verify this order."))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
		})
		Context("hardware create with verify order succeed", func() {
			BeforeEach(func() {
				fakeHardwareManager.VerifyOrderReturns(datatypes.Container_Product_Order{}, nil)
			})
			It("return order", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "hourly", "-t")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Total monthly cost"))
			})
		})
		Context("hardware create with export succeed", func() {
			It("return file", func() {
				if os.Getenv("OS") == "Windows_NT" {
					Skip("Test doesn't work in windows.")
				}
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "hourly", "-x", "/tmp/template.json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Hardware server template is exported to: /tmp/template.json."))
			})
		})
		Context("hardware create with cancel", func() {
			It("return no error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "hourly")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This action will incur charges on your account. Continue?"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted."))
			})
		})
		Context("hardware create with place order fails", func() {
			BeforeEach(func() {
				fakeHardwareManager.PlaceOrderReturns(datatypes.Container_Product_Order_Receipt{}, errors.New("Internal server error"))
			})
			It("return error", func() {
				fakeUI.Inputs("Yes")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "hourly")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to place this order."))
				Expect(err.Error()).To(ContainSubstring("Internal server error"))
			})
		})
		Context("hardware create with place order succeed", func() {
			BeforeEach(func() {
				fakeHardwareManager.PlaceOrderReturns(datatypes.Container_Product_Order_Receipt{
					OrderId: sl.Int(123456),
				}, nil)
			})
			It("return order receipt", func() {
				fakeUI.Inputs("Yes")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "hourly")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 123456 was placed."))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Run 'ibmcloud sl hardware list --order 123456' to find this hardware server after it is ready."))
			})
			It("return order receipt", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "hourly", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 123456 was placed."))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Run 'ibmcloud sl hardware list --order 123456' to find this hardware server after it is ready."))
			})
			It("return order receipt", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "monthly", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 123456 was placed."))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Run 'ibmcloud sl hardware list --order 123456' to find this hardware server after it is ready."))
			})
			It("return order receipt", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "monthly", "-i", "https://postinstall.sh", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 123456 was placed."))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Run 'ibmcloud sl hardware list --order 123456' to find this hardware server after it is ready."))
			})
			It("return order receipt", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "monthly", "-i", "https://postinstall.sh", "-n", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 123456 was placed."))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Run 'ibmcloud sl hardware list --order 123456' to find this hardware server after it is ready."))
			})
			It("Success with SSH keys", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "monthly", "-i", "https://postinstall.sh", "-n", "-k", "123", "-k", "234", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 123456 was placed."))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Run 'ibmcloud sl hardware list --order 123456' to find this hardware server after it is ready."))
				_, callData := fakeHardwareManager.GenerateCreateTemplateArgsForCall(0)
				sshKeys := callData["sshKeys"]
				// Expect(len(callData.SshKeys[0].SshKeyIds)).To(Equal(1))
				Expect(sshKeys).To(Equal([]int{123, 234}))
				
			})
			It("return order receipt", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "monthly", "-i", "https://postinstall.sh", "-n", "-k", "123", "-k", "234", "-e", "1_IPV6_ADDRESS", "-e", "64_BLOCK_STATIC_PUBLIC_IPV6_ADDRESSES", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 123456 was placed."))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Run 'ibmcloud sl hardware list --order 123456' to find this hardware server after it is ready."))
			})
		})
	})
})
