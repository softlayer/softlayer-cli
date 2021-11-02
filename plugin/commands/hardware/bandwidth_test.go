package hardware_test

import (

	"time"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/hardware"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/datatypes"
)

var _ = Describe("Hardware bandwidth", func() {
	var (
		fakeUI        *terminal.FakeUI
		fakeManager   *testhelpers.FakeHardwareServerManager
		cmd           *hardware.BandwidthCommand
		cliCommand    cli.Command
		fakeTransport *testhelpers.FakeTransportHandler
		fakeSession   *session.Session
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeManager = new(testhelpers.FakeHardwareServerManager)
		bleg := []string{}
		fakeSession = testhelpers.NewFakeSoftlayerSession(bleg)
		fakeTransport = new(testhelpers.FakeTransportHandler)
		cmd = hardware.NewBandwidthCommand(fakeUI, fakeManager)
		cliCommand = cli.Command{
			Name:        metadata.HardwareBandwidthMetaData().Name,
			Description: metadata.HardwareBandwidthMetaData().Description,
			Usage:       metadata.HardwareBandwidthMetaData().Usage,
			Flags:       metadata.HardwareBandwidthMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Hardware bandwidth", func() {
		Context("Argument Checking", func() {
			It("Error on missing ID", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
			It("Rollup specified", func() {
				testTime := "2021-08-01"
				err := testhelpers.RunCommand(cliCommand, "123456", "-s", testTime, "-e", testTime, "-r", "300")
				Expect(err).NotTo(HaveOccurred())
				// Expect(fakeUI.Outputs()).To(ContainSubstring("2021-08-10"))
				arg1, arg2, arg3, arg4 := fakeManager.GetBandwidthDataArgsForCall(0)
				Expect(arg1).To(Equal(123456))
				Expect(arg2.Format("2006-01-02")).To(Equal(testTime))
				Expect(arg3.Format("2006-01-02")).To(Equal(testTime))
				Expect(arg4).To(Equal(300))
			})
		})
		Context("DateTime parsing checks", func() {
			It("2006-01-02 Parsing works properly", func() {
				testTime := "2021-08-01"
				err := testhelpers.RunCommand(cliCommand, "123456", "-s", testTime, "-e", testTime)
				Expect(err).NotTo(HaveOccurred())
				// Expect(fakeUI.Outputs()).To(ContainSubstring("2021-08-10"))
				arg1, arg2, arg3, arg4 := fakeManager.GetBandwidthDataArgsForCall(0)
				Expect(arg1).To(Equal(123456))
				Expect(arg2.Format("2006-01-02")).To(Equal(testTime))
				Expect(arg3.Format("2006-01-02")).To(Equal(testTime))
				Expect(arg4).To(Equal(3600))
			})
			It("2006-01-02T15:04 Parsing works properly", func() {
				testTime := "2021-01-02T00:01"
				err := testhelpers.RunCommand(cliCommand, "123456", "-s", testTime, "-e", testTime)
				Expect(err).NotTo(HaveOccurred())
				arg1, arg2, arg3, arg4 := fakeManager.GetBandwidthDataArgsForCall(0)
				Expect(arg1).To(Equal(123456))
				Expect(arg2.Format("2006-01-02T15:04")).To(Equal(testTime))
				Expect(arg3.Format("2006-01-02T15:04")).To(Equal(testTime))
				Expect(arg4).To(Equal(3600))
			})
			It("2006-01-02T15:04:05-07:00 Parsing works properly", func() {
				testTime := "2021-01-02T00:01-05:00"
				err := testhelpers.RunCommand(cliCommand, "123456", "-s", testTime, "-e", testTime)
				Expect(err).NotTo(HaveOccurred())
				arg1, arg2, arg3, arg4 := fakeManager.GetBandwidthDataArgsForCall(0)
				Expect(arg1).To(Equal(123456))
				Expect(arg2.Format("2006-01-02T15:04-07:00")).To(Equal(testTime))
				Expect(arg3.Format("2006-01-02T15:04-07:00")).To(Equal(testTime))
				Expect(arg4).To(Equal(3600))
			})
			It("No time specified works properly", func() {
				testTime := time.Now()
				format := "2006-01-02T15:04-07:00"
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).NotTo(HaveOccurred())
				arg1, arg2, arg3, arg4 := fakeManager.GetBandwidthDataArgsForCall(0)
				Expect(arg1).To(Equal(123456))
				// Time has microsecond precision, so need to make sure we drop that part of when checking
				Expect(arg2.Format(format)).To(Equal(testTime.Format(format)))
				Expect(arg3.Format(format)).To(Equal(testTime.AddDate(0, -1, 0).Format(format)))
				Expect(arg4).To(Equal(3600))	
			})
			It("Bad Time", func() {
				testTime := "2021/01/03 00:01-05:00"
				err := testhelpers.RunCommand(cliCommand, "123456", "-s", testTime, "-e", "2021-01-02")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid start date: parsing time"))
				err = testhelpers.RunCommand(cliCommand, "123456", "-s", "2021-01-02", "-e", testTime)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid end date: parsing time"))
			})
		})
		Context("Build a proper table", func() {
			var returnData []datatypes.Metric_Tracking_Object_Data 
			var testTime string
			BeforeEach(func() {
				errAPI := fakeTransport.DoRequest(fakeSession, "SoftLayer_Metric_Tracking_Object",
												  "getBandwidthData", nil, nil, &returnData)
				Expect(errAPI).NotTo(HaveOccurred())
				testTime = "2021-08-01"
			})
			It("Default output", func() {
				fakeManager.GetBandwidthDataReturns(returnData, nil)
				err := testhelpers.RunCommand(cliCommand, "123456", "-s", testTime, "-e", testTime)
				Expect(err).NotTo(HaveOccurred())
				outputs := fakeUI.Outputs()
				Expect(outputs).To(ContainSubstring("Pub In    0.0032   0.2689         0.0016   2021-07-31 23:00"))
				Expect(outputs).To(ContainSubstring("2021-07-31 23:00   0.0016   0.0017    0.0000   0.0000"))
				
			})
			It("Quiet output", func() {
				fakeManager.GetBandwidthDataReturns(returnData, nil)
				err := testhelpers.RunCommand(cliCommand, "123456", "-s", testTime, "-e", testTime, "-q")
				Expect(err).NotTo(HaveOccurred())
				outputs := fakeUI.Outputs()
				Expect(outputs).To(ContainSubstring("Pub In    0.0032   0.2689         0.0016   2021-07-31 23:00"))
				Expect(outputs).NotTo(ContainSubstring("2021-07-31 23:00   0.0016   0.0017    0.0000   0.0000"))
				
			})
			It("Empty Response", func() {
				fakeManager.GetBandwidthDataReturns([]datatypes.Metric_Tracking_Object_Data{}, nil)
				err := testhelpers.RunCommand(cliCommand, "123456", "-s", testTime, "-e", testTime)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("No data"))
				
			})
		})
		
	})
})
