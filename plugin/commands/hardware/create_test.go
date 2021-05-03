package hardware_test

import (
	"errors"
	"fmt"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/cgallo/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/cgallo/softlayer-cli/plugin/metadata"
	"github.ibm.com/cgallo/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("hardware create", func() {
	var (
		fakeUI              *terminal.FakeUI
		fakeHardwareManager *testhelpers.FakeHardwareServerManager
		cmd                 *hardware.CreateCommand
		cliCommand          cli.Command
		context             plugin.PluginContext
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeHardwareManager = new(testhelpers.FakeHardwareServerManager)
		context = plugin.InitPluginContext("softlayer")
		cmd = hardware.NewCreateCommand(fakeUI, fakeHardwareManager, context)
		cliCommand = cli.Command{
			Name:        metadata.HardwareCreateMetaData().Name,
			Description: metadata.HardwareCreateMetaData().Description,
			Usage:       metadata.HardwareCreateMetaData().Usage,
			Flags:       metadata.HardwareCreateMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("hardware create", func() {
		Context("hardware create with non-exist template file", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-m", "/tmp/template.txt")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Template file: /tmp/template.txt does not exist."))
			})
		})
		Context("hardware create with no size", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-s|--size' is required"))
			})
		})
		Context("hardware create with no hostname", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-s", "S1270_32GB_2X960GBSSD_NORAID")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-H|--hostname' is required"))
			})
		})
		Context("hardware create with no domain", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-D|--domain' is required"))
			})
		})
		Context("hardware create with no osName", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-o|--os' is required"))
			})
		})
		Context("hardware create with no datacenter", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-d|--datacenter' is required"))
			})
		})
		Context("hardware create with no port speed", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '-p|--port-speed' is required"))
			})
		})
		Context("hardware create with wrong billing", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: -b|--billing has to be either hourly or monthly."))
			})
		})
		Context("hardware create with get package fails", func() {
			BeforeEach(func() {
				fakeHardwareManager.GetPackageReturns(datatypes.Product_Package{}, errors.New("Internal server error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "hourly", "-t")
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
				err := testhelpers.RunCommand(cliCommand, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "hourly", "-t")
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
				err := testhelpers.RunCommand(cliCommand, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "hourly", "-t")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Total monthly cost"))
			})
		})
		Context("hardware create with export succeed", func() {
			It("return file", func() {
				err := testhelpers.RunCommand(cliCommand, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "hourly", "-x", "/tmp/template.json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Hardware server template is exported to: /tmp/template.json."))
			})
		})
		Context("hardware create with cancel", func() {
			It("return no error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCommand(cliCommand, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "hourly")
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
				err := testhelpers.RunCommand(cliCommand, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "hourly")
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
				err := testhelpers.RunCommand(cliCommand, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "hourly")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 123456 was placed."))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{fmt.Sprintf("Run '%s sl hardware list --order 123456' to find this hardware server after it is ready.", cmd.Context.CLIName())}))
			})
			It("return order receipt", func() {
				err := testhelpers.RunCommand(cliCommand, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "hourly", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 123456 was placed."))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{fmt.Sprintf("Run '%s sl hardware list --order 123456' to find this hardware server after it is ready.", cmd.Context.CLIName())}))
			})
			It("return order receipt", func() {
				err := testhelpers.RunCommand(cliCommand, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "monthly", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 123456 was placed."))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{fmt.Sprintf("Run '%s sl hardware list --order 123456' to find this hardware server after it is ready.", cmd.Context.CLIName())}))
			})
			It("return order receipt", func() {
				err := testhelpers.RunCommand(cliCommand, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "monthly", "-i", "https://postinstall.sh", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 123456 was placed."))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{fmt.Sprintf("Run '%s sl hardware list --order 123456' to find this hardware server after it is ready.", cmd.Context.CLIName())}))
			})
			It("return order receipt", func() {
				err := testhelpers.RunCommand(cliCommand, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "monthly", "-i", "https://postinstall.sh", "-n", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 123456 was placed."))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{fmt.Sprintf("Run '%s sl hardware list --order 123456' to find this hardware server after it is ready.", cmd.Context.CLIName())}))
			})
			It("return order receipt", func() {
				err := testhelpers.RunCommand(cliCommand, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "monthly", "-i", "https://postinstall.sh", "-n", "-k", "123", "-k", "234", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 123456 was placed."))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{fmt.Sprintf("Run '%s sl hardware list --order 123456' to find this hardware server after it is ready.", cmd.Context.CLIName())}))
			})
			It("return order receipt", func() {
				err := testhelpers.RunCommand(cliCommand, "-s", "S1270_32GB_2X960GBSSD_NORAID", "-H", "ibmcloud-cli", "-D", "ibm.com", "-o", "UBUNTU_16_64", "-d", "dal10", "-p", "1000", "-b", "monthly", "-i", "https://postinstall.sh", "-n", "-k", "123", "-k", "234", "-e", "1_IPV6_ADDRESS", "-e", "64_BLOCK_STATIC_PUBLIC_IPV6_ADDRESSES", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("Order 123456 was placed."))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{fmt.Sprintf("Run '%s sl hardware list --order 123456' to find this hardware server after it is ready.", cmd.Context.CLIName())}))
			})
		})
	})
})
