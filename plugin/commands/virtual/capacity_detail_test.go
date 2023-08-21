package virtual_test

import (
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

var _ = Describe("VS capacity-detail", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *virtual.CapacityDetailCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeVSManager *testhelpers.FakeVirtualServerManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = virtual.NewCapacityDetailCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.VirtualServerManager = fakeVSManager
	})
	Describe("VS capacity-detail", func() {
		Context("Capacity-Detail without vs ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})
		Context("VS capacity-detail with wrong VS ID", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Reserved Capacity Group Virtual server ID"))
			})
		})
		Context("VS capacity detail successfull", func() {
			BeforeEach(func() {
				fakeVSManager.GetCapacityDetailReturns(datatypes.Virtual_ReservedCapacityGroup{
					Id:   sl.Int(123456),
					Name: sl.String("test"),
					Instances: []datatypes.Virtual_ReservedCapacityGroup_Instance{
						datatypes.Virtual_ReservedCapacityGroup_Instance{
							Id: sl.Int(1234567),
							Guest: &datatypes.Virtual_Guest{
								Hostname:                sl.String("unitest"),
								Domain:                  sl.String("techsupport"),
								PrimaryIpAddress:        sl.String("168.192.0.12"),
								PrimaryBackendIpAddress: sl.String("192.168.1.2"),
							},
						}, {},
					},
				}, nil)
			})
			It("return successfully", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
