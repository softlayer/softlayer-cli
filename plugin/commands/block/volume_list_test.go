package block_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/block"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Volume list", func() {
	var (
		fakeUI             *terminal.FakeUI
		FakeStorageManager *testhelpers.FakeStorageManager
		cliCommand         *block.VolumeListCommand
		fakeSession        *session.Session
		slCommand          *metadata.SoftlayerStorageCommand
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeStorageManager = new(testhelpers.FakeStorageManager)
		slCommand = metadata.NewSoftlayerStorageCommand(fakeUI, fakeSession, "block")
		cliCommand = block.NewVolumeListCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		cliCommand.StorageManager = FakeStorageManager
	})

	Describe("Volume list", func() {
		Context("Volume list with wrong column", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--column", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --column abc is not supported."))
			})
		})
		Context("Volume list with wrong columns", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--column", "id", "--column", "username", "--column", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --column abc is not supported."))
			})
		})
		Context("Volume list with wrong sortby", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --sortby abc is not supported."))
			})
		})

		Context("Volume list but server API call fails", func() {
			BeforeEach(func() {
				FakeStorageManager.ListVolumesReturns(nil, errors.New("Server Internal Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to list volumes on your account."))
				Expect(err.Error()).To(ContainSubstring("Server Internal Error"))
			})
		})

		Context("Volume list with sortby=id", func() {
			BeforeEach(func() {
				FakeStorageManager.ListVolumesReturns([]datatypes.Network_Storage{
					datatypes.Network_Storage{
						Id: sl.Int(123458),
					},
					datatypes.Network_Storage{
						Id: sl.Int(123457),
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "id")
				Expect(err).NotTo(HaveOccurred())
				result := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(result[1], "123457")).To(BeTrue())
				Expect(strings.Contains(result[2], "123458")).To(BeTrue())
			})
		})

		Context("Volume list with sortby=name", func() {
			BeforeEach(func() {
				FakeStorageManager.ListVolumesReturns([]datatypes.Network_Storage{
					datatypes.Network_Storage{
						Id:       sl.Int(123458),
						Username: sl.String("myvolume"),
					},
					datatypes.Network_Storage{
						Id:       sl.Int(123457),
						Username: sl.String("hisvolume"),
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "username")
				Expect(err).NotTo(HaveOccurred())
				result := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(result[1], "hisvolume")).To(BeTrue())
				Expect(strings.Contains(result[2], "myvolume")).To(BeTrue())
			})
		})

		Context("Volume list with sortby=datacenter", func() {
			BeforeEach(func() {
				FakeStorageManager.ListVolumesReturns([]datatypes.Network_Storage{
					datatypes.Network_Storage{
						ServiceResource: &datatypes.Network_Service_Resource{
							Datacenter: &datatypes.Location{
								Name: sl.String("tok02"),
							},
						},
					},
					datatypes.Network_Storage{
						ServiceResource: &datatypes.Network_Service_Resource{
							Datacenter: &datatypes.Location{
								Name: sl.String("dal10"),
							},
						},
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "datacenter")
				Expect(err).NotTo(HaveOccurred())
				result := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(result[1], "dal10")).To(BeTrue())
				Expect(strings.Contains(result[2], "tok02")).To(BeTrue())
			})
		})

		Context("Volume list with sortby=storage_type", func() {
			BeforeEach(func() {
				FakeStorageManager.ListVolumesReturns([]datatypes.Network_Storage{
					datatypes.Network_Storage{
						StorageType: &datatypes.Network_Storage_Type{
							KeyName: sl.String("performance"),
						},
					},
					datatypes.Network_Storage{
						StorageType: &datatypes.Network_Storage_Type{
							KeyName: sl.String("enduration"),
						},
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "storage_type")
				Expect(err).NotTo(HaveOccurred())
				result := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(result[1], "enduration")).To(BeTrue())
				Expect(strings.Contains(result[2], "performance")).To(BeTrue())
			})
		})

		Context("Volume list with sortby=capacity_gb", func() {
			BeforeEach(func() {
				FakeStorageManager.ListVolumesReturns([]datatypes.Network_Storage{
					datatypes.Network_Storage{
						CapacityGb: sl.Int(1000),
					},
					datatypes.Network_Storage{
						CapacityGb: sl.Int(2000),
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "capacity_gb")
				Expect(err).NotTo(HaveOccurred())
				result := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(result[1], "1000")).To(BeTrue())
				Expect(strings.Contains(result[2], "2000")).To(BeTrue())
			})
		})

		Context("Volume list with sortby=bytes_used", func() {
			BeforeEach(func() {
				FakeStorageManager.ListVolumesReturns([]datatypes.Network_Storage{
					datatypes.Network_Storage{
						BytesUsed: sl.String("1020"),
					},
					datatypes.Network_Storage{
						BytesUsed: sl.String("600"),
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "bytes_used", "--column", "bytes_used")
				Expect(err).NotTo(HaveOccurred())
				result := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(result[1], "600")).To(BeTrue())
				Expect(strings.Contains(result[2], "1020")).To(BeTrue())
			})
		})

		Context("Volume list with sortby=ip_addr", func() {
			BeforeEach(func() {
				FakeStorageManager.ListVolumesReturns([]datatypes.Network_Storage{
					datatypes.Network_Storage{
						ServiceResourceBackendIpAddress: sl.String("9.4.6.4"),
					},
					datatypes.Network_Storage{
						ServiceResourceBackendIpAddress: sl.String("6.7.8.9"),
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "ip_addr", "--column", "ip_addr")
				Expect(err).NotTo(HaveOccurred())
				result := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(result[1], "6.7.8.9")).To(BeTrue())
				Expect(strings.Contains(result[2], "9.4.6.4")).To(BeTrue())
			})
		})

		Context("Volume list with sortby=lunId", func() {
			BeforeEach(func() {
				FakeStorageManager.ListVolumesReturns([]datatypes.Network_Storage{
					datatypes.Network_Storage{
						LunId: sl.String("125"),
					},
					datatypes.Network_Storage{
						LunId: sl.String("67"),
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "lunId")
				Expect(err).NotTo(HaveOccurred())
				result := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(result[1], "67")).To(BeTrue())
				Expect(strings.Contains(result[2], "125")).To(BeTrue())
			})
		})

		Context("Volume list with sortby=active_transactions", func() {
			BeforeEach(func() {
				FakeStorageManager.ListVolumesReturns([]datatypes.Network_Storage{
					datatypes.Network_Storage{
						ActiveTransactionCount: sl.Uint(uint(2)),
					},
					datatypes.Network_Storage{
						ActiveTransactionCount: sl.Uint(uint(1)),
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "active_transactions", "--column", "active_transactions")
				Expect(err).NotTo(HaveOccurred())
				result := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(result[1], "1")).To(BeTrue())
				Expect(strings.Contains(result[2], "2")).To(BeTrue())
			})
		})

		Context("Volume list with sortby=created_by", func() {
			BeforeEach(func() {
				FakeStorageManager.ListVolumesReturns([]datatypes.Network_Storage{
					datatypes.Network_Storage{
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
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "created_by", "--column", "created_by")
				Expect(err).NotTo(HaveOccurred())
				result := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(result[1], "Anne Clark")).To(BeTrue())
				Expect(strings.Contains(result[2], "Bill Jones")).To(BeTrue())
			})

		})
	})
})
