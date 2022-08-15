package reports_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/softlayer/softlayer-go/session"


	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/reports"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
)

var _ = Describe("reports bandwidth", func() {
    var (
        fakeUI          *terminal.FakeUI
        cliCommand      *reports.BandwidthCommand
        fakeSession     *session.Session
        slCommand       *metadata.SoftlayerCommand
        fakeReportManager *testhelpers.FakeReportManager
    )
    BeforeEach(func() {
        fakeUI = terminal.NewFakeUI()
        fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
        fakeReportManager = new(testhelpers.FakeReportManager)
        slCommand  = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
        cliCommand = reports.NewBandwidthCommand(slCommand)
        cliCommand.ReportManager = fakeReportManager
        cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")

    })

	Describe("reports bandwidth", func() {

		Context("Return error", func() {
			It("Set invalid output", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})

			It("Set invalid --sortby option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby=id")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid --sortBy option."))
			})

			It("Set invalid --start option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--start=20220305")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid format date to --start."))
			})

			It("Set invalid --end option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--end=20220305")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid format date to --end."))
			})

			It("--end is not greater than --start", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--start=2022-04-05", "--end=2022-03-05")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: End Date must be greater than Start Date."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeReportManager.GetVirtualGuestsReturns([]datatypes.Virtual_Guest{}, errors.New("Failed to get virtual guests on your account."))
			})
			It("Failed get virtual guests with --virtual option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--virtual")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get virtual guests on your account."))
			})

			It("Failed get virtual guests without --virtual option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get virtual guests on your account."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeReportManager.GetHardwareServersReturns([]datatypes.Hardware{}, errors.New("Failed to get hardware servers on your account."))
			})
			It("Failed get hardware servers with --server option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--server")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get hardware servers on your account."))
			})

			It("Failed get hardware servers without --server option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get hardware servers on your account."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeReportManager.GetVirtualDedicatedRacksReturns([]datatypes.Network_Bandwidth_Version1_Allotment{}, errors.New("Failed to get virtual dedicated racks on your account."))
			})
			It("Failed get pools with --pool option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--pool")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get virtual dedicated racks on your account."))
			})

			It("Failed get pools without --pool option", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get virtual dedicated racks on your account."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakerVirtualGuests := []datatypes.Virtual_Guest{
					datatypes.Virtual_Guest{
						Id:                     sl.Int(111111),
						Hostname:               sl.String("virtualGuestHostname"),
						MetricTrackingObjectId: sl.Int(222222),
						VirtualRack: &datatypes.Network_Bandwidth_Version1_Allotment{
							BandwidthAllotmentTypeId: sl.Int(2),
							Name:                     sl.String("virtualGuestPool"),
						},
					},
				}
				fakeReportManager.GetVirtualGuestsReturns(fakerVirtualGuests, nil)
				fakeReportManager.GetMetricTrackingSummaryDataReturns([]datatypes.Metric_Tracking_Object_Data{}, errors.New("Failed to get metric tracking summary"))
			})
			It("Failed to get virtual guest metric tracking summary", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--virtual")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakerVirtualDedicatedRacks := []datatypes.Network_Bandwidth_Version1_Allotment{
					datatypes.Network_Bandwidth_Version1_Allotment{
						Id:                     sl.Int(111111),
						Name:                   sl.String("poolName"),
						MetricTrackingObjectId: sl.Int(222222),
					},
				}
				fakeReportManager.GetVirtualDedicatedRacksReturns(fakerVirtualDedicatedRacks, nil)
				fakeReportManager.GetMetricTrackingSummaryDataReturns([]datatypes.Metric_Tracking_Object_Data{}, errors.New("Failed to get metric tracking summary"))
			})
			It("Failed to get pool metric tracking summary", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--pool")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerVirtualGuests := []datatypes.Virtual_Guest{
					datatypes.Virtual_Guest{
						Id:                     sl.Int(111111),
						Hostname:               sl.String("virtualGuestHostname"),
						MetricTrackingObjectId: sl.Int(222222),
						VirtualRack: &datatypes.Network_Bandwidth_Version1_Allotment{
							BandwidthAllotmentTypeId: sl.Int(2),
							Name:                     sl.String("virtualGuestPool"),
						},
					},
				}
				fakeMetricTrackingSummaryData := []datatypes.Metric_Tracking_Object_Data{
					datatypes.Metric_Tracking_Object_Data{
						Counter: sl.Float(1000000000),
						Type:    sl.String("publicIn_net_octet"),
					},
					datatypes.Metric_Tracking_Object_Data{
						Counter: sl.Float(2000000000),
						Type:    sl.String("publicOut_net_octet"),
					},
					datatypes.Metric_Tracking_Object_Data{
						Counter: sl.Float(3000000000),
						Type:    sl.String("privateIn_net_octet"),
					},
					datatypes.Metric_Tracking_Object_Data{
						Counter: sl.Float(4000000000),
						Type:    sl.String("privateOut_net_octet"),
					},
				}
				fakeReportManager.GetVirtualGuestsReturns(fakerVirtualGuests, nil)
				fakeReportManager.GetMetricTrackingSummaryDataReturns(fakeMetricTrackingSummaryData, nil)
			})
			It("Display virtual guests bandwidth summary with diferents --sortby options", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--virtual", "--sortby=type")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("virtual"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("virtualGuestHostname"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("3.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("4.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("virtualGuestPool"))
			})

			It("Display virtual guests bandwidth summary", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--virtual", "--sortby=hostname")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("virtual"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("virtualGuestHostname"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("3.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("4.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("virtualGuestPool"))
			})

			It("Display virtual guests bandwidth summary", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--virtual", "--sortby=publicIn")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("virtual"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("virtualGuestHostname"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("3.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("4.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("virtualGuestPool"))
			})

			It("Display virtual guests bandwidth summary", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--virtual", "--sortby=publicOut")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("virtual"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("virtualGuestHostname"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("3.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("4.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("virtualGuestPool"))
			})

			It("Display virtual guests bandwidth summary", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--virtual", "--sortby=privateIn")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("virtual"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("virtualGuestHostname"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("3.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("4.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("virtualGuestPool"))
			})

			It("Display virtual guests bandwidth summary", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--virtual", "--sortby=privateOut")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("virtual"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("virtualGuestHostname"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("3.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("4.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("virtualGuestPool"))
			})

			It("Display virtual guests bandwidth summary", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--virtual", "--sortby=pool")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("virtual"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("virtualGuestHostname"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("3.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("4.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("virtualGuestPool"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				fakerVirtualDedicatedRacks := []datatypes.Network_Bandwidth_Version1_Allotment{
					datatypes.Network_Bandwidth_Version1_Allotment{
						Id:                     sl.Int(111111),
						Name:                   sl.String("poolName"),
						MetricTrackingObjectId: sl.Int(222222),
					},
				}
				fakeMetricTrackingSummaryData := []datatypes.Metric_Tracking_Object_Data{
					datatypes.Metric_Tracking_Object_Data{
						Counter: sl.Float(1000000000),
						Type:    sl.String("publicIn_net_octet"),
					},
					datatypes.Metric_Tracking_Object_Data{
						Counter: sl.Float(2000000000),
						Type:    sl.String("publicOut_net_octet"),
					},
					datatypes.Metric_Tracking_Object_Data{
						Counter: sl.Float(3000000000),
						Type:    sl.String("privateIn_net_octet"),
					},
					datatypes.Metric_Tracking_Object_Data{
						Counter: sl.Float(4000000000),
						Type:    sl.String("privateOut_net_octet"),
					},
				}
				fakeReportManager.GetVirtualDedicatedRacksReturns(fakerVirtualDedicatedRacks, nil)
				fakeReportManager.GetMetricTrackingSummaryDataReturns(fakeMetricTrackingSummaryData, nil)
			})
			It("Display pool bandwidth summary", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--pool")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("pool"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("poolName"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("3.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("4.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("-"))
			})
		})
	})
})
