package virtual_test

import (
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"strings"
)

var _ = Describe("VS capacity-detail", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeVSManager *testhelpers.FakeVirtualServerManager
		cmd           *virtual.CapacityDetailCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		cmd = virtual.NewCapacityDetailCommand(fakeUI, fakeVSManager)
		cliCommand = cli.Command{
			Name:        virtual.VSCapacityDetailMetaData().Name,
			Description: virtual.VSCapacityDetailMetaData().Description,
			Usage:       virtual.VSCapacityDetailMetaData().Usage,
			Flags:       virtual.VSCapacityDetailMetaData().Flags,
			Action:      cmd.Run,
		}
	})
	Describe("VS capacity-detail", func() {
		Context("Capacity-Detail without vs ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})
		Context("VS capacity-detail with wrong VS ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Reserved Capacity Group Virtual server ID")).To(BeTrue())
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
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
