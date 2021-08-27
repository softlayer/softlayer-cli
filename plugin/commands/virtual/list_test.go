package virtual_test

import (
	"errors"
	"strings"

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

var _ = Describe("VS list", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeVSManager *testhelpers.FakeVirtualServerManager
		cmd           *virtual.ListCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		cmd = virtual.NewListCommand(fakeUI, fakeVSManager)
		cliCommand = cli.Command{
			Name:        metadata.VSListMetaData().Name,
			Description: metadata.VSListMetaData().Description,
			Usage:       metadata.VSListMetaData().Usage,
			Flags:       metadata.VSListMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("VS list", func() {
		Context("VS list with wrong parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--hourly", "--monthly")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: '[--hourly]', '[--monthly]' are exclusive.")).To(BeTrue())
			})
		})

		Context("VS list with wrong parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--column", "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: --column abc is not supported.")).To(BeTrue())
			})
		})
		Context("VS list with wrong parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--columns", "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: --columns abc is not supported.")).To(BeTrue())
			})
		})
		Context("VS list with wrong parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: --sortby abc is not supported.")).To(BeTrue())
			})
		})

		Context("VS list with server fails", func() {
			BeforeEach(func() {
				fakeVSManager.ListInstancesReturns([]datatypes.Virtual_Guest{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to list virtual server instances on your account.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
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
				err := testhelpers.RunCommand(cliCommand, "--sortby", "id")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "789")).To(BeTrue())
				Expect(strings.Contains(results[2], "987")).To(BeTrue())
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "hostname")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "abc-vs")).To(BeTrue())
				Expect(strings.Contains(results[2], "vs-abc")).To(BeTrue())
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "domain")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "abc.com")).To(BeTrue())
				Expect(strings.Contains(results[2], "wilma.com")).To(BeTrue())
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "datacenter")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "dal10")).To(BeTrue())
				Expect(strings.Contains(results[2], "tok02")).To(BeTrue())
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "cpu")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "1")).To(BeTrue())
				Expect(strings.Contains(results[2], "4")).To(BeTrue())
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "memory")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "1024")).To(BeTrue())
				Expect(strings.Contains(results[2], "4096")).To(BeTrue())
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "public_ip")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "9.9.8.9")).To(BeTrue())
				Expect(strings.Contains(results[2], "9.9.9.9")).To(BeTrue())
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "private_ip")
				Expect(err).NotTo(HaveOccurred())
				results := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(results[1], "1.1.1.0")).To(BeTrue())
				Expect(strings.Contains(results[2], "1.1.1.1")).To(BeTrue())
			})
		})
	})
})
