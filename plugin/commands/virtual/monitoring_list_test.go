package virtual_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"

	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS monitoring list", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeVSManager *testhelpers.FakeVirtualServerManager
		cmd           *virtual.MonitoringListCommand
		cliCommand    cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		cmd = virtual.NewMonitoringListCommand(fakeUI, fakeVSManager)
		cliCommand = cli.Command{
			Name:        virtual.VSMonitoringListMetaData().Name,
			Description: virtual.VSMonitoringListMetaData().Description,
			Usage:       virtual.VSMonitoringListMetaData().Usage,
			Flags:       virtual.VSMonitoringListMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("VS monitoring list", func() {
		Context("Return error", func() {
			It("Set command without id", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})

			It("Set command with an invalid id", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Virtual server ID'. It must be a positive integer."))
			})

			It("Set command with an invalid output format", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeVSManager.GetInstanceReturns(datatypes.Virtual_Guest{}, errors.New("Internal Server Error"))
			})
			It("Command fails to get VS", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get virtual server"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerVS := datatypes.Virtual_Guest{
					Id:                      sl.Int(123456),
					Domain:                  sl.String("domain.com"),
					PrimaryIpAddress:        sl.String("9.9.9.9"),
					PrimaryBackendIpAddress: sl.String("1.1.1.1"),
					Datacenter: &datatypes.Location{
						LongName: sl.String("Dallas 10"),
					},
					NetworkMonitors: []datatypes.Network_Monitor_Version1_Query_Host{
						datatypes.Network_Monitor_Version1_Query_Host{
							Id:        sl.Int(678),
							IpAddress: sl.String("2.2.2.2"),
							Status:    sl.String("ON"),
							QueryType: &datatypes.Network_Monitor_Version1_Query_Type{
								Name: sl.String("SERVICE PING"),
							},
							ResponseAction: &datatypes.Network_Monitor_Version1_Query_ResponseType{
								ActionDescription: sl.String("Do Nothing"),
							},
						},
					},
				}
				fakeVSManager.GetInstanceReturns(fakerVS, nil)
			})
			It("Set command with correct virtual server id", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("domain.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("9.9.9.9"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1.1.1.1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Dallas 10"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("678"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2.2.2.2"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("ON"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("SERVICE PING"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Do Nothing"))
			})
		})
	})
})
