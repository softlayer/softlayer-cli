package account_test

import (
	"errors"
	"time"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/account"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("account bandwidth_pools_details", func() {
	var (
		fakeUI             *terminal.FakeUI
		fakeAccountManager *testhelpers.FakeAccountManager
		cmd                *account.BandwidthPoolsDetailCommand
		cliCommand         cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeAccountManager = new(testhelpers.FakeAccountManager)
		cmd = account.NewBandwidthPoolsDetailCommand(fakeUI, fakeAccountManager)
		cliCommand = cli.Command{
			Name:        account.BandwidthPoolsDetailMetaData().Name,
			Description: account.BandwidthPoolsDetailMetaData().Description,
			Usage:       account.BandwidthPoolsDetailMetaData().Usage,
			Flags:       account.BandwidthPoolsDetailMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("account bandwidth_pools_details", func() {

		Context("Return error", func() {
			It("Set command without Id", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument."))
			})

			It("Set command with an invalid Id", func() {
				err := testhelpers.RunCommand(cliCommand, "abcde")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Bandwidth Pool ID'. It must be a positive integer."))
			})

			It("Set invalid output", func() {
				err := testhelpers.RunCommand(cliCommand, "123456", "--output=xml")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: Invalid output format, only JSON is supported now."))
			})
		})

		Context("Return error", func() {
			BeforeEach(func() {
				fakeAccountManager.GetBandwidthPoolDetailReturns(datatypes.Network_Bandwidth_Version1_Allotment{}, errors.New("Failed to get Bandwidth Pool."))
			})
			It("Failed Bandwidth Pool", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to get Bandwidth Pool."))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2017-11-08T00:00:00Z")
				fakerBandwidthPool := datatypes.Network_Bandwidth_Version1_Allotment{
					Id:         sl.Int(123456),
					Name:       sl.String("Bandwidth Pool 1"),
					CreateDate: sl.Time(created),
					BillingCyclePublicBandwidthUsage: &datatypes.Network_Bandwidth_Usage{
						AmountOut: sl.Float(1000000000.0),
					},
					ProjectedPublicBandwidthUsage: sl.Float(2000000000.0),
					InboundPublicBandwidthUsage:   sl.Float(3000000000.0),
					Hardware: []datatypes.Hardware{
						datatypes.Hardware{
							Id:                       sl.Int(111111),
							FullyQualifiedDomainName: sl.String("hardware.mydomain.com"),
							PrimaryIpAddress:         sl.String("11.11.11.11"),
							BandwidthAllotmentDetail: &datatypes.Network_Bandwidth_Version1_Allotment_Detail{
								Allocation: &datatypes.Network_Bandwidth_Version1_Allocation{
									Amount: sl.Float(4000000000.0),
								},
							},
							OutboundBandwidthUsage: sl.Float(5000000000.0),
						},
					},
					VirtualGuests: []datatypes.Virtual_Guest{
						datatypes.Virtual_Guest{
							Id:                       sl.Int(222222),
							FullyQualifiedDomainName: sl.String("virtualguest.mydomain.com"),
							PrimaryIpAddress:         sl.String("22.22.22.22"),
							BandwidthAllotmentDetail: &datatypes.Network_Bandwidth_Version1_Allotment_Detail{
								Allocation: &datatypes.Network_Bandwidth_Version1_Allocation{
									Amount: sl.Float(6000000000.0),
								},
							},
							OutboundPublicBandwidthUsage: sl.Float(7000000000.0),
						},
					},
					BareMetalInstances: []datatypes.Hardware{
						datatypes.Hardware{
							Id:                       sl.Int(333333),
							FullyQualifiedDomainName: sl.String("baremetal.mydomain.com"),
							PrimaryIpAddress:         sl.String("33.33.33.33"),
							BandwidthAllotmentDetail: &datatypes.Network_Bandwidth_Version1_Allotment_Detail{
								Allocation: &datatypes.Network_Bandwidth_Version1_Allocation{
									Amount: sl.Float(8000000000.0),
								},
							},
							OutboundBandwidthUsage: sl.Float(9000000000.0),
						},
					},
				}
				fakeAccountManager.GetBandwidthPoolDetailReturns(fakerBandwidthPool, nil)
			})
			It("Get Bandwidth Pool with devices", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("2017-11-08T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Bandwidth Pool 1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("2.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("3.00 GB"))

				Expect(fakeUI.Outputs()).To(ContainSubstring("111111"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("hardware.mydomain.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("1.11.11.11"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("4.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("5.00 GB"))

				Expect(fakeUI.Outputs()).To(ContainSubstring("222222"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("virtualguest.mydomain.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("22.22.22.22"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("6.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("7.00 GB"))

				Expect(fakeUI.Outputs()).To(ContainSubstring("333333"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("baremetal.mydomain.com"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("3.33.33.33"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("8.00 GB"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("9.00 GB"))
			})
		})

		Context("Return no error", func() {
			BeforeEach(func() {
				created, _ := time.Parse(time.RFC3339, "2017-11-08T00:00:00Z")
				fakerBandwidthPool := datatypes.Network_Bandwidth_Version1_Allotment{
					Id:         sl.Int(123456),
					Name:       sl.String("Bandwidth Pool 1"),
					CreateDate: sl.Time(created),
					BillingCyclePublicBandwidthUsage: &datatypes.Network_Bandwidth_Usage{
						AmountOut: sl.Float(1000000000.0),
					},
					ProjectedPublicBandwidthUsage: sl.Float(2000000000.0),
					InboundPublicBandwidthUsage:   sl.Float(3000000000.0),
					Hardware:                      []datatypes.Hardware{},
					VirtualGuests:                 []datatypes.Virtual_Guest{},
					BareMetalInstances:            []datatypes.Hardware{},
				}
				fakeAccountManager.GetBandwidthPoolDetailReturns(fakerBandwidthPool, nil)
			})
			It("Get Bandwidth Pool with devices", func() {
				err := testhelpers.RunCommand(cliCommand, "123456")
				Expect(err).NotTo(HaveOccurred())
				Expect(fakeUI.Outputs()).To(ContainSubstring("2017-11-08T00:00:00Z"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("123456"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Bandwidth Pool 1"))
				Expect(fakeUI.Outputs()).To(ContainSubstring("Not Found"))
			})
		})
	})
})
