package file_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/file"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Volume list", func() {
	var (
		fakeUI             *terminal.FakeUI
		cliCommand         *file.VolumeListCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerStorageCommand
		FakeStorageManager *testhelpers.FakeStorageManager
		fakeHandler 	   *testhelpers.FakeTransportHandler
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		fakeHandler = testhelpers.GetSessionHandler(fakeSession)
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "file")
		cliCommand = file.NewVolumeListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		
	})
	AfterEach(func() {
		fakeHandler.ClearApiCallLogs()
		fakeHandler.ClearErrors()
	})
	Describe("Volume list", func() {
		Context("Usage Errors", func() {
			BeforeEach(func() {
				// Dont use fixture data for these tests
				cliCommand.StorageManager = FakeStorageManager
			})
			It("Volume list with wrong column", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--column", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --column abc is not supported."))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--column", "id", "--column", "username", "--column", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --column abc is not supported."))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --sortby abc is not supported."))
			})
		})

		Context("API Errors", func() {
			BeforeEach(func() {
				// Dont use fixture data for these tests
				cliCommand.StorageManager = FakeStorageManager
				FakeStorageManager.ListVolumesReturns(nil, errors.New("Server Internal Error"))
			})
			It("Volume list but server API call fails", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to list volumes on your account."))
				Expect(err.Error()).To(ContainSubstring("Server Internal Error"))
			})
		})

		Context("Sorting", func() {
			BeforeEach(func() {
				// Dont use fixture data for these tests
				cliCommand.StorageManager = FakeStorageManager
				FakeStorageManager.ListVolumesReturns([]datatypes.Network_Storage{
					datatypes.Network_Storage{
						Id: sl.Int(123458),
						Username: sl.String("myvolume"),
						ServiceResource: &datatypes.Network_Service_Resource{
							Datacenter: &datatypes.Location{
								Name: sl.String("tok02"),
							},
						},
						ServiceResourceBackendIpAddress: sl.String("9.4.6.4"),
						StorageType: &datatypes.Network_Storage_Type{
							KeyName: sl.String("performance"),
						},
						CapacityGb: sl.Int(1000),
						BytesUsed: sl.String("1020"),
						ActiveTransactionCount: sl.Uint(uint(2)),
						BillingItem: &datatypes.Billing_Item{
							OrderItem: &datatypes.Billing_Order_Item{
								Order: &datatypes.Billing_Order{
									UserRecord: &datatypes.User_Customer{
										Username: sl.String("Bill Jones"),
									},
								},
							},
						},
					},
					datatypes.Network_Storage{
						Id: sl.Int(123457),
						Username: sl.String("hisvolume"),
						ServiceResource: &datatypes.Network_Service_Resource{
							Datacenter: &datatypes.Location{
								Name: sl.String("dal10"),
							},
						},
						ServiceResourceBackendIpAddress: sl.String("6.7.8.9"),
						StorageType: &datatypes.Network_Storage_Type{
							KeyName: sl.String("enduration"),
						},
						CapacityGb: sl.Int(2000),
						BytesUsed: sl.String("600"),
						ActiveTransactionCount: sl.Uint(uint(1)),
						BillingItem: &datatypes.Billing_Item{
							OrderItem: &datatypes.Billing_Order_Item{
								Order: &datatypes.Billing_Order{
									UserRecord: &datatypes.User_Customer{
										Username: sl.String("Anne Clark"),
									},
								},
							},
						},
					},
				}, nil)
			})
			It("Volume list with sortby=id", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "id")
				Expect(err).NotTo(HaveOccurred())
				result := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(result[1], "123457"))
				Expect(strings.Contains(result[2], "123458"))
			})
			It("Volume list with sortby=name", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "username")
				Expect(err).NotTo(HaveOccurred())
				result := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(result[1], "hisvolume"))
				Expect(strings.Contains(result[2], "myvolume"))
			})
			It("Volume list with sortby=datacenter", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "datacenter")
				Expect(err).NotTo(HaveOccurred())
				result := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(result[1], "dal10"))
				Expect(strings.Contains(result[2], "tok02"))
			})
			It("Volume list with sortby=storage_type", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "storage_type")
				Expect(err).NotTo(HaveOccurred())
				result := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(result[1], "enduration"))
				Expect(strings.Contains(result[2], "performance"))
			})
			It("Volume list with sortby=capacity_gb", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "capacity_gb")
				Expect(err).NotTo(HaveOccurred())
				result := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(result[1], "1000"))
				Expect(strings.Contains(result[2], "2000"))
			})
			It("Volume list with sortby=bytes_used", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "bytes_used", "--column", "bytes_used")
				Expect(err).NotTo(HaveOccurred())
				result := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(result[1], "600"))
				Expect(strings.Contains(result[2], "1020"))
			})
			It("Volume list with sortby=ip_addr", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "ip_addr")
				Expect(err).NotTo(HaveOccurred())
				result := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(result[1], "6.7.8.9"))
				Expect(strings.Contains(result[2], "9.4.6.4"))
			})
			It("Volume list with sortby=active_transactions", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "active_transactions", "--column", "active_transactions")
				Expect(err).NotTo(HaveOccurred())
				result := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(result[1], "1"))
				Expect(strings.Contains(result[2], "2"))
			})
			It("Volume list with sortby=created_by", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "created_by", "--column", "created_by")
				Expect(err).NotTo(HaveOccurred())
				result := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(result[1], "Anne Clark"))
				Expect(strings.Contains(result[2], "Bill Jon1es"))
			})
		})
		Context("Github Issues #937", func() {
			// v1.6.0 of github.com/IBM-Cloud/ibm-cloud-cli-sdk introduced this change in line breaks
			It("Volume list with long column data", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).NotTo(HaveOccurred())
				result := strings.Split(fakeUI.Outputs(), "\n")
				Expect(result[2]).To(ContainSubstring("Lorem ipsum dolor sit amet,"))
				// This is where we started getting unexpected line breaks
				Expect(result[2]).To(ContainSubstring("eiusmod tempor incididunt ut labore"))
				
			})

		})
	})
})
