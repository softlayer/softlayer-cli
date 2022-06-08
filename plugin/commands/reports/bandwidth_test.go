package reports_test

import (
	"errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/reports"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("reports bandwidth", func() {
	var (
		fakeUI            *terminal.FakeUI
		fakeReportManager *testhelpers.FakeReportManager
		cmd               *reports.BandwidthCommand
		cliCommand        cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeReportManager = new(testhelpers.FakeReportManager)
		cmd = reports.NewBandwidthCommand(fakeUI, fakeReportManager)
		cliCommand = cli.Command{
			Name:        reports.ReportBandwidthMetaData().Name,
			Description: reports.ReportBandwidthMetaData().Description,
			Usage:       reports.ReportBandwidthMetaData().Usage,
			Flags:       reports.ReportBandwidthMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("reports bandwidth", func() {

		Context("Return error", func() {
			It("Set invalid output", func() {
				err := testhelpers.RunCommand(cliCommand, "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})

			It("Set invalid --sortby option", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby=Type")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid --sortBy option."))
			})

			It("Set invalid --start option", func() {
				err := testhelpers.RunCommand(cliCommand, "--start=20220305")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid format date to --start."))
			})

			It("Set invalid --end option", func() {
				err := testhelpers.RunCommand(cliCommand, "--end=20220305")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid format date to --end."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeReportManager.GetVirtualGuestsReturns([]datatypes.Virtual_Guest{}, errors.New("Failed to get virtual guests on your account."))
			})
			It("Failed get virtual guests with --virtual option", func() {
				err := testhelpers.RunCommand(cliCommand, "--virtual")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get virtual guests on your account."))
			})

			It("Failed get virtual guests without --virtual option", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get virtual guests on your account."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeReportManager.GetHardwareServersReturns([]datatypes.Hardware{}, errors.New("Failed to get hardware servers on your account."))
			})
			It("Failed get hardware servers with --server option", func() {
				err := testhelpers.RunCommand(cliCommand, "--server")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get hardware servers on your account."))
			})

			It("Failed get hardware servers without --server option", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get hardware servers on your account."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeReportManager.GetVirtualDedicatedRacksReturns([]datatypes.Network_Bandwidth_Version1_Allotment{}, errors.New("Failed to get virtual dedicated racks on your account."))
			})
			It("Failed get pools with --pool option", func() {
				err := testhelpers.RunCommand(cliCommand, "--pool")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get virtual dedicated racks on your account."))
			})

			It("Failed get pools without --pool option", func() {
				err := testhelpers.RunCommand(cliCommand)
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
				err := testhelpers.RunCommand(cliCommand, "--virtual")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get metric tracking summary"))
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
				err := testhelpers.RunCommand(cliCommand, "--pool")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get metric tracking summary"))
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
				err := testhelpers.RunCommand(cliCommand, "--virtual", "--sortby=type")
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
				err := testhelpers.RunCommand(cliCommand, "--virtual", "--sortby=hostname")
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
				err := testhelpers.RunCommand(cliCommand, "--virtual", "--sortby=publicIn")
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
				err := testhelpers.RunCommand(cliCommand, "--virtual", "--sortby=publicOut")
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
				err := testhelpers.RunCommand(cliCommand, "--virtual", "--sortby=privateIn")
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
				err := testhelpers.RunCommand(cliCommand, "--virtual", "--sortby=privateOut")
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
				err := testhelpers.RunCommand(cliCommand, "--virtual", "--sortby=pool")
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
				err := testhelpers.RunCommand(cliCommand, "--pool")
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
