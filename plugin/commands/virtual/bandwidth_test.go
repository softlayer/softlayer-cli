package virtual_test

import (
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/virtual"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("VS bandwidth", func() {
	var (
		fakeUI        *terminal.FakeUI
		cliCommand    *virtual.BandwidthCommand
		fakeSession   *session.Session
		slCommand     *metadata.SoftlayerCommand
		fakeVSManager *testhelpers.FakeVirtualServerManager
		fakeTransport *testhelpers.FakeTransportHandler
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeVSManager = new(testhelpers.FakeVirtualServerManager)
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = virtual.NewBandwidthCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.VirtualServerManager = fakeVSManager
		fakeTransport = new(testhelpers.FakeTransportHandler)
	})

	Describe("VS bandwidth", func() {
		Context("Argument Checking", func() {
			It("Error on missing ID", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
			It("Rollup specified", func() {
				testTime := "2021-08-01"
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "-s", testTime, "-e", testTime, "-r", "300")
				Expect(err).NotTo(HaveOccurred())
				// Expect(fakeUI.Outputs()).To(ContainSubstring("2021-08-10"))
				arg1, arg2, arg3, arg4 := fakeVSManager.GetBandwidthDataArgsForCall(0)
				Expect(arg1).To(Equal(123456))
				Expect(arg2.Format("2006-01-02")).To(Equal(testTime))
				Expect(arg3.Format("2006-01-02")).To(Equal(testTime))
				Expect(arg4).To(Equal(300))
			})
		})
		Context("DateTime parsing checks", func() {
			It("2006-01-02 Parsing works properly", func() {
				testTime := "2021-08-01"
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "-s", testTime, "-e", testTime)
				Expect(err).NotTo(HaveOccurred())
				// Expect(fakeUI.Outputs()).To(ContainSubstring("2021-08-10"))
				arg1, arg2, arg3, arg4 := fakeVSManager.GetBandwidthDataArgsForCall(0)
				Expect(arg1).To(Equal(123456))
				Expect(arg2.Format("2006-01-02")).To(Equal(testTime))
				Expect(arg3.Format("2006-01-02")).To(Equal(testTime))
				Expect(arg4).To(Equal(3600))
			})
			It("2006-01-02T15:04 Parsing works properly", func() {
				testTime := "2021-01-02T00:01"
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "-s", testTime, "-e", testTime)
				Expect(err).NotTo(HaveOccurred())
				arg1, arg2, arg3, arg4 := fakeVSManager.GetBandwidthDataArgsForCall(0)
				Expect(arg1).To(Equal(123456))
				Expect(arg2.Format("2006-01-02T15:04")).To(Equal(testTime))
				Expect(arg3.Format("2006-01-02T15:04")).To(Equal(testTime))
				Expect(arg4).To(Equal(3600))
			})
			It("2006-01-02T15:04:05-07:00 Parsing works properly", func() {
				testTime := "2021-01-02T00:01-05:00"
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "-s", testTime, "-e", testTime)
				Expect(err).NotTo(HaveOccurred())
				arg1, arg2, arg3, arg4 := fakeVSManager.GetBandwidthDataArgsForCall(0)
				Expect(arg1).To(Equal(123456))
				Expect(arg2.Format("2006-01-02T15:04-07:00")).To(Equal(testTime))
				Expect(arg3.Format("2006-01-02T15:04-07:00")).To(Equal(testTime))
				Expect(arg4).To(Equal(3600))
			})
			It("No time specified works properly", func() {
				testTime := time.Now()
				format := "2006-01-02T15:04-07:00"
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456")
				Expect(err).NotTo(HaveOccurred())
				arg1, arg2, arg3, arg4 := fakeVSManager.GetBandwidthDataArgsForCall(0)
				Expect(arg1).To(Equal(123456))
				// Time has microsecond precision, so need to make sure we drop that part of when checking
				Expect(arg2.Format(format)).To(Equal(testTime.Format(format)))
				Expect(arg3.Format(format)).To(Equal(testTime.AddDate(0, -1, 0).Format(format)))
				Expect(arg4).To(Equal(3600))
			})
			It("Bad Time", func() {
				testTime := "2021/01/03 00:01-05:00"
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "-s", testTime, "-e", "2021-01-02")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid start date: parsing time"))
				err = testhelpers.RunCobraCommand(cliCommand.Command, "123456", "-s", "2021-01-02", "-e", testTime)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid end date: parsing time"))
			})
		})
		Context("Build a proper table", func() {
			var returnData []datatypes.Metric_Tracking_Object_Data
			var testTime string
			BeforeEach(func() {
				// Loads data from testfixtrues into returnData
				errAPI := fakeTransport.DoRequest(fakeSession, "SoftLayer_Metric_Tracking_Object",
					"getBandwidthData", nil, nil, &returnData)
				Expect(errAPI).NotTo(HaveOccurred())
				testTime = "2021-08-01"
			})
			It("Default output", func() {
				fakeVSManager.GetBandwidthDataReturns(returnData, nil)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "-s", testTime, "-e", testTime)
				Expect(err).NotTo(HaveOccurred())
				outputs := fakeUI.Outputs()
				Expect(outputs).To(ContainSubstring("Pub In    0.0032   0.2689         0.0016   2021-07-31 23:00"))
				Expect(outputs).To(ContainSubstring("2021-07-31 23:00   0.0016   0.0017    0.0000   0.0000"))

			})
			It("Quiet output", func() {
				fakeVSManager.GetBandwidthDataReturns(returnData, nil)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "-s", testTime, "-e", testTime, "-q")
				Expect(err).NotTo(HaveOccurred())
				outputs := fakeUI.Outputs()
				Expect(outputs).To(ContainSubstring("Pub In    0.0032   0.2689         0.0016   2021-07-31 23:00"))
				Expect(outputs).NotTo(ContainSubstring("2021-07-31 23:00   0.0016   0.0017    0.0000   0.0000"))

			})
			It("Empty Response", func() {
				fakeVSManager.GetBandwidthDataReturns([]datatypes.Metric_Tracking_Object_Data{}, nil)
				err := testhelpers.RunCobraCommand(cliCommand.Command, "123456", "-s", testTime, "-e", testTime)
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("No data"))

			})
		})

	})
})
