package virtual_test

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS create", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *virtual.CreateCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeVSManager *testhelpers.FakeVirtualServerManager
		fakeServer    datatypes.Virtual_Guest
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = virtual.NewCreateCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.VirtualServerManager = fakeVSManager
		// These APIS are generally called for vs create
		created, _ := time.Parse(time.RFC3339, "2017-01-03T00:00:00Z")
		fakeServer = datatypes.Virtual_Guest{
			Id:                       sl.Int(1234),
			FullyQualifiedDomainName: sl.String("vs-abc.wilma.com"),
			GlobalIdentifier:         sl.String("dhtyengodyhebt"),
			CreateDate:               sl.Time(created),
		}
		fakeVSManager.GenerateInstanceCreationTemplateReturns(&datatypes.Virtual_Guest{}, nil)
		fakeVSManager.VerifyInstanceCreationReturns(datatypes.Container_Product_Order{}, nil)
		fakeVSManager.CreateInstanceReturns(fakeServer, nil)
		fakeVSManager.InstanceIsReadyReturns(true, "", nil)
	})

	Describe("VS create", func() {
		Context("VS create with incorrect parameters", func() {
			BeforeEach(func() {
				fakeUI.Inputs("yes", "")
			})
			It("Flavor/CPU Conflict", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--flavor", "C1_1X1X100", "-c", "1")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '[-c|--cpu]', '[--flavor]' are exclusive."))
			})
			It("Flavor/Memory Conflict", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--flavor", "C1_1X1X100", "-m", "1024")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '[-m|--memory]', '[--flavor]' are exclusive."))
			})
			It("Flavor/Dedicated Conflict", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--flavor", "C1_1X1X100", "--dedicated")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '[--dedicated]', '[--flavor]' are exclusive."))
			})
			It("Flavor/Host Conflict", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--flavor", "C1_1X1X100", "--host-id", "12345")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '[--host-id]', '[--flavor]' are exclusive."))
			})
			It("OS/Image Conflict", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-o", "CENTOS", "--image", "111")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '[-o|--os]', '[--image]' are exclusive."))
			})
			It("Billing Conflict", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--billing", "yearly")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: [--billing] billing rate must be either hourly or monthly."))
			})
			It("UserData/UserFile Conflict", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-u", "CENTOS", "-F", "/tmp/file")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: '[-u|--userdata]', '[-F|--userfile]' are exclusive."))
			})
			It("SAN/Local Conflict", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--san", "--local")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("[local san] were all set"))
			})
		})
		Context("VS create with --export", func() {
			It("Success", func() {
				tmpFile, tmpErr := ioutil.TempFile(os.TempDir(), "create_tests-")
				if tmpErr != nil {
					Skip("Cannot create temporary file")
				}
				fileName := tmpFile.Name()
				defer os.Remove(fileName)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-H", "vs-abc", "-D", "wilma.com", "-c", "2", "-m", "4096", "--datacenter", "dal10", "-o", "CENTOS", "--export", fileName)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Virtual server template is exported to: " + fileName + "."))
			})
		})
		Context("VS create with --test", func() {
			It("API Error", func() {
				fakeVSManager.VerifyInstanceCreationReturns(datatypes.Container_Product_Order{}, errors.New("Internal Server Error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-H", "vs-abc", "-D", "wilma.com", "-c", "2", "-m", "4096", "--datacenter", "dal10", "-o", "CENTOS", "--test")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to verify virtual server creation."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
			It("Success", func() {
				fakeVSManager.VerifyInstanceCreationReturns(datatypes.Container_Product_Order{}, nil)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-H", "vs-abc", "-D", "wilma.com", "-c", "2", "-m", "4096", "--datacenter", "dal10", "-o", "CENTOS", "--test")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("OK"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("The order is correct."))
			})
		})
		Context("Check for -s and -S", func() {
			It("Make sure -S sets public-security-group", func() {
				err1 := testhelpers.RunCobraCommand(cliCommand.Command, "--hostname", "vs111abc", "-D", "wilma1.com", "--flavor", "C1_1X1X100", "--datacenter", "dal10", "-o", "CENTOS", "-S", "9999", "--test")
				Expect(err1).NotTo(HaveOccurred())
				_, call1 := fakeVSManager.GenerateInstanceCreationTemplateArgsForCall(0)
				Expect(call1["public-security-group"]).To(Equal([]int{9999}))

			})
			It("Make sure -s sets private-security-group", func() {
				err1 := testhelpers.RunCobraCommand(cliCommand.Command, "--hostname", "vs111abc", "-D", "wilma.com", "--flavor", "C1_1X1X100", "--datacenter", "dal10", "-o", "CENTOS", "-s", "9999", "--test")
				Expect(err1).NotTo(HaveOccurred())
				_, call1 := fakeVSManager.GenerateInstanceCreationTemplateArgsForCall(0)
				Expect(call1["private-security-group"]).To(Equal([]int{9999}))
			})
		})
		Context("VS create -f", func() {
			It("Aborted", func() {
				fakeUI.Inputs("No")
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-H", "vs-abc", "-D", "wilma.com", "-c", "2", "-m", "4096", "--datacenter", "dal10", "-o", "CENTOS")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("This action will incur charges on your account. Continue?"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Aborted."))
			})
			It("API Error", func() {
				fakeVSManager.CreateInstanceReturns(datatypes.Virtual_Guest{}, errors.New("Internal Server Error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-H", "vs-abc", "-D", "wilma.com", "-c", "2", "-m", "4096", "--datacenter", "dal10", "-o", "CENTOS", "-f")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to create virtual server instance."))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
			It("Success", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-H", "vs-abc", "-D", "wilma.com", "-c", "2", "-m", "4096", "--datacenter", "dal10", "-o", "CENTOS", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1234"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("vs-abc.wilma.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2017-01-03T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dhtyengodyhebt"))
			})
		})
		Context("VS create with tags", func() {
			It("return error", func() {
				fakeVSManager.SetTagsReturns(errors.New("Internal Server Error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-H", "vs-abc", "-D", "wilma.com", "-c", "2", "-m", "4096", "--datacenter", "dal10", "-o", "CENTOS", "-f", "--tag", "mytag")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1234"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("vs-abc.wilma.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2017-01-03T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dhtyengodyhebt"))
				Expect(strings.Contains(err.Error(), "Failed to update the tag of virtual server instance: 1234."))
				Expect(strings.Contains(err.Error(), "Internal Server Error"))
			})
		})
		Context("VS create with succeed but get ready fails", func() {
			It("Read API Error", func() {
				fakeVSManager.InstanceIsReadyReturns(false, "", errors.New("Internal Server Error"))
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-H", "vs-abc", "-D", "wilma.com", "-c", "2", "-m", "4096", "--datacenter", "dal10", "-o", "CENTOS", "-f", "--wait", "1")
				Expect(err).To(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1234"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("vs-abc.wilma.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2017-01-03T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dhtyengodyhebt"))
				Expect(strings.Contains(err.Error(), "Failed to get ready status of virtual server instance: 1234."))
				Expect(strings.Contains(err.Error(), "Internal Server Error"))
			})
		})
		Context("Happy Path", func() {
			It("Created with ready check", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-H", "vs-abc", "-D", "wilma.com", "-c", "2", "-m", "4096", "--datacenter", "dal10", "-o", "CENTOS", "-f", "--wait", "1")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1234"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("vs-abc.wilma.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2017-01-03T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dhtyengodyhebt"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("true"))
			})
			It("Created", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "-H", "vs-abc", "-D", "wilma.com", "--flavor", "M1_1X8X100", "--datacenter", "dal10", "-o", "CENTOS", "-f")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("1234"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("vs-abc.wilma.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2017-01-03T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("dhtyengodyhebt"))
			})
			It("Test setting local flag", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--local", "-f", "-H=vs-abc", "-D=wilma.com", "--flavor=SOME_FLAVOR", "--datacenter=dal10", "-o=CENTOS")
				Expect(err).NotTo(HaveOccurred())
				_, vsTemplate := fakeVSManager.GenerateInstanceCreationTemplateArgsForCall(0)
				Expect(vsTemplate["san"]).To(BeFalse())
			})
			It("Test setting san flag", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--san", "-f", "-H=vs-abc", "-D=wilma.com", "--flavor=SOME_FLAVOR", "--datacenter=dal10", "-o=CENTOS")
				Expect(err).NotTo(HaveOccurred())
				_, vsTemplate := fakeVSManager.GenerateInstanceCreationTemplateArgsForCall(0)
				Expect(vsTemplate["san"]).To(BeTrue())
			})
		})
	})
})
