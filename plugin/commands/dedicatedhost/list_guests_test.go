package dedicatedhost_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dedicatedhost"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Dedicated host guests list", func() {
	var (
		fakeUI                   *terminal.FakeUI
		cliCommand               *dedicatedhost.ListGuestsCommand
		fakeSession              *session.Session
		slCommand                *metadata.SoftlayerCommand
		FakeDedicatedhostManager *testhelpers.FakeDedicatedHostManager
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		fakeSession = testhelpers.NewFakeSoftlayerSession([]string{})
		slCommand = metadata.NewSoftlayerCommand(fakeUI, fakeSession)
		cliCommand = dedicatedhost.NewListGuestsCommand(slCommand)
		cliCommand.Command.PersistentFlags().Var(cliCommand.OutputFlag, "output", "--output=JSON for json output.")
		FakeDedicatedhostManager = new(testhelpers.FakeDedicatedHostManager)
		cliCommand.DedicatedHostManager = FakeDedicatedhostManager
	})

	Describe("Guests list", func() {
		Context("Guests list without host id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: This command requires one argument"))
			})
		})

		Context("Guests list with wrong host id", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "abc")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid input for 'Host ID'. It must be a positive integer."))
			})
		})

		Context("Guests list with wrong column", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--column", "abc", "1234567")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --column abc is not supported."))
			})
		})

		Context("Guests list with wrong column", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--column", "abc", "1234567")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --column abc is not supported."))
			})
		})

		Context("Guests list with wrong column", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--column", "id", "--column", "hostname", "--column", "abc", "1234567")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --column abc is not supported."))
			})
		})
		Context("Guests list with wrong sortby", func() {
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "abc", "1234567")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Incorrect Usage: --sortby abc is not supported."))
			})
		})

		Context("Guests list but server API call fails", func() {
			BeforeEach(func() {
				FakeDedicatedhostManager.ListGuestsReturns(nil, errors.New("Server Internal Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "1234567")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Failed to list the host guest on your account."))
				Expect(err.Error()).To(ContainSubstring("Server Internal Error"))
			})
		})

		Context("Guests list with sortby=id", func() {
			BeforeEach(func() {
				FakeDedicatedhostManager.ListGuestsReturns([]datatypes.Virtual_Guest{
					datatypes.Virtual_Guest{
						Id: sl.Int(1234567),
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "id", "1234567")
				Expect(err).NotTo(HaveOccurred())
				result := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(result[1], "1234567")).To(BeTrue())
			})
		})

		Context("Guests list with sortby=hostname", func() {
			BeforeEach(func() {
				FakeDedicatedhostManager.ListGuestsReturns([]datatypes.Virtual_Guest{
					datatypes.Virtual_Guest{
						Id:       sl.Int(1234567),
						Hostname: sl.String("test"),
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "hostname", "1234567")
				Expect(err).NotTo(HaveOccurred())
				result := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(result[1], "test")).To(BeTrue())
			})
		})

		Context("Guests list with sortby=datacenter", func() {
			BeforeEach(func() {
				FakeDedicatedhostManager.ListGuestsReturns([]datatypes.Virtual_Guest{
					datatypes.Virtual_Guest{
						Id:       sl.Int(1234567),
						Hostname: sl.String("test"),
						Datacenter: &datatypes.Location{
							Name: sl.String("dal10"),
						},
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--sortby", "datacenter", "1234567")
				Expect(err).NotTo(HaveOccurred())
				result := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(result[1], "dal10")).To(BeTrue())
			})
		})

		Context("Guests list with colum=created_by", func() {
			BeforeEach(func() {
				FakeDedicatedhostManager.ListGuestsReturns([]datatypes.Virtual_Guest{
					datatypes.Virtual_Guest{
						BillingItem: &datatypes.Billing_Item_Virtual_Guest{
							Billing_Item: datatypes.Billing_Item{
								OrderItem: &datatypes.Billing_Order_Item{
									Order: &datatypes.Billing_Order{
										UserRecord: &datatypes.User_Customer{
											Username: sl.String("Anne Clark"),
										},
									},
								},
							},
						},
					},
				}, nil)
			})
			It("return no error", func() {
				err := testhelpers.RunCobraCommand(cliCommand.Command, "--column", "created_by", "1234567")
				Expect(err).NotTo(HaveOccurred())
				result := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(result[1], "Anne Clark")).To(BeTrue())
			})
		})
	})
})
