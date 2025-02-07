package vlan_test

import (
	"errors"
	"strings"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/vlan"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VLAN List", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *vlan.ListCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerCommand
		fakeNetworkManager *testhelpers.FakeNetworkManager
		fakeHandler     *testhelpers.FakeTransportHandler
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession(nil)
		fakeHandler = testhelpers.GetSessionHandler(fakeSession)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = vlan.NewListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		fakeNetworkManager = new(testhelpers.FakeNetworkManager)
		cliCommand.NetworkManager = fakeNetworkManager
	})
	AfterEach(func() {
		// Clear API call logs and any errors that might have been set after every test
		fakeHandler.ClearApiCallLogs()
		fakeHandler.ClearErrors()
	})
	Describe("VLAN list", func() {
		Context("VLAN list but server API call fails", func() {
			BeforeEach(func() {
				fakeNetworkManager.ListVlansReturns([]datatypes.Network_Vlan{}, errors.New("Internal Server Error"))
			})
			It("Account::getNetworkVlans() failed", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to list VLANs on your account."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
			It("Network_Pod::getPods failed", func() {
				fakeNetworkManager.ListVlansReturns([]datatypes.Network_Vlan{}, nil)
				fakeNetworkManager.GetPodsReturns([]datatypes.Network_Pod{}, errors.New("Internal Server Error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get Pods."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("VLAN list with wrong --sortby", func() {
			BeforeEach(func() {
				fakeNetworkManager.ListVlansReturns([]datatypes.Network_Vlan{}, nil)
			})
			It("Incorrect Usage: --sortby", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "abcd")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --sortby abcd is not supported."))
			})
		})

		Context("VLAN list Empty", func() {
			BeforeEach(func() {
				fakeNetworkManager.ListVlansReturns([]datatypes.Network_Vlan{}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{""}))
			})
		})

		Context("Happy Path Tests", func() {

			BeforeEach(func() {
				// Set to a real network manager to get results from test fixtures
				cliCommand.NetworkManager = managers.NewNetworkManager(fakeSession)
			})
			It("--sortby id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "id")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(len(results)).To(Equal(5))
				Expect(results[1]).To(ContainSubstring("12345"))
				Expect(results[2]).To(ContainSubstring("554422"))
				Expect(results[3]).To(ContainSubstring("663322"))
			})

			It("--sortby number", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "number")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(len(results)).To(Equal(5))
				Expect(results[1]).To(ContainSubstring("554422   1111     dal14.bcr03.1111"))
				Expect(results[2]).To(ContainSubstring("663322   2222     dal15.bcr03.2222"))
			})

			It("--sortby name", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "name")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(len(results)).To(Equal(5))
				Expect(results[1]).To(ContainSubstring("554422   1111     dal14.bcr03.1111"))
				Expect(results[2]).To(ContainSubstring("663322   2222     dal15.bcr03.2222"))
			})

			It("--sortby firewall", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "firewall")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("12345    3333     dal13.bcr03.3333"))
				Expect(results[2]).To(ContainSubstring("554422   1111     dal14.bcr03.1111"))
			})

			It("--sortby datacenter", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "datacenter")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("12345    3333     dal13.bcr03.3333"))
				Expect(results[2]).To(ContainSubstring("554422   1111     dal14.bcr03.1111"))
			})

			It("--sortby hardware", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "hardware")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("12345    3333     dal13.bcr03.3333"))
				Expect(results[2]).To(ContainSubstring("663322   2222     dal15.bcr03.2222"))
			})

			It("--sortby virtual_servers", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "virtual_servers")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("12345    3333     dal13.bcr03.3333"))
				Expect(results[2]).To(ContainSubstring("554422   1111     dal14.bcr03.1111"))
			})

			It("--sortby public_ips", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "public_ips")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("12345    3333     dal13.bcr03.3333"))
				Expect(results[2]).To(ContainSubstring("554422   1111     dal14.bcr03.1111"))
			})

			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output", "json")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring(`[`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`{`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`}`))
				Expect(fakeUI.Outputs()).To(ContainSubstring(`]`))
			})
		})
		Context("Issues844", func() {
			BeforeEach(func() {
				fakeSession = testhelpers.NewFakeSoftlayerSession([]string{"getNetworkVlans_844"})
				cliCommand.NetworkManager = managers.NewNetworkManager(fakeSession)
			})
			It("Handle empty Datacenter Name for some routers", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("dal13.fcr01.1362"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("3087"))
			})
		})
	})
})
