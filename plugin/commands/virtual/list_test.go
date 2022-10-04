package virtual_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS list", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *virtual.ListCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeVSManager *testhelpers.FakeVirtualServerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = virtual.NewListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.VirtualServerManager = fakeVSManager
	})

	Describe("VS list", func() {
		Context("VS list with wrong parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--hourly", "--monthly")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '[--hourly]', '[--monthly]' are exclusive."))
			})
		})

		Context("VS list with wrong parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--column", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --column abc is not supported."))
			})
		})
		Context("VS list with wrong parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --sortby abc is not supported."))
			})
		})

		Context("VS list with server fails", func() {
			BeforeEach(func() {
				fakeVSManager.ListInstancesReturns([]datatypes.Virtual_Guest{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to list virtual server instances on your account."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("VS list with different --sortby", func() {
			BeforeEach(func() {
				fakeVSManager.ListInstancesReturns([]datatypes.Virtual_Guest{
					datatypes.Virtual_Guest{
						Id:       sl.Int(987),
						Hostname: sl.String("vs-abc"),
						Domain:   sl.String("wilma.com"),
						Datacenter: &datatypes.Location{
							Name: sl.String("tok02"),
						},
						MaxCpu:                  sl.Int(4),
						MaxMemory:               sl.Int(4096),
						PrimaryIpAddress:        sl.String("9.9.9.9"),
						PrimaryBackendIpAddress: sl.String("1.1.1.1"),
					},
					datatypes.Virtual_Guest{
						Id:       sl.Int(789),
						Hostname: sl.String("abc-vs"),
						Domain:   sl.String("abc.com"),
						Datacenter: &datatypes.Location{
							Name: sl.String("dal10"),
						},
						MaxCpu:                  sl.Int(1),
						MaxMemory:               sl.Int(1024),
						PrimaryIpAddress:        sl.String("9.9.8.9"),
						PrimaryBackendIpAddress: sl.String("1.1.1.0"),
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "id")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("789"))
				Expect(results[2]).To(ContainSubstring("987"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "hostname")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("abc-vs"))
				Expect(results[2]).To(ContainSubstring("vs-abc"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "domain")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("abc.com"))
				Expect(results[2]).To(ContainSubstring("wilma.com"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "datacenter")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("dal10"))
				Expect(results[2]).To(ContainSubstring("tok02"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "cpu")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("1"))
				Expect(results[2]).To(ContainSubstring("4"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "memory")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("1024"))
				Expect(results[2]).To(ContainSubstring("4096"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "public_ip")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("9.9.8.9"))
				Expect(results[2]).To(ContainSubstring("9.9.9.9"))
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "private_ip")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(results[1]).To(ContainSubstring("1.1.1.0"))
				Expect(results[2]).To(ContainSubstring("1.1.1.1"))
			})
		})
	})
})
