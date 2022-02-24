package dedicatedhost_test

import (
	"errors"
	"time"

	. "github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/matchers"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dedicatedhost"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Dedicated host detail", func() {
	var (
		fakeUI                   *terminal.FakeUI
		FakeDedicatedhostManager *testhelpers.FakeDedicatedhostManager
		cmd                      *dedicatedhost.DetailCommand
		cliCommand               cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeDedicatedhostManager = new(testhelpers.FakeDedicatedhostManager)
		cmd = dedicatedhost.NewDetailCommand(fakeUI, FakeDedicatedhostManager)
		cliCommand = cli.Command{
			Name:        dedicatedhost.DedicatedhostDetailMetaData().Name,
			Description: dedicatedhost.DedicatedhostDetailMetaData().Description,
			Usage:       dedicatedhost.DedicatedhostDetailMetaData().Usage,
			Flags:       dedicatedhost.DedicatedhostDetailMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Dedicatedhost detail", func() {
		Context("Dedicatedhost detail without ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})
		})
		Context("Dedicatedhost detail with wrong VS ID", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Host ID'. It must be a positive integer."))
			})
		})

		Context("Dedicatedhost detail with server fails", func() {
			BeforeEach(func() {
				FakeDedicatedhostManager.GetInstanceReturns(datatypes.Virtual_DedicatedHost{}, errors.New("Internal Server Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get dedicatedhost instance: 1234.\n"))
				Expect(err.Error()).To(ContainSubstring("Internal Server Error"))
			})
		})

		Context("Dedicatedhost detail with correct VS ID ", func() {
			created, _ := time.Parse(time.RFC3339, "2022-02-01T00:00:00Z")
			modified, _ := time.Parse(time.RFC3339, "2022-02-01T00:00:00Z")
			BeforeEach(func() {
				FakeDedicatedhostManager.GetInstanceReturns(datatypes.Virtual_DedicatedHost{
					Id:             sl.Int(1234),
					Name:           sl.String("dedicatedhost"),
					CpuCount:       sl.Int(56),
					CreateDate:     sl.Time(created),
					DiskCapacity:   sl.Int(1200),
					MemoryCapacity: sl.Int(242),
					ModifyDate:     sl.Time(modified),
					GuestCount:     sl.Uint(3),
					BillingItem: &datatypes.Billing_Item_Virtual_DedicatedHost{
						Billing_Item: datatypes.Billing_Item{
							Id: sl.Int(1234567),
							Children: []datatypes.Billing_Item{
								datatypes.Billing_Item{
									NextInvoiceTotalRecurringAmount: sl.Float(10),
								},
							},
							OrderItem: &datatypes.Billing_Order_Item{
								Order: &datatypes.Billing_Order{
									UserRecord: &datatypes.User_Customer{
										Username: sl.String("wilmawang"),
									},
								},
							},
							NextInvoiceTotalRecurringAmount: sl.Float(1000.00),
						},
					},
					Datacenter: &datatypes.Location{
						Id:       sl.Int(1854895),
						LongName: sl.String("Dallas 13"),
						Name:     sl.String("dal13"),
					},
					Guests: []datatypes.Virtual_Guest{
						datatypes.Virtual_Guest{
							Domain:   sl.String("test.com"),
							Hostname: sl.String("test"),
							Id:       sl.Int(1234567),
							Uuid:     sl.String("9131111-2222-6a10-3333-992c544444"),
						},
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1234"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"dedicatedhost"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"56"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1200"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"242"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2022-02-01T00:00:00Z"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2022-02-01T00:00:00Z"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"3"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"dal13"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"wilmawang"}))
			})
			It("return no error", func() {
				err := testhelpers.RunCommand(cliCommand, "1234", "--guests", "--price")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1234"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"dedicatedhost"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"56"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1200"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"242"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2022-02-01T00:00:00Z"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"2022-02-01T00:00:00Z"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"3"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"dal13"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"10"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"wilmawang"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"test.com"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"test"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"9131111-2222-6a10-3333-992c544444"}))
				Expect(fakeUI.Outputs()).To(ContainSubstrings([]string{"1234567"}))
			})
		})
	})
})
