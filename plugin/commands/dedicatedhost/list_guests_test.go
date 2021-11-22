package dedicatedhost_test

import (
	"errors"
	"strings"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/sl"
	"github.com/urfave/cli"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/commands/dedicatedhost"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/testhelpers"
)

var _ = Describe("Dedicated host guests list", func() {
	var (
		fakeUI                   *terminal.FakeUI
		FakeDedicatedhostManager *testhelpers.FakeDedicatedhostManager
		cmd                      *dedicatedhost.ListGuestsCommand
		cliCommand               cli.Command
	)
	BeforeEach(func() {
		fakeUI = terminal.NewFakeUI()
		FakeDedicatedhostManager = new(testhelpers.FakeDedicatedhostManager)
		cmd = dedicatedhost.NewListGuestsCommand(fakeUI, FakeDedicatedhostManager)
		cliCommand = cli.Command{
			Name:        metadata.DedicatedhostListGuestsMetaData().Name,
			Description: metadata.DedicatedhostListGuestsMetaData().Description,
			Usage:       metadata.DedicatedhostListGuestsMetaData().Usage,
			Flags:       metadata.DedicatedhostListGuestsMetaData().Flags,
			Action:      cmd.Run,
		}
	})

	Describe("Guests list", func() {
		Context("Guests list without host id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand)
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: This command requires one argument.")).To(BeTrue())
			})
		})

		Context("Guests list with wrong host id", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Invalid input for 'Host ID'. It must be a positive integer.")).To(BeTrue())
			})
		})

		Context("Guests list with wrong column", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--column", "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: --column abc is not supported.")).To(BeTrue())
			})
		})

		Context("Guests list with wrong columns", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--column", "id", "--column", "username", "--column", "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: --column abc is not supported.")).To(BeTrue())
			})
		})
		Context("Guests list with wrong column", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--columns", "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: --columns abc is not supported.")).To(BeTrue())
			})
		})

		Context("Guests list with wrong columns", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--columns", "id", "--columns", "username", "--columns", "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: --columns abc is not supported.")).To(BeTrue())
			})
		})
		Context("Guests list with wrong sortby", func() {
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "--sortby", "abc")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Incorrect Usage: --sortby abc is not supported.")).To(BeTrue())
			})
		})

		Context("Guests list but server API call fails", func() {
			BeforeEach(func() {
				FakeDedicatedhostManager.ListGuestsReturns(nil, errors.New("Server Internal Error"))
			})
			It("return error", func() {
				err := testhelpers.RunCommand(cliCommand, "")
				Expect(err).To(HaveOccurred())
				Expect(strings.Contains(err.Error(), "Failed to list the host guest on your account.")).To(BeTrue())
				Expect(strings.Contains(err.Error(), "Server Internal Error")).To(BeTrue())
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
				err := testhelpers.RunCommand(cliCommand, "--sortby", "id")
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
				err := testhelpers.RunCommand(cliCommand, "--sortby", "hostname")
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
				err := testhelpers.RunCommand(cliCommand, "--sortby", "datacenter")
				Expect(err).NotTo(HaveOccurred())
				result := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(result[1], "dal10")).To(BeTrue())
			})
		})

		Context("Guests list with sortby=created_by", func() {
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
				err := testhelpers.RunCommand(cliCommand, "--sortby", "created_by", "--column", "created_by")
				Expect(err).NotTo(HaveOccurred())
				result := strings.Split(fakeUI.Outputs(), "\n")
				Expect(strings.Contains(result[1], "Anne Clark")).To(BeTrue())
			})
		})
	})
})
