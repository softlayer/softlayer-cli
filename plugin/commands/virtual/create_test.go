package virtual_test

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/plugin"
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

var _ = Describe("VS create", func() {
	var (
		fakeUI           *terminal.FakeUI
		fakeVSManager    *testhelpers.FakeVirtualServerManager
		fakeImageManager *testhelpers.FakeImageManager
		cmd              *virtual.CreateCommand
		cliCommand       cli.Command
		context          plugin.PluginContext
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		fakeImageManager = new(testhelpers.FakeImageManager)
		context = plugin.InitPluginContext("softlayer")
		cmd = virtual.NewCreateCommand(fakeUI, fakeVSManager, fakeImageManager, context)
		cliCommand = cli.Command{
			Name:        metadata.VSCreateMataData().Name,
			Description: metadata.VSCreateMataData().Description,
			Usage:       metadata.VSCreateMataData().Usage,
			Flags:       metadata.VSCreateMataData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("VS create", func() {
		Context("VS create with incorrect parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--flavor", "C1_1X1X100", "-c", "1")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: '[-c|--cpu]', '[--flavor]' are exclusive.")).To(BeTrue())
			})
		})
		Context("VS create with incorrect parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--flavor", "C1_1X1X100", "-m", "1024")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: '[-m|--memory]', '[--flavor]' are exclusive.")).To(BeTrue())
			})
		})
		Context("VS create with incorrect parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--flavor", "C1_1X1X100", "--dedicated")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: '[--dedicated]', '[--flavor]' are exclusive.")).To(BeTrue())
			})
		})
		Context("VS create with incorrect parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--flavor", "C1_1X1X100", "--host-id", "12345")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: '[--host-id]', '[--flavor]' are exclusive.")).To(BeTrue())
			})
		})
		Context("VS create with incorrect parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-o", "CENTOS", "--image", "111")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: '[-o|--os]', '[--image]' are exclusive.")).To(BeTrue())
			})
		})
		Context("VS create with incorrect parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--billing", "yearly")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [--billing] billing rate must be either hourly or monthly.")).To(BeTrue())
			})
		})
		Context("VS create with incorrect parameters", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-u", "CENTOS", "-F", "/tmp/file")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: '[-u|--userdata]', '[-F|--userfile]' are exclusive.")).To(BeTrue())
			})
		})
		Context("VS create with incorrect parameters", func() {
			BeforeEach(func() {
				fakeVSManager.GenerateInstanceCreationTemplateReturns(&datatypes.Virtual_Guest{}, errors.New("Incorrect Usage: [-H|--hostname] is required."))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-o", "CENTOS")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [-H|--hostname] is required.")).To(BeTrue())
			})
		})
		Context("VS create with incorrect parameters", func() {
			BeforeEach(func() {
				fakeVSManager.GenerateInstanceCreationTemplateReturns(&datatypes.Virtual_Guest{}, errors.New("Incorrect Usage: [-D|--domain] is required."))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-H", "vs-abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [-D|--domain] is required.")).To(BeTrue())
			})
		})
		Context("VS create with incorrect parameters", func() {
			BeforeEach(func() {
				fakeVSManager.GenerateInstanceCreationTemplateReturns(&datatypes.Virtual_Guest{}, errors.New("Incorrect Usage: [-c|--cpu] is required and must be positive integer."))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-H", "vs-abc", "-D", "wilma.com")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [-c|--cpu] is required and must be positive integer.")).To(BeTrue())
			})
		})
		Context("VS create with incorrect parameters", func() {
			BeforeEach(func() {
				fakeVSManager.GenerateInstanceCreationTemplateReturns(&datatypes.Virtual_Guest{}, errors.New("Incorrect Usage: [-m|--memory] is required and must be positive integer."))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-H", "vs-abc", "-D", "wilma.com", "-c", "2")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [-m|--memory] is required and must be positive integer.")).To(BeTrue())
			})
		})
		Context("VS create with incorrect parameters", func() {
			BeforeEach(func() {
				fakeVSManager.GenerateInstanceCreationTemplateReturns(&datatypes.Virtual_Guest{}, errors.New("Incorrect Usage: [--datacenter] is required."))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-H", "vs-abc", "-D", "wilma.com", "-c", "2", "-m", "4096")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: [--datacenter] is required.")).To(BeTrue())
			})
		})
		Context("VS create with incorrect parameters", func() {
			BeforeEach(func() {
				fakeVSManager.GenerateInstanceCreationTemplateReturns(&datatypes.Virtual_Guest{}, errors.New("Incorrect Usage: either [-o|--os] or [--image] is required."))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-H", "vs-abc", "-D", "wilma.com", "-c", "2", "-m", "4096", "--datacenter", "dal10")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: either [-o|--os] or [--image] is required.")).To(BeTrue())
			})
		})
		Context("VS create with incorrect parameters", func() {
			BeforeEach(func() {
				fakeVSManager.GenerateInstanceCreationTemplateReturns(&datatypes.Virtual_Guest{}, errors.New("Incorrect Usage: Template file: /abc/def/tmplate does not exist."))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-H", "vs-abc", "-D", "wilma.com", "-c", "2", "-m", "4096", "--datacenter", "dal10", "-o", "CENTOS", "--template", "/abc/def/tmplate")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: Template file: /abc/def/tmplate does not exist.")).To(BeTrue())
			})
		})
		Context("VS create with incorrect parameters", func() {
			BeforeEach(func() {
				fakeVSManager.GenerateInstanceCreationTemplateReturns(&datatypes.Virtual_Guest{}, errors.New("Incorrect Usage: Local disk number cannot excceed two."))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-H", "vs-abc", "-D", "wilma.com", "-c", "2", "-m", "4096", "--datacenter", "dal10", "-o", "CENTOS", "disk", "100", "disk", "100", "disk", "100")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: Local disk number cannot excceed two.")).To(BeTrue())
			})
		})
		Context("VS create with incorrect parameters", func() {
			BeforeEach(func() {
				fakeVSManager.GenerateInstanceCreationTemplateReturns(&datatypes.Virtual_Guest{}, errors.New("Incorrect Usage: San disk number cannot excceed five."))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-H", "vs-abc", "-D", "wilma.com", "-c", "2", "-m", "4096", "--datacenter", "dal10", "-o", "CENTOS", "--san", "disk", "100", "disk", "100", "disk", "100", "disk", "100", "disk", "100", "disk", "100")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: San disk number cannot excceed five.")).To(BeTrue())
			})
		})

		Context("VS create with --export fails", func() {
			BeforeEach(func() {
				fakeVSManager.GenerateInstanceCreationTemplateReturns(&datatypes.Virtual_Guest{}, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-H", "vs-abc", "-D", "wilma.com", "-c", "2", "-m", "4096", "--datacenter", "dal10", "-o", "CENTOS", "--export", "/root/template")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to write virtual server template file to: /root/template.")).To(BeTrue())
			})
		})
		Context("VS create with --export succeed", func() {
			BeforeEach(func() {
				fakeVSManager.GenerateInstanceCreationTemplateReturns(&datatypes.Virtual_Guest{}, nil)
				
			})
			It("return no error", func() {
				tmpFile, tmpErr := ioutil.TempFile(os.TempDir(), "create_tests-")
				if tmpErr != nil {
					Skip("Cannot create temporary file")
				}
				fileName := tmpFile.Name()
				defer os.Remove(fileName)
				err := testhelpers.RunCommand(cliCommand, "-H", "vs-abc", "-D", "wilma.com", "-c", "2", "-m", "4096", "--datacenter", "dal10", "-o", "CENTOS", "--export", fileName)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Virtual server template is exported to: " + fileName + "."}))
			})
		})

		Context("VS create with --test fails", func() {
			BeforeEach(func() {
				fakeVSManager.GenerateInstanceCreationTemplateReturns(&datatypes.Virtual_Guest{}, nil)
				fakeVSManager.VerifyInstanceCreationReturns(datatypes.Container_Product_Order{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-H", "vs-abc", "-D", "wilma.com", "-c", "2", "-m", "4096", "--datacenter", "dal10", "-o", "CENTOS", "--test")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to verify virtual server creation.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})
		Context("VS create with --test succeed", func() {
			BeforeEach(func() {
				fakeVSManager.GenerateInstanceCreationTemplateReturns(&datatypes.Virtual_Guest{}, nil)
				fakeVSManager.VerifyInstanceCreationReturns(datatypes.Container_Product_Order{}, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-H", "vs-abc", "-D", "wilma.com", "-c", "2", "-m", "4096", "--datacenter", "dal10", "-o", "CENTOS", "--test")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"OK"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"The order is correct."}))
			})
		})
		Context("Check for -s and -S from #3489", func() {
			BeforeEach(func() {
				fakeVSManager.GenerateInstanceCreationTemplateReturns(&datatypes.Virtual_Guest{}, nil)
				fakeVSManager.VerifyInstanceCreationReturns(datatypes.Container_Product_Order{}, nil)
			})
			It("Make sure -S sets public-security-group", func() {
				err1 := testhelpers.RunCommand(cliCommand, "--hostname", "vs111abc", "-D", "wilma.com", "--flavor", "C1_1X1X100", "--datacenter", "dal10", "-o", "CENTOS", "-S", "9999", "--test")
				err2 := testhelpers.RunCommand(cliCommand, "--hostname", "vs111abc", "-D", "wilma.com", "--flavor", "C1_1X1X100", "--datacenter", "dal10", "-o", "CENTOS", "--public-security-group", "9999", "--test")
				Expect(err1).NotTo(HaveOccurred())
				Expect(err2).NotTo(HaveOccurred())
				_, call1 := fakeVSManager.GenerateInstanceCreationTemplateArgsForCall(0)
				_, call2 := fakeVSManager.GenerateInstanceCreationTemplateArgsForCall(1)
				Expect(call1["public-security-group"]).To(Equal([]int{9999}))
				Expect(call2["public-security-group"]).To(Equal([]int{9999}))
			})
			It("Make sure -s sets private-security-group", func() {
				err1 := testhelpers.RunCommand(cliCommand, "--hostname", "vs111abc", "-D", "wilma.com", "--flavor", "C1_1X1X100", "--datacenter", "dal10", "-o", "CENTOS", "-s", "9999", "--test")
				err2 := testhelpers.RunCommand(cliCommand, "--hostname", "vs111abc", "-D", "wilma.com", "--flavor", "C1_1X1X100", "--datacenter", "dal10", "-o", "CENTOS", "--private-security-group", "9999", "--test")
				Expect(err1).NotTo(HaveOccurred())
				Expect(err2).NotTo(HaveOccurred())
				_, call1 := fakeVSManager.GenerateInstanceCreationTemplateArgsForCall(0)
				_, call2 := fakeVSManager.GenerateInstanceCreationTemplateArgsForCall(1)
				Expect(call1["private-security-group"]).To(Equal([]int{9999}))
				Expect(call2["private-security-group"]).To(Equal([]int{9999}))
			})
		})

		Context("VS create without -f and not continue", func() {
			BeforeEach(func() {
				fakeVSManager.GenerateInstanceCreationTemplateReturns(&datatypes.Virtual_Guest{}, nil)
				fakeVSManager.VerifyInstanceCreationReturns(datatypes.Container_Product_Order{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCommand(cliCommand, "-H", "vs-abc", "-D", "wilma.com", "-c", "2", "-m", "4096", "--datacenter", "dal10", "-o", "CENTOS")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"This action will incur charges on your account. Continue?"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"Aborted."}))
			})
		})

		Context("VS create with -f but server fails", func() {
			BeforeEach(func() {
				fakeVSManager.GenerateInstanceCreationTemplateReturns(&datatypes.Virtual_Guest{}, nil)
				fakeVSManager.VerifyInstanceCreationReturns(datatypes.Container_Product_Order{}, errors.New("Internal Server Error"))
				fakeVSManager.CreateInstanceReturns(datatypes.Virtual_Guest{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-H", "vs-abc", "-D", "wilma.com", "-c", "2", "-m", "4096", "--datacenter", "dal10", "-o", "CENTOS", "-f")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to create virtual server instance.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Internal Server Error")).To(BeTrue())
			})
		})

		Context("VS create with -f and succeed", func() {
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2017-01-03T00:00:00Z")
				fakeVSManager.GenerateInstanceCreationTemplateReturns(&datatypes.Virtual_Guest{}, nil)
				fakeVSManager.VerifyInstanceCreationReturns(datatypes.Container_Product_Order{}, errors.New("Internal Server Error"))
				fakeVSManager.CreateInstanceReturns(datatypes.Virtual_Guest{
					Id:                       sl.Int(1234),
					FullyQualifiedDomainName: sl.String("vs-abc.wilma.com"),
					GlobalIdentifier:         sl.String("dhtyengodyhebt"),
					CreateDate:               sl.Time(created),
				}, nil)
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-H", "vs-abc", "-D", "wilma.com", "-c", "2", "-m", "4096", "--datacenter", "dal10", "-o", "CENTOS", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1234"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"vs-abc.wilma.com"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2017-01-03T00:00:00Z"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"dhtyengodyhebt"}))
			})
		})

		Context("VS create with succeed but set tag fails", func() {
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2017-01-03T00:00:00Z")
				fakeVSManager.GenerateInstanceCreationTemplateReturns(&datatypes.Virtual_Guest{}, nil)
				fakeVSManager.VerifyInstanceCreationReturns(datatypes.Container_Product_Order{}, nil)
				fakeVSManager.CreateInstanceReturns(datatypes.Virtual_Guest{
					Id:                       sl.Int(1234),
					FullyQualifiedDomainName: sl.String("vs-abc.wilma.com"),
					GlobalIdentifier:         sl.String("dhtyengodyhebt"),
					CreateDate:               sl.Time(created),
				}, nil)
				fakeVSManager.SetTagsReturns(errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-H", "vs-abc", "-D", "wilma.com", "-c", "2", "-m", "4096", "--datacenter", "dal10", "-o", "CENTOS", "-f", "--tag", "mytag")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1234"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"vs-abc.wilma.com"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2017-01-03T00:00:00Z"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"dhtyengodyhebt"}))
				Expect(strings.Contains(err.Error(), "Failed to update the tag of virtual server instance: 1234."))
				Expect(strings.Contains(err.Error(), "Internal Server Error"))
			})
		})

		Context("VS create with succeed but get ready fails", func() {
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2017-01-03T00:00:00Z")
				fakeVSManager.GenerateInstanceCreationTemplateReturns(&datatypes.Virtual_Guest{}, nil)
				fakeVSManager.VerifyInstanceCreationReturns(datatypes.Container_Product_Order{}, nil)
				fakeVSManager.CreateInstanceReturns(datatypes.Virtual_Guest{
					Id:                       sl.Int(1234),
					FullyQualifiedDomainName: sl.String("vs-abc.wilma.com"),
					GlobalIdentifier:         sl.String("dhtyengodyhebt"),
					CreateDate:               sl.Time(created),
				}, nil)
				fakeVSManager.InstanceIsReadyReturns(false, "", errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "-H", "vs-abc", "-D", "wilma.com", "-c", "2", "-m", "4096", "--datacenter", "dal10", "-o", "CENTOS", "-f", "--wait", "1")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1234"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"vs-abc.wilma.com"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2017-01-03T00:00:00Z"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"dhtyengodyhebt"}))
				Expect(strings.Contains(err.Error(), "Failed to get ready status of virtual server instance: 1234."))
				Expect(strings.Contains(err.Error(), "Internal Server Error"))
			})
		})

		Context("VS create  succeed and get ready succeed", func() {
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2017-01-03T00:00:00Z")
				fakeVSManager.GenerateInstanceCreationTemplateReturns(&datatypes.Virtual_Guest{}, nil)
				fakeVSManager.VerifyInstanceCreationReturns(datatypes.Container_Product_Order{}, nil)
				fakeVSManager.CreateInstanceReturns(datatypes.Virtual_Guest{
					Id:                       sl.Int(1234),
					FullyQualifiedDomainName: sl.String("vs-abc.wilma.com"),
					GlobalIdentifier:         sl.String("dhtyengodyhebt"),
					CreateDate:               sl.Time(created),
				}, nil)
				fakeVSManager.InstanceIsReadyReturns(true, "", nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "-H", "vs-abc", "-D", "wilma.com", "-c", "2", "-m", "4096", "--datacenter", "dal10", "-o", "CENTOS", "-f", "--wait", "1")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1234"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"vs-abc.wilma.com"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2017-01-03T00:00:00Z"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"dhtyengodyhebt"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"true"}))
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "-H", "vs-abc", "-D", "wilma.com", "--flavor", "M1_1X8X100", "--datacenter", "dal10", "-o", "CENTOS", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1234"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"vs-abc.wilma.com"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2017-01-03T00:00:00Z"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"dhtyengodyhebt"}))
			})
		})
	})
})
